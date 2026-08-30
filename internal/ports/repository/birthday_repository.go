package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/birthday"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// Birthday Repository --
type BirthdayRepository interface {
	Set(ctx context.Context, arg dto.SetBirthdayParams) (birthday.Birthday, error)
	Get(ctx context.Context, guildID, userID string) (birthday.Birthday, error)
	ListByGuild(ctx context.Context, guildID string) ([]birthday.Birthday, error)
	ListByMonthDay(ctx context.Context, guildID string, month, day int32) ([]birthday.Birthday, error)
	Delete(ctx context.Context, guildID, userID string) error
}

type sqlBirthdayRepository struct {
	store *SQLStore
}

func NewBirthdayRepository(store *SQLStore) BirthdayRepository {
	return &sqlBirthdayRepository{store: store}
}

func (r *sqlBirthdayRepository) Set(ctx context.Context, arg dto.SetBirthdayParams) (birthday.Birthday, error) {
	b, err := r.store.queries.UpsertBirthday(ctx, postgres.UpsertBirthdayParams{
		GuildID:  arg.GuildID,
		UserID:   arg.UserID,
		Birthday: pgtype.Date{Time: arg.Birthday, Valid: true},
	})
	if err != nil {
		return birthday.Birthday{}, err
	}
	return toentitiesBirthday(b), nil
}

func (r *sqlBirthdayRepository) Get(ctx context.Context, guildID, userID string) (birthday.Birthday, error) {
	b, err := r.store.queries.GetBirthday(ctx, postgres.GetBirthdayParams{
		GuildID: guildID,
		UserID:  userID,
	})
	if err != nil {
		return birthday.Birthday{}, err
	}
	return toentitiesBirthday(b), nil
}

func (r *sqlBirthdayRepository) ListByGuild(ctx context.Context, guildID string) ([]birthday.Birthday, error) {
	birthdays_, err := r.store.queries.ListBirthdaysByGuild(ctx, guildID)
	if err != nil {
		return nil, err
	}
	return toentitiesBirthdays(birthdays_), nil
}

func (r *sqlBirthdayRepository) ListByMonthDay(ctx context.Context, guildID string, month, day int32) ([]birthday.Birthday, error) {
	birthdays_, err := r.store.queries.ListBirthdaysTodayByGuild(ctx, postgres.ListBirthdaysTodayByGuildParams{
		GuildID: guildID,
		Month:   month,
		Day:     day,
	})
	if err != nil {
		return nil, err
	}
	return toentitiesBirthdays(birthdays_), nil
}

func (r *sqlBirthdayRepository) Delete(ctx context.Context, guildID, userID string) error {
	return r.store.queries.DeleteBirthday(ctx, postgres.DeleteBirthdayParams{
		GuildID: guildID,
		UserID:  userID,
	})
}

func toentitiesBirthdays(birthdays_ []postgres.Birthday) []birthday.Birthday {
	entitiesBirthdays := make([]birthday.Birthday, len(birthdays_))
	for i, b := range birthdays_ {
		entitiesBirthdays[i] = toentitiesBirthday(b)
	}
	return entitiesBirthdays
}

func toentitiesBirthday(b postgres.Birthday) birthday.Birthday {
	return birthday.Birthday{
		ID:        b.ID,
		GuildID:   b.GuildID,
		UserID:    b.UserID,
		Birthday:  b.Birthday.Time,
		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
	}
}
