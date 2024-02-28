package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type ForgotTokenRepository interface {
	CreateForgotToken(ctx context.Context, token entity.ForgotPasswordToken) (*entity.ForgotPasswordToken, error)
	FindForgotTokenByUserId(ctx context.Context, userId int64) (*entity.ForgotPasswordToken, error)
	FindForgotTokenByToken(ctx context.Context, token string) (*entity.ForgotPasswordToken, error)
	DeactivateForgotToken(ctx context.Context, token entity.ForgotPasswordToken) (*entity.ForgotPasswordToken, error)
}

type ForgotTokenRepositoryImpl struct {
	db *sql.DB
}

func NewForgotTokenRepository(db *sql.DB) *ForgotTokenRepositoryImpl {
	repo := ForgotTokenRepositoryImpl{db: db}
	return &repo
}

func (repo *ForgotTokenRepositoryImpl) CreateForgotToken(ctx context.Context, token entity.ForgotPasswordToken) (*entity.ForgotPasswordToken, error) {
	const createForgotToken = `
	INSERT INTO forgot_password_tokens(token, is_valid, expired_at, user_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id, token, is_valid, expired_at, user_id, created_at, updated_at, deleted_at
	`

	row := repo.db.QueryRowContext(ctx, createForgotToken,
		token.Token,
		token.IsValid,
		token.ExpiredAt,
		token.UserId,
	)

	var createdToken entity.ForgotPasswordToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.UserId,
		&createdToken.CreatedAt,
		&createdToken.UpdatedAt,
		&createdToken.DeletedAt,
	)

	return &createdToken, err
}

func (repo *ForgotTokenRepositoryImpl) FindForgotTokenByUserId(ctx context.Context, userId int64) (*entity.ForgotPasswordToken, error) {
	const getForgotTokenByUserId = `
	SELECT id, token, is_valid, expired_at, user_id, created_at, updated_at, deleted_at FROM forgot_password_tokens
	WHERE user_id = $1
	`

	row := repo.db.QueryRowContext(ctx, getForgotTokenByUserId,
		userId,
	)

	var createdToken entity.ForgotPasswordToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.UserId,
		&createdToken.CreatedAt,
		&createdToken.UpdatedAt,
		&createdToken.DeletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return &createdToken, err
}

func (repo *ForgotTokenRepositoryImpl) FindForgotTokenByToken(ctx context.Context, token string) (*entity.ForgotPasswordToken, error) {
	const getForgotTokenByToken = `-- name: GetForgotTokenByToken :one
	SELECT id, token, is_valid, expired_at, user_id, created_at, updated_at, deleted_at FROM forgot_password_tokens
	WHERE token = $1
	`

	row := repo.db.QueryRowContext(ctx, getForgotTokenByToken,
		token,
	)

	var createdToken entity.ForgotPasswordToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.UserId,
		&createdToken.CreatedAt,
		&createdToken.UpdatedAt,
		&createdToken.DeletedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}

	if err != nil {
		return nil, err
	}

	return &createdToken, err

}

func (repo *ForgotTokenRepositoryImpl) DeactivateForgotToken(ctx context.Context, token entity.ForgotPasswordToken) (*entity.ForgotPasswordToken, error) {
	const deleteToken = `
	DELETE FROM forgot_password_tokens WHERE token = $1
	`

	row := repo.db.QueryRowContext(ctx, deleteToken,
		token.Token,
	)

	return &token, row.Err()
}
