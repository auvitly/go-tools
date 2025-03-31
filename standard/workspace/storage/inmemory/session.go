package inmemory

import (
	"context"
	"slices"
	"sync"
	"time"

	"github.com/auvitly/go-tools/standard/workspace/entity"
	"github.com/auvitly/go-tools/standard/workspace/storage"
	"github.com/auvitly/go-tools/stderrs"
	"github.com/google/uuid"
)

type SessionStorage struct {
	mu      sync.RWMutex
	storage map[uuid.UUID]*entity.Session
	config  SessionConfig
}

type SessionConfig struct {
	DeleteCompleted bool
}

func NewSessionStorage(config SessionConfig) *SessionStorage {
	return &SessionStorage{
		storage: map[uuid.UUID]*entity.Session{},
		config:  config,
	}
}

func (s *SessionStorage) New(ctx context.Context, params storage.SessionNewParams) (*entity.Session, *stderrs.Error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var id = uuid.New()

	s.storage[id] = &entity.Session{
		ID:        id,
		WorkerID:  params.WorkerID,
		CreatedAT: time.Now(),
		UpdatedAT: time.Now(),
		DoneAT:    nil,
	}

	return s.storage[id], nil
}

func (s *SessionStorage) Get(ctx context.Context, params storage.SessionGetParams) (*entity.Session, *stderrs.Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.storage[params.SessionID]
	if !ok {
		return nil, stderrs.NotFound.SetMessage("not found session with id=%s", params.SessionID.String())
	}

	return session, nil
}

func (s *SessionStorage) List(ctx context.Context, params storage.SessionListParams) ([]*entity.Session, *stderrs.Error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var list []*entity.Session

	for _, session := range s.storage {
		var (
			condWorkerID      = params.WorkerID == nil || session.WorkerID == *params.WorkerID
			condSessionID     = len(params.SessionIDs) == 0 || slices.Contains(params.SessionIDs, session.ID)
			condOnlyCompleted = !params.OnlyCompleted || params.OnlyCompleted && session.DoneAT != nil
		)

		if condOnlyCompleted && condSessionID && condWorkerID {
			list = append(list, session)
		}
	}

	return list, nil
}

func (s *SessionStorage) Drop(ctx context.Context, params storage.SessionDropParams) *stderrs.Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.storage, params.SessionID)

	return nil
}

func (s *SessionStorage) Done(ctx context.Context, params storage.SessionDoneParams) *stderrs.Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, ok := s.storage[params.SessionID]
	if !ok {
		return stderrs.NotFound.SetMessage("not found session with id=%s", params.SessionID.String())
	}

	if s.config.DeleteCompleted {
		delete(s.storage, params.SessionID)
	} else {
		var ts = time.Now()

		session.DoneAT = &ts
	}

	return nil
}
