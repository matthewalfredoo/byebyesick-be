package repository

import (
	"context"
	"database/sql"
	"halodeksik-be/app/entity"
)

type ConsultationMessageRepository interface {
	Create(ctx context.Context, message entity.ConsultationMessage) (*entity.ConsultationMessage, error)
}

type ConsultationMessageRepositoryImpl struct {
	db *sql.DB
}

func NewConsultationMessageRepositoryImpl(db *sql.DB) *ConsultationMessageRepositoryImpl {
	return &ConsultationMessageRepositoryImpl{db: db}
}

func (repo *ConsultationMessageRepositoryImpl) Create(ctx context.Context, message entity.ConsultationMessage) (*entity.ConsultationMessage, error) {
	const create = `
	INSERT INTO consultation_messages (session_id, sender_id, message_type, message, attachment)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, session_id, sender_id, message_type, message, attachment, created_at, updated_at`

	row := repo.db.QueryRowContext(ctx, create, message.SessionId, message.SenderId, message.MessageType, message.Message, message.Attachment)
	var created entity.ConsultationMessage
	err := row.Scan(&created.Id, &created.SessionId, &created.SenderId, &created.MessageType, &created.Message, &created.Attachment, &created.CreatedAt, &created.UpdatedAt)

	return &created, err
}
