package models

import (
	"errors"
	"strings"
	"time"
)

// Request struct
type Request struct {
	ID                string    `json:"id,omitempty"`
	File              string    `json:"file,omitempty"`
	Dir               string    `json:"dir,omitempty"`
	Status            string    `json:"status,omitempty"`
	FilteredNames     string    `json:"filteredNames,omitempty"`
	Aligorithm        string    `json:"aligorithm,omitempty"`
	PartialExtraction string    `json:"partialExtraction,omitempty"`
	Background        bool      `json:"background,omitempty"`
	CreatedOn         time.Time `json:"-"`
}

// Validate check if request is valid
func (r *Request) Validate(action string) error {
	switch strings.ToLower(action) {
	case "archive":
		if r.File == "" {
			return errors.New("file is required")
		}
		if r.Dir == "" {
			return errors.New("dir is required")
		}
		return nil
	case "extract":
		if r.File == "" {
			return errors.New("file is required")
		}
		if r.Dir == "" {
			return errors.New("dir is required")
		}
		return nil
	case "get":
		if r.ID == "" {
			return errors.New("id is required")
		}
		return nil
	default:
		return errors.New("invalid validation operation")
	}
}

// PrepareOutput initiliazes the response object
func (r *Request) PrepareOutput(request Request) {
	r.ID = request.ID
	r.Status = request.Status
}
