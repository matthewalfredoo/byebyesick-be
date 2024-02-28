package usecase

import (
	"context"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ConsultationMessageUseCase interface {
	Add(ctx context.Context, message entity.ConsultationMessage) (*entity.ConsultationMessage, error)
}

type ConsultationMessageUseCaseImpl struct {
	repo repository.ConsultationMessageRepository
}

func NewConsultationMessageUseCaseImpl(repo repository.ConsultationMessageRepository) *ConsultationMessageUseCaseImpl {
	return &ConsultationMessageUseCaseImpl{repo: repo}
}

func (uc *ConsultationMessageUseCaseImpl) Add(ctx context.Context, message entity.ConsultationMessage) (*entity.ConsultationMessage, error) {
	added, err := uc.repo.Create(ctx, message)
	if err != nil {
		return nil, err
	}
	return added, nil
}
