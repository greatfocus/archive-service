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

// Extract struct
type Extract struct {
	extractService *services.ExtractService
}

// ServeHTTP checks if is valid method
func (f Extract) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		f.getStatus(w, r)
		return
	}
	if r.Method == http.MethodPost {
		f.createExtract(w, r)
		return
	}

	// catch all
	// if no method is satisfied return an error
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Header().Add("Allow", "GET, POST")
}

// Init method
func (f *Extract) Init(ExtractService *services.ExtractService) {
	f.extractService = ExtractService
}

// create prepares Extract
func (f *Extract) createExtract(w http.ResponseWriter, r *http.Request) {
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

	err = req.Validate("extract")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		Error(w, r, err)
		return
	}

	res, err := f.extractService.CreateExtract(ctx, &req)
	if err != nil {
		log.Printf("%s", err)
		Error(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	Success(w, r, res)
}

// getExtracts method
func (f *Extract) getStatus(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Minute))
	defer cancel()

	id := r.FormValue("id")
	if id != "" {
		Extract, err := f.extractService.GetStatus(ctx, id)
		if err != nil {
			log.Printf("%s", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			Error(w, r, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		Success(w, r, Extract)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	Error(w, r, errors.New("invalid payload request"))
}
