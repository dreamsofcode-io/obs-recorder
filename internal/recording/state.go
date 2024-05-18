package recording

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrNotFound      = errors.New("No recording for the given ID")
	ErrAlreadyExists = errors.New("Recording already exists for ID")
)

type Recording struct {
	ID   uuid.UUID `json:"id"`
	Path string    `json:"-"`
}

type Repository struct {
	lock       sync.RWMutex
	recordings map[uuid.UUID]Recording
}

func NewRepository() *Repository {
	return &Repository{
		recordings: make(map[uuid.UUID]Recording),
	}
}

func (r *Repository) Insert(ctx context.Context, recording Recording) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, exists := r.recordings[recording.ID]
	if exists {
		return ErrAlreadyExists
	}
	r.recordings[recording.ID] = recording

	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (Recording, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	rec, exists := r.recordings[id]
	if !exists {
		return Recording{}, ErrNotFound
	}

	return rec, nil
}

func (r *Repository) Update(ctx context.Context, recording Recording) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, exists := r.recordings[recording.ID]
	if !exists {
		return ErrNotFound
	}

	r.recordings[recording.ID] = recording

	return nil
}
