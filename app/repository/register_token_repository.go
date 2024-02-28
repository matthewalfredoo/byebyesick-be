package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type RegisterTokenRepository interface {
	CreateRegisterToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error)
	FindRegisterTokenByToken(ctx context.Context, token string) (*entity.VerificationToken, error)
	FindRegisterTokenByEmail(ctx context.Context, email string) (*entity.VerificationToken, error)
	DeactivateRegisterToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error)
}

type RegisterTokenRepositoryImpl struct {
	db *sql.DB
}

func NewRegisterTokenRepository(db *sql.DB) *RegisterTokenRepositoryImpl {
	repo := RegisterTokenRepositoryImpl{db: db}
	return &repo
}

func (repo *RegisterTokenRepositoryImpl) DeactivateRegisterToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error) {
	const deleteToken = `
	DELETE FROM verification_tokens WHERE token = $1
	`

	row := repo.db.QueryRowContext(ctx, deleteToken,
		token.Token,
	)

	return &token, row.Err()
}

func (repo *RegisterTokenRepositoryImpl) FindRegisterTokenByEmail(ctx context.Context, email string) (*entity.VerificationToken, error) {
	const getActiveVerifyTokenByEmail = `
	SELECT id, token, is_valid, expired_at, email, created_at, updated_at, deleted_at FROM verification_tokens
	WHERE email = $1
	`

	row := repo.db.QueryRowContext(ctx, getActiveVerifyTokenByEmail,
		email,
	)

	var createdToken entity.VerificationToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.Email,
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

func (repo *RegisterTokenRepositoryImpl) FindRegisterTokenByToken(ctx context.Context, token string) (*entity.VerificationToken, error) {
	const getTokenByToken = `
	SELECT id, token, is_valid, expired_at, email, created_at, updated_at, deleted_at FROM verification_tokens
	WHERE token = $1
	`

	row := repo.db.QueryRowContext(ctx, getTokenByToken,
		token,
	)

	var createdToken entity.VerificationToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.Email,
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

func (repo *RegisterTokenRepositoryImpl) CreateRegisterToken(ctx context.Context, token entity.VerificationToken) (*entity.VerificationToken, error) {

	const createVerifyToken = `
	INSERT INTO verification_tokens(token, is_valid, expired_at, email)
	VALUES ($1, $2, $3, $4)
	RETURNING id, token, is_valid, expired_at, email, created_at, updated_at, deleted_at
	`
	row := repo.db.QueryRowContext(ctx, createVerifyToken,
		token.Token,
		token.IsValid,
		token.ExpiredAt,
		token.Email,
	)

	var createdToken entity.VerificationToken

	err := row.Scan(
		&createdToken.Id,
		&createdToken.Token,
		&createdToken.IsValid,
		&createdToken.ExpiredAt,
		&createdToken.Email,
		&createdToken.CreatedAt,
		&createdToken.UpdatedAt,
		&createdToken.DeletedAt,
	)

	return &createdToken, err

}
