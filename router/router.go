package router

import (
	"log"
	"net/http"

	"github.com/greatfocus/archive-service/database"
	"github.com/greatfocus/archive-service/handler"
	"github.com/greatfocus/archive-service/services"
)

// LoadRouter creates service handlers in MUx
func LoadRouter(db *database.Conn) *http.ServeMux {
	mux := http.NewServeMux()
	createHanlders(db, mux)
	log.Println("Created routes with handler")
	return mux
}

// createHanlders prepares handlers with services requires
func createHanlders(db *database.Conn, mux *http.ServeMux) {

	archiveService := services.ArchiveService{}
	archiveService.Init(db)

	archiveHandler := handler.Archive{}
	archiveHandler.Init(&archiveService)
	mux.Handle("/archive", archiveHandler)

	extractService := services.ExtractService{}
	extractService.Init(db)
	extractHandler := handler.Extract{}
	extractHandler.Init(&extractService)
	mux.Handle("/extract", extractHandler)
}
