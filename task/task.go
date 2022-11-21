package task

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/greatfocus/archive-service/database"
	"github.com/greatfocus/archive-service/models"
	"github.com/greatfocus/archive-service/services"
)

// Tasks struct
type Tasks struct {
	archiveService *services.ArchiveService
	extractService *services.ExtractService
	database       *database.Conn
}

// Init required parameters
func (t *Tasks) Init(db *database.Conn) {
	t.archiveService = &services.ArchiveService{}
	t.archiveService.Init(db)

	t.extractService = &services.ExtractService{}
	t.extractService.Init(db)

	t.database = db
}

// ExtractBackgroundFile start the job to extract files in the background
func (t *Tasks) ExtractBackgroundFile() {
	log.Println("Scheduler_ExtractBackgroundFile started")
	list, err := t.getBackgroundExtracts()
	if err != nil {
		log.Println("Scheduler_ExtractBackgroundFile failed to fetch archives")
		return
	}
	if len(list) > 0 {
		t.extractBulk(list)
	} else {
		log.Println("Scheduler_ExtractBackgroundFile queued is empty")
	}
	log.Println("Scheduler_ExtractBackgroundFile finished")
}

func (t *Tasks) extractBulk(list []models.Request) {
	for i := 0; i < len(list); i++ {

		go func(req *models.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
			defer cancel()
			t.extractService.InitiateExtract(ctx, req)
		}(&list[i])
	}
}

func (t *Tasks) getBackgroundExtracts() ([]models.Request, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
	defer cancel()
	query := `
	select id, fileName, dir, status, createdOn
	from extract
	where status = ? and background = ?
	LIMIT 10;
	`
	rows, err := t.database.Query(ctx, query, "new", true)
	if err != nil {
		return nil, err
	}
	result, err := archiveMapper(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ArchiveBackgroundFile start the job to archive files in the background
func (t *Tasks) ArchiveBackgroundFile() {
	log.Println("Scheduler_ArchiveBackgroundFile started")
	list, err := t.getBackgroundArchives()
	if err != nil {
		log.Println("Scheduler_ArchiveBackgroundFile failed to fetch archives")
		return
	}
	if len(list) > 0 {
		t.archiveBulk(list)
	} else {
		log.Println("Scheduler_ArchiveBackgroundFile queued is empty")
	}
	log.Println("Scheduler_ArchiveBackgroundFile finished")
}

func (t *Tasks) archiveBulk(list []models.Request) {
	for i := 0; i < len(list); i++ {

		go func(req *models.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
			defer cancel()
			t.archiveService.InitiateArchive(ctx, req)
		}(&list[i])
	}
}

func (t *Tasks) getBackgroundArchives() ([]models.Request, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
	defer cancel()
	query := `
	select id, fileName, dir, status, createdOn
	from archive
	where status = ? and background = ?
	LIMIT 10;
	`
	rows, err := t.database.Query(ctx, query, "new", true)
	if err != nil {
		return nil, err
	}
	result, err := archiveMapper(rows)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// prepare row
func archiveMapper(rows *sql.Rows) ([]models.Request, error) {
	requests := []models.Request{}
	for rows.Next() {
		var channel models.Request
		err := rows.Scan(&channel.ID, &channel.File, &channel.Dir, &channel.Status, &channel.CreatedOn)
		if err != nil {
			return nil, err
		}
		requests = append(requests, channel)
	}

	return requests, nil
}
