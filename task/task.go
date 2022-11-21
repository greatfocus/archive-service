package task

import (
	"github.com/greatfocus/archive-service/database"
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
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Minute)
	// defer cancel()

}
