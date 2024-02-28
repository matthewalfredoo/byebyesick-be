package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type PrescriptionUseCase interface {
	Add(ctx context.Context, prescription entity.Prescription) (*entity.Prescription, error)
	GetBySessionId(ctx context.Context, sessionId int64) (*entity.Prescription, error)
	EditBySessionId(ctx context.Context, sessionId int64, prescription entity.Prescription) (*entity.Prescription, error)
}

type PrescriptionUseCaseImpl struct {
	prescriptionRepo repository.PrescriptionRepository
	sessionRepo      repository.ConsultationSessionRepository
}

func NewPrescriptionUseCaseImpl(prescriptionRepo repository.PrescriptionRepository, sessionRepo repository.ConsultationSessionRepository) *PrescriptionUseCaseImpl {
	return &PrescriptionUseCaseImpl{prescriptionRepo: prescriptionRepo, sessionRepo: sessionRepo}
}

func (uc *PrescriptionUseCaseImpl) Add(ctx context.Context, prescription entity.Prescription) (*entity.Prescription, error) {
	sessionDb, err := uc.sessionRepo.FindById(ctx, prescription.SessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(sessionDb, "Id", prescription.SessionId)
		}
		return nil, err
	}

	doctorId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	if sessionDb.DoctorId != doctorId {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	added, err := uc.prescriptionRepo.Create(ctx, prescription)
	if err != nil {
		return nil, err
	}
	return uc.GetBySessionId(ctx, added.SessionId)
}

func (uc *PrescriptionUseCaseImpl) GetBySessionId(ctx context.Context, sessionId int64) (*entity.Prescription, error) {
	sessionDb, err := uc.sessionRepo.FindById(ctx, sessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(&entity.Prescription{}, "SessionId", sessionId)
		}
		return nil, err
	}

	roleId := ctx.Value(appconstant.ContextKeyRoleId).(int64)
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	if roleId == appconstant.UserRoleIdDoctor && sessionDb.DoctorId != userId {
		return nil, apperror.ErrForbiddenViewEntity
	}

	if roleId == appconstant.UserRoleIdUser && sessionDb.UserId != userId {
		return nil, apperror.ErrForbiddenViewEntity
	}

	prescription, err := uc.prescriptionRepo.FindBySessionIdDetailed(ctx, sessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(prescription, "SessionId", sessionId)
		}
		return nil, err
	}
	return prescription, nil
}

func (uc *PrescriptionUseCaseImpl) EditBySessionId(ctx context.Context, sessionId int64, prescription entity.Prescription) (*entity.Prescription, error) {
	sessionDb, err := uc.sessionRepo.FindById(ctx, sessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(&entity.Prescription{}, "SessionId", sessionId)
		}
		return nil, err
	}

	doctorId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	if sessionDb.DoctorId != doctorId {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	prescriptionDb, err := uc.prescriptionRepo.FindBySessionId(ctx, sessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(prescriptionDb, "SessionId", sessionId)
		}
		return nil, err
	}

	prescription.SessionId = sessionId

	edited, err := uc.prescriptionRepo.UpdateBySessionId(ctx, prescription)
	if err != nil {
		return nil, err
	}
	return uc.GetBySessionId(ctx, edited.SessionId)
}
