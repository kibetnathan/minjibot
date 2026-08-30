package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/kibetnathan/minjibot/infrastructure/postgres"
	"github.com/kibetnathan/minjibot/internal/domain/user"
	"github.com/kibetnathan/minjibot/internal/ports/dto"
)

// User Repository --
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (user.User, error)
	GetByDiscordID(ctx context.Context, discordID string) (user.User, error)
	GetByEmail(ctx context.Context, email string) (user.User, error)
	List(ctx context.Context) ([]user.User, error)
	Create(ctx context.Context, arg dto.CreateUserParams) (user.User, error)
	Update(ctx context.Context, id int64, arg dto.UpdateUserParams) (user.User, error)
	SetPassword(ctx context.Context, id int64, passwordhash string) (user.User, error)
	Deactivate(ctx context.Context, id int64) error
	Reactivate(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int64, error)
}

type sqlUserRepository struct {
	store *SQLStore
}

func NewUserRepository(store *SQLStore) UserRepository {
	return &sqlUserRepository{store: store}
}

func (r *sqlUserRepository) GetByID(ctx context.Context, id int64) (user.User, error) {
	u, err := r.store.queries.GetUser(ctx, id)
	if err != nil {
		return user.User{}, err
	}
	return toentitiesUser(u), nil
}

func (r *sqlUserRepository) GetByDiscordID(ctx context.Context, discordID string) (user.User, error) {
	u, err := r.store.queries.GetUserByDiscordID(ctx, discordID)
	if err != nil {
		return user.User{}, err
	}
	return toentitiesUser(u), nil
}

func (r *sqlUserRepository) GetByEmail(ctx context.Context, email string) (user.User, error) {
	u, err := r.store.queries.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: email != ""})
	if err != nil {
		return user.User{}, err
	}
	return toentitiesUser(u), nil
}

func (r *sqlUserRepository) List(ctx context.Context) ([]user.User, error) {
	users, err := r.store.queries.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return toentitiesUsers(users), nil
}

func (r *sqlUserRepository) Create(ctx context.Context, arg dto.CreateUserParams) (user.User, error) {
	u, err := r.store.queries.CreateUser(ctx, postgres.CreateUserParams{
		UserID: arg.UserID,
		Email: pgtype.Text{
			String: strOrEmpty(arg.Email),
			Valid:  arg.Email != nil,
		},
		Passwordhash: pgtype.Text{
			String: strOrEmpty(arg.Passwordhash),
			Valid:  arg.Passwordhash != nil,
		},
	})
	if err != nil {
		return user.User{}, err
	}
	return toentitiesUser(u), nil
}

func (r *sqlUserRepository) Update(ctx context.Context, id int64, arg dto.UpdateUserParams) (user.User, error) {
	u, err := r.store.queries.UpdateUser(ctx, postgres.UpdateUserParams{
		ID: id,
		Email: pgtype.Text{
			String: strOrEmpty(arg.Email),
			Valid:  arg.Email != nil,
		},
		Passwordhash: pgtype.Text{
			String: strOrEmpty(arg.Passwordhash),
			Valid:  arg.Passwordhash != nil,
		},
		IsActive: arg.IsActive,
	})
	if err != nil {
		return user.User{}, err
	}
	return toentitiesUser(u), nil
}

func (r *sqlUserRepository) SetPassword(ctx context.Context, id int64, passwordhash string) (user.User, error) {
	u, err := r.store.queries.SetUserPassword(ctx, postgres.SetUserPasswordParams{
		ID:           id,
		Passwordhash: pgtype.Text{String: passwordhash, Valid: passwordhash != ""},
	})
	if err != nil {
		return user.User{}, err
	}
	return toentitiesUser(u), nil
}

func (r *sqlUserRepository) Deactivate(ctx context.Context, id int64) error {
	return r.store.queries.DeactivateUser(ctx, id)
}

func (r *sqlUserRepository) Reactivate(ctx context.Context, id int64) error {
	return r.store.queries.ReactivateUser(ctx, id)
}

func (r *sqlUserRepository) Delete(ctx context.Context, id int64) error {
	return r.store.queries.DeleteUser(ctx, id)
}

func (r *sqlUserRepository) Count(ctx context.Context) (int64, error) {
	return r.store.queries.CountUsers(ctx)
}

func toentitiesUsers(users []postgres.User) []user.User {
	entitiesUsers := make([]user.User, len(users))
	for i, u := range users {
		entitiesUsers[i] = toentitiesUser(u)
	}
	return entitiesUsers
}

func toentitiesUser(u postgres.User) user.User {
	return user.User{
		ID:           u.ID,
		UserID:       u.UserID,
		Email:        strOrEmptyPtr(u.Email),
		Passwordhash: strOrEmptyPtr(u.Passwordhash),
		IsActive:     u.IsActive,
		CreatedAt:    u.CreatedAt.Time,
		UpdatedAt:    u.UpdatedAt.Time,
	}
}

func strOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func strOrEmptyPtr(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}
