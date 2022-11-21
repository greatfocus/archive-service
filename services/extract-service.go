package services

import (
	"context"
	"database/sql"
	"errors"

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
		return e.archiveFiles(ctx, req)
	}

}

func (e *ExtractService) archiveFiles(ctx context.Context, req *models.Request) (*models.Request, error) {

	return req, nil
}

func (e *ExtractService) insertRecordToDB(ctx context.Context, req *models.Request) (*models.Request, error) {
	req.ID = uuid.New().String()
	req.Status = "new"
	query := `
	insert into extract (id, fileName, dir, status, aligorithm, filters, partialExtraction, background)
	archive VALUES(?,?,?,?,?,?,?,?)
	returning id
	`
	_, inserted := e.database.Insert(ctx, query, req.ID, req.File, req.Dir, req.Status, req.Aligorithm, req.Filters, req.PartialExtraction, req.Background)
	if !inserted {
		return req, errors.New("failed to insert archive")
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
