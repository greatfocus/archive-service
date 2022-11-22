package services

import (
	"archive/zip"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/greatfocus/archive-service/database"
	"github.com/greatfocus/archive-service/models"
)

// ArchiveService struct
type ArchiveService struct {
	database *database.Conn
}

// Init method
func (a *ArchiveService) Init(db *database.Conn) {
	a.database = db
}

func (a *ArchiveService) CreateArchive(ctx context.Context, req *models.Request) (*models.Request, error) {
	_, err := a.insertRecordToDB(ctx, req)
	if err != nil {
		return req, err
	}

	// Check for background execution
	if req.Background {
		return req, nil
	} else {
		// update the database
		return a.InitiateArchive(ctx, req)
	}
}

func (a *ArchiveService) InitiateArchive(ctx context.Context, req *models.Request) (*models.Request, error) {
	_, err := a.archiveFiles(ctx, req)
	if err != nil {
		return req, err
	}

	req.Status = "done"
	err = a.updateStatus(ctx, req)
	if err != nil {
		return req, err
	}
	return req, nil
}

func (a *ArchiveService) archiveFiles(ctx context.Context, req *models.Request) (*models.Request, error) {
	fileNames, err := getListOfFileNames(req)
	if err != nil {
		return req, err
	}

	err = compress(fileNames, req)
	if err != nil {
		return req, err
	}

	return req, nil
}

func (a *ArchiveService) insertRecordToDB(ctx context.Context, req *models.Request) (*models.Request, error) {
	req.ID = uuid.New().String()
	req.Status = "new"
	query := `
	INSERT INTO archive(id, fileName, dir, status, aligorithm, filteredNames, background)
	VALUES($1,$2,$3,$4,$5,$6,$7);
	`
	_, inserted := a.database.Insert(ctx, query, req.ID, req.File, req.Dir, req.Status, req.Aligorithm, req.FilteredNames, req.Background)
	if !inserted {
		return req, errors.New("failed to insert archive")
	}
	return req, nil
}

func (a *ArchiveService) updateStatus(ctx context.Context, req *models.Request) error {
	query := `
    UPDATE archive SET status=? WHERE id=?;
  	`
	updated := a.database.Update(ctx, query, req.Status, req.ID)
	if !updated {
		return errors.New("failed to update archive")
	}
	return nil
}

func (a *ArchiveService) GetStatus(ctx context.Context, id string) (models.Request, error) {
	query := `
	select id, fileName, dir, status, createdOn
	from archive
	where id = ?
	`
	row := a.database.Select(ctx, query, id)
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

func compress(files []string, req *models.Request) error {
	zipPath := req.Dir + req.File
	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	_ = os.Remove(zipPath) // remove a single file
	file, err := os.OpenFile(zipPath, flags, 0644)
	if err != nil {
		return errors.New("failed to open zip for writing")
	}
	defer file.Close()

	zipw := zip.NewWriter(file)
	defer zipw.Close()

	for _, filename := range files {
		if len(filename) > 1 {
			if err := appendFiles(req.Dir, filename, zipw); err != nil {
				return err
			}
		}
	}
	return nil
}

func appendFiles(dir, filename string, zipw *zip.Writer) error {
	fileLoc := dir + "/" + filename
	file, err := os.Open(fileLoc)
	if err != nil {
		return fmt.Errorf("failed to open %s: %s", filename, err)
	}
	defer file.Close()

	wr, err := zipw.Create(filename)
	if err != nil {
		msg := "failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("failed to write %s to zip: %s", filename, err)
	}

	return nil
}

func getListOfFileNames(req *models.Request) ([]string, error) {
	// filter names
	var fileNames = make(map[string]string)
	if len(req.FilteredNames) > 1 {
		names := strings.Split(req.FilteredNames, "|")
		for _, n := range names {
			fileNames[n] = n
		}
	}

	limit := 10
	result := make([]string, 10)
	count := 1
	files, err := ioutil.ReadDir(req.Dir)
	if err != nil {
		return nil, err
	}

	if len(fileNames) > 0 {
		for _, file := range files {
			if len(file.Name()) > 0 {
				name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				namedFile := fileNames[name]
				if !file.IsDir() && count <= limit && len(namedFile) > 0 {
					result[count-1] = file.Name()
					count++
				}
			}
		}
	} else {
		sortFileSizeDescend(files)
		for _, file := range files {
			if !file.IsDir() && count <= limit && !strings.Contains(req.File, file.Name()) {
				result[count-1] = file.Name()
				count++
			}
		}
	}

	return result, nil
}

func sortFileSizeDescend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Size() > files[j].Size()
	})
}
