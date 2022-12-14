package services

import (
	"archive/zip"
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/greatfocus/archive-service/database"
	"github.com/greatfocus/archive-service/models"
)

// ExtractService struct
type ExtractService struct {
	database *database.Conn
}

// Init method
func (e *ExtractService) Init(db *database.Conn) {
	e.database = db
}
func (e *ExtractService) CreateExtract(ctx context.Context, req *models.Request) (*models.Request, error) {
	_, err := e.insertRecordToDB(ctx, req)
	if err != nil {
		return req, err
	}

	// Check for background execution
	if req.Background {
		return req, nil
	} else {
		// update the database
		return e.InitiateExtract(ctx, req)
	}
}

func (e *ExtractService) InitiateExtract(ctx context.Context, req *models.Request) (*models.Request, error) {
	_, err := e.extractFiles(req)
	if err != nil {
		return req, err
	}

	req.Status = "done"
	err = e.updateStatus(ctx, req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (e *ExtractService) extractFiles(req *models.Request) (*models.Request, error) {
	hasFilteredName := false
	hasPartialExtraction := false
	zipPath := req.Dir + req.File
	read, err := zip.OpenReader(zipPath)
	if err != nil {
		return req, errors.New("failed to open file")
	}
	defer read.Close()

	// filter names
	var fileNames = make(map[string]string)
	var partiallyFiltered = strings.Split(req.PartialExtraction, "|")
	if len(req.FilteredNames) > 1 {
		hasFilteredName = true
		names := strings.Split(req.FilteredNames, "|")
		for _, n := range names {
			fileNames[n] = n
		}
	} else if len(partiallyFiltered) > 1 {
		hasPartialExtraction = true
	}

	for i, file := range read.File {
		if hasFilteredName && len(fileNames) > 0 {
			name := strings.TrimSuffix(file.Name, filepath.Ext(file.Name))
			namedFile := fileNames[name]
			if len(namedFile) < 1 {
				continue
			}
		} else if hasPartialExtraction {
			found := false
			for _, r := range partiallyFiltered {
				record, _ := strconv.ParseInt(r, 6, 12)
				if int(record-1) == i {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}

		zippedFile, err := file.Open()
		if err != nil {
			return req, err
		}
		defer zippedFile.Close()

		extractedFilePath := filepath.Join(
			req.Dir,
			file.Name,
		)
		if file.FileInfo().IsDir() {
			log.Println("Directory Created:", extractedFilePath)
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			log.Println("File extracted:", file.Name)

			outputFile, err := os.OpenFile(
				extractedFilePath,
				os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
				file.Mode(),
			)
			if err != nil {
				return req, err
			}
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedFile)
			if err != nil {
				return req, err
			}
		}
	}
	return req, nil
}

func (e *ExtractService) updateStatus(ctx context.Context, req *models.Request) error {
	query := `
    UPDATE extract SET status=? WHERE id=?;
  	`
	updated := e.database.Update(ctx, query, req.Status, req.ID)
	if !updated {
		return errors.New("failed to update extract")
	}
	return nil
}

func (e *ExtractService) insertRecordToDB(ctx context.Context, req *models.Request) (*models.Request, error) {
	req.ID = uuid.New().String()
	req.Status = "new"
	query := `
	insert into extract (id, fileName, dir, status, aligorithm, filteredNames, partialExtraction, background)
	VALUES(?,?,?,?,?,?,?,?);
	`
	_, inserted := e.database.Insert(ctx, query, req.ID, req.File, req.Dir, req.Status, req.Aligorithm, req.FilteredNames, req.PartialExtraction, req.Background)
	if !inserted {
		return req, errors.New("failed to insert extract")
	}
	return req, nil
}

func (e *ExtractService) GetStatus(ctx context.Context, id string) (models.Request, error) {
	query := `
	select id, fileName, dir, status, createdOn
	from extract
	where id = ?
	`
	row := e.database.Select(ctx, query, id)
	result := models.Request{}
	err := row.Scan(&result.ID, &result.File, &result.Dir, &result.Status, &result.CreatedOn)
	switch err {
	case sql.ErrNoRows:
		return result, nil
	case nil:
		// update cache
		return result, nil
	default:
		return result, err
	}
}
