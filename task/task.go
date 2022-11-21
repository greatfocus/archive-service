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
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
	// defer cancel()

}

// ArchiveBackgroundFile start the job to archive files in the background
func (t *Tasks) ArchiveBackgroundFile() {
	log.Println("Scheduler_ArchiveBackgroundFile started")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Minute)
	defer cancel()

	list, err := t.getBackgroundArchives(ctx)
	if err != nil {
		log.Println("Scheduler_ArchiveBackgroundFile failed to fetch archives")
		return
	}
	if len(list) > 0 {
		t.archiveBulk(ctx, list)
	} else {
		log.Println("Scheduler_ArchiveBackgroundFile Email queued is empty")
	}
	log.Println("Scheduler_ArchiveBackgroundFile finished")
}

func (t *Tasks) archiveBulk(ctx context.Context, list []models.Request) {
	for i := 0; i < len(list); i++ {
		go t.archiveService.InitiateArchive(ctx, &list[i])
	}
}

func (t *Tasks) getBackgroundArchives(ctx context.Context) ([]models.Request, error) {
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
