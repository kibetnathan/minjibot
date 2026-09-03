package repository

import (
	"context"

	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/diary"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// DiaryRepository --
type DiaryRepository interface {
	Create(ctx context.Context, arg dto.CreateDiaryEntryParams) (diary.Entry, error)
	ListByUser(ctx context.Context, userID string) ([]diary.Entry, error)
	Get(ctx context.Context, id int64, userID string) (diary.Entry, error)
	Delete(ctx context.Context, id int64, userID string) error
}

type sqlDiaryRepository struct {
	store *SQLStore
}

func NewDiaryRepository(store *SQLStore) DiaryRepository {
	return &sqlDiaryRepository{store: store}
}

func (r *sqlDiaryRepository) Create(ctx context.Context, arg dto.CreateDiaryEntryParams) (diary.Entry, error) {
	e, err := r.store.queries.CreateDiaryEntry(ctx, postgres.CreateDiaryEntryParams{
		UserID:  arg.UserID,
		Content: arg.Content,
	})
	if err != nil {
		return diary.Entry{}, err
	}
	return toentitiesDiaryEntry(e), nil
}

func (r *sqlDiaryRepository) ListByUser(ctx context.Context, userID string) ([]diary.Entry, error) {
	entries_, err := r.store.queries.ListDiaryEntriesByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	entries := make([]diary.Entry, len(entries_))
	for i, e := range entries_ {
		entries[i] = toentitiesDiaryEntry(e)
	}
	return entries, nil
}

func (r *sqlDiaryRepository) Get(ctx context.Context, id int64, userID string) (diary.Entry, error) {
	e, err := r.store.queries.GetDiaryEntry(ctx, postgres.GetDiaryEntryParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		return diary.Entry{}, err
	}
	return toentitiesDiaryEntry(e), nil
}

func (r *sqlDiaryRepository) Delete(ctx context.Context, id int64, userID string) error {
	return r.store.queries.DeleteDiaryEntry(ctx, postgres.DeleteDiaryEntryParams{
		ID:     id,
		UserID: userID,
	})
}

func toentitiesDiaryEntry(e postgres.DiaryEntry) diary.Entry {
	return diary.Entry{
		ID:        e.ID,
		UserID:    e.UserID,
		Content:   e.Content,
		CreatedAt: e.CreatedAt.Time,
	}
}
