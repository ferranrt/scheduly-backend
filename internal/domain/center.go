package domain

import "github.com/google/uuid"

type Center struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id"`
}
