package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
	"time"
)

type ConsultationSessionUseCase interface {
	Add(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error)
	GetById(ctx context.Context, id int64) (*entity.ConsultationSession, error)
	GetAllByUserIdOrDoctorId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	EditTime(ctx context.Context, id int64) (*entity.ConsultationSession, error)
	EditStatusAsEnded(ctx context.Context, id int64) (*entity.ConsultationSession, error)
}

type ConsultationSessionUseCaseImpl struct {
	sessionRepo      repository.ConsultationSessionRepository
	prescriptionRepo repository.PrescriptionRepository
	sickLeaveRepo    repository.SickLeaveFormRepository
	userRepo         repository.UserRepository
}

func NewConsultationSessionUseCaseImpl(
	sessionRepo repository.ConsultationSessionRepository,
	prescriptionRepo repository.PrescriptionRepository,
	sickLeaveRepo repository.SickLeaveFormRepository,
	userRepo repository.UserRepository,
) *ConsultationSessionUseCaseImpl {
	return &ConsultationSessionUseCaseImpl{
		sessionRepo: sessionRepo, prescriptionRepo: prescriptionRepo, sickLeaveRepo: sickLeaveRepo, userRepo: userRepo,
	}
}

func (uc *ConsultationSessionUseCaseImpl) Add(ctx context.Context, session entity.ConsultationSession) (*entity.ConsultationSession, error) {
	doctor, err := uc.userRepo.FindById(ctx, session.DoctorId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(doctor, "Id", session.DoctorId)
		}
		return nil, err
	}
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	session.UserId = userId

	sessionDb, err := uc.sessionRepo.FindByUserIdAndDoctorId(ctx, session.UserId, session.DoctorId)
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, err
	}

	if !errors.Is(err, apperror.ErrRecordNotFound) && sessionDb.ConsultationSessionStatusId == appconstant.ConsultationSessionStatusOngoing {
		return sessionDb, apperror.ErrChatStillOngoing
	}

	session.ConsultationSessionStatusId = appconstant.ConsultationSessionStatusOngoing
	added, err := uc.sessionRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}
	return added, nil
}

func (uc *ConsultationSessionUseCaseImpl) GetById(ctx context.Context, id int64) (*entity.ConsultationSession, error) {
	sessionDb, err := uc.sessionRepo.FindByIdJoinAll(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(sessionDb, "Id", id)
		}
		return nil, err
	}

	clientIdCtx := ctx.Value(appconstant.ContextKeyUserId)
	clientId := clientIdCtx.(int64)

	roleIdCtx := ctx.Value(appconstant.ContextKeyRoleId)
	roleId := roleIdCtx.(int64)

	if roleId == appconstant.UserRoleIdDoctor {
		if sessionDb.DoctorId != clientId {
			return nil, apperror.ErrForbiddenViewEntity
		}
	}

	if roleId == appconstant.UserRoleIdUser {
		if sessionDb.UserId != clientId {
			return nil, apperror.ErrForbiddenViewEntity
		}
	}

	prescription, err := uc.prescriptionRepo.FindBySessionId(ctx, sessionDb.Id)
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, err
	}

	if prescription != nil {
		sessionDb.Prescription = prescription
	}

	sickLeave, err := uc.sickLeaveRepo.FindBySessionId(ctx, sessionDb.Id)
	if err != nil && !errors.Is(err, apperror.ErrRecordNotFound) {
		return nil, err
	}

	if sickLeave != nil {
		sessionDb.SickLeaveForm = sickLeave
	}

	return sessionDb, nil
}

func (uc *ConsultationSessionUseCaseImpl) GetAllByUserIdOrDoctorId(ctx context.Context, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userIdOrDoctorId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	sessions, err := uc.sessionRepo.FindAllByUserIdOrDoctorId(ctx, userIdOrDoctorId, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.sessionRepo.CountFindAllByUserIdOrDoctorId(ctx, userIdOrDoctorId, param)
	if err != nil {
		return nil, err
	}
	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(totalItems, totalPages, int64(len(sessions)), int64(*param.PageId), sessions)
	return paginatedItems, nil
}

func (uc *ConsultationSessionUseCaseImpl) getById(ctx context.Context, id int64) (*entity.ConsultationSession, error) {
	sessionDb, err := uc.sessionRepo.FindByIdJoinAll(ctx, id)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(sessionDb, "Id", id)
		}
		return nil, err
	}
	return sessionDb, nil
}

func (uc *ConsultationSessionUseCaseImpl) EditTime(ctx context.Context, id int64) (*entity.ConsultationSession, error) {
	sessionDb, err := uc.getById(ctx, id)
	if err != nil {
		return nil, err
	}

	if sessionDb.ConsultationSessionStatusId == appconstant.ConsultationSessionStatusEnded {
		return nil, apperror.ErrChatAlreadyEnded
	}

	sessionDb.UpdatedAt = time.Now()

	updated, err := uc.sessionRepo.Update(ctx, *sessionDb)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *ConsultationSessionUseCaseImpl) EditStatusAsEnded(ctx context.Context, id int64) (*entity.ConsultationSession, error) {
	sessionDb, err := uc.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if sessionDb.ConsultationSessionStatusId == appconstant.ConsultationSessionStatusEnded {
		return nil, apperror.ErrChatAlreadyEnded
	}

	sessionDb.ConsultationSessionStatusId = appconstant.ConsultationSessionStatusEnded

	updated, err := uc.sessionRepo.Update(ctx, *sessionDb)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
