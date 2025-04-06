package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	WorkerID  uuid.UUID
	CreatedAT time.Time
	UpdatedAT time.Time
	DoneAT    *time.Time
}

type IsSession interface{ Impl() *Session }

func (s *Session) Impl() *Session { return s }
