package services

import (
	"context"
	"database/sql"
	"errors"

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
		return a.archiveFiles(ctx, req)
	}

}

func (a *ArchiveService) archiveFiles(ctx context.Context, req *models.Request) (*models.Request, error) {

	return req, nil
}

func (a *ArchiveService) insertRecordToDB(ctx context.Context, req *models.Request) (*models.Request, error) {
	req.ID = uuid.New().String()
	req.Status = "new"
	query := `
	insert into archive (id, fileName, dir, status, aligorithm, filters, background)
	archive VALUES(?,?,?,?,?,?,?)
	returning id
	`
	_, inserted := a.database.Insert(ctx, query, req.ID, req.File, req.Dir, req.Status, req.Aligorithm, req.Filters, req.Background)
	if !inserted {
		return req, errors.New("failed to insert archive")
	}
	return req, nil
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
		return result, err
	case nil:
		// update cache
		return result, nil
	default:
		return result, err
	}
}
