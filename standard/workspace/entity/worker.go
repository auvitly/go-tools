package entity

import (
	"cmp"
	"time"

	"github.com/google/uuid"
)

type Worker[T cmp.Ordered] struct {
	ID        uuid.UUID
	Type      T
	Version   string
	Labels    map[string]string
	CreatedAT time.Time
	UpdatedAT time.Time
}
