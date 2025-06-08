package api

import "github.com/google/uuid"

type CallbackResponse struct {
	RequestID uuid.UUID       `json:"RequestID"`
	Site      string          `json:"Site"`
	Results   map[string]bool `json:"Results"`
}
