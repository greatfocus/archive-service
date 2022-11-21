package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/greatfocus/archive-service/models"
	"github.com/greatfocus/archive-service/services"
)

// Archive struct
type Archive struct {
	archiveHandler func(http.ResponseWriter, *http.Request)
	archiveService *services.ArchiveService
}

// ServeHTTP checks if is valid method
func (a Archive) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		a.getStatus(w, r)
		return
	}
	if r.Method == http.MethodPost {
		a.createArchive(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET, POST")
}

// Init method
func (a *Archive) Init(ArchiveService *services.ArchiveService) {
	a.archiveService = ArchiveService
}

// create prepares Archive
func (a *Archive) createArchive(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	req := models.Request{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		derr := errors.New("invalid payload request")
		w.WriteHeader(http.StatusBadRequest)
		Error(w, r, derr)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		derr := errors.New("invalid payload request")
		w.WriteHeader(http.StatusBadRequest)
		Error(w, r, derr)
		return
	}

	err = req.Validate("archive")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		Error(w, r, err)
		return
	}

	res, err := a.archiveService.CreateArchive(ctx, &req)
	if err != nil {
		log.Printf("%s", err)
		Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	Success(w, r, res)
}

// getArchives method
func (a *Archive) getStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	id := r.FormValue("id")
	if id != "" {
		Archive, err := a.archiveService.GetStatus(ctx, id)
		if err != nil {
			log.Printf("%s", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		Success(w, r, Archive)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	Error(w, r, errors.New("invalid payload request"))
	return
}

// Success returns object as json
func Success(w http.ResponseWriter, r *http.Request, data interface{}) {
	if data != nil {
		response(w, r, data, "success")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	response(w, r, nil, "success")
}

// Error returns error as json
func Error(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		response(w, r, struct {
			Error string `json:"error"`
		}{Error: err.Error()}, "error")
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	response(w, r, nil, "error")
}

// response returns payload
func response(w http.ResponseWriter, r *http.Request, data interface{}, message string) {
	out, _ := json.Marshal(data)
	res := models.Response{
		Result: string(out),
	}
	_ = json.NewEncoder(w).Encode(res)
}
