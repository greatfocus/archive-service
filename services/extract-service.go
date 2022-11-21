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
	req.ID = uuid.New().String()
	file := req.Dir + "/" + req.File
	dst := req.Dir

	archive, err := zip.OpenReader(file)
	if err != nil {
		return req, errors.New("files does not exist")
	}
	defer archive.Close()

	for _, f := range archive.File {
		var filePath = filepath.Join(dst, f.Name)
		log.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return req, errors.New("invalid file path")
		}
		if f.FileInfo().IsDir() {
			log.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return req, err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return req, err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return req, err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return req, err
		}

		dstFile.Close()
		fileInArchive.Close()
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
		return result, err
	case nil:
		// update cache
		return result, nil
	default:
		return result, err
	}
}
