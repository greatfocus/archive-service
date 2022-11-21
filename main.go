package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/greatfocus/archive-service/database"
	"github.com/greatfocus/archive-service/router"
	"github.com/greatfocus/archive-service/task"
	gfcron "github.com/greatfocus/gf-cron"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	// initialize services
	var db = database.Conn{}
	db.Connect()

	// background task
	tasks := task.Tasks{}
	tasks.Init(&db)
	gfcron.New().MustAddJob("0 0 * * *", tasks.ExtractBackgroundFile)
	gfcron.New().MustAddJob("0 0 * * *", tasks.ArchiveBackgroundFile)

	mux := router.LoadRouter(&db)
	serve(mux)
}

// serve creates server instance
func serve(mux *http.ServeMux) {
	timeout, err := strconv.ParseUint(os.Getenv("SERVER_TIMEOUT"), 0, 64)
	if err != nil {
		log.Fatal(fmt.Println(err))
	}

	addr := ":" + os.Getenv("SERVER_PORT")
	srv := &http.Server{
		Addr:           addr,
		ReadTimeout:    time.Duration(timeout) * time.Second,
		WriteTimeout:   time.Duration(timeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}
	log.Println("Listening to port HTTP", addr)
	log.Fatal(srv.ListenAndServe())
}
