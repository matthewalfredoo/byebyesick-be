package usecase

import (
	"context"
	"errors"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type SickLeaveFormUseCase interface {
	Add(ctx context.Context, form entity.SickLeaveForm) (*entity.SickLeaveForm, error)
	GetBySessionId(ctx context.Context, sessionId int64) (*entity.SickLeaveForm, error)
	EditBySessionId(ctx context.Context, sessionId int64, form entity.SickLeaveForm) (*entity.SickLeaveForm, error)
}

type SickLeaveFormUseCaseImpl struct {
	formRepo         repository.SickLeaveFormRepository
	sessionRepo      repository.ConsultationSessionRepository
	prescriptionRepo repository.PrescriptionRepository
	messageRepo      repository.ConsultationMessageRepository
}

func NewSickLeaveFormUseCaseImpl(
	formRepo repository.SickLeaveFormRepository,
	sessionRepo repository.ConsultationSessionRepository,
	prescriptionRepo repository.PrescriptionRepository,
	messageRepo repository.ConsultationMessageRepository,
) *SickLeaveFormUseCaseImpl {
	return &SickLeaveFormUseCaseImpl{
		formRepo: formRepo, sessionRepo: sessionRepo, prescriptionRepo: prescriptionRepo, messageRepo: messageRepo,
	}
}

func (uc *SickLeaveFormUseCaseImpl) Add(ctx context.Context, form entity.SickLeaveForm) (*entity.SickLeaveForm, error) {
	session, err := uc.sessionRepo.FindById(ctx, form.SessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(session, "Id", form.SessionId)
		}
		return nil, err
	}

	doctorId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	if session.DoctorId != doctorId {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	_, err = uc.prescriptionRepo.FindBySessionId(ctx, form.SessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.ErrConsultationSessionPrescriptionMustExistBeforeIssuingSickLeave
		}
		return nil, err
	}

	added, err := uc.formRepo.Create(ctx, form)
	if err != nil {
		return nil, err
	}

	if session.ConsultationSessionStatusId == appconstant.ConsultationSessionStatusEnded {
		userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
		roleId := ctx.Value(appconstant.ContextKeyRoleId).(int64)

		if roleId != appconstant.UserRoleIdDoctor {
			return added, nil
		}

		message := entity.ConsultationMessage{
			SessionId:   appdb.NewSqlNullInt64(session.Id),
			SenderId:    appdb.NewSqlNullInt64(userId),
			MessageType: appdb.NewSqlNullInt64(appconstant.MessageTypeAlert),
			Message:     appdb.NewSqlNullString(appconstant.MessageDoctorCreateLeaveSick),
			Attachment:  appdb.NewSqlNullString(""),
		}
		_, _ = uc.messageRepo.Create(ctx, message)
	}

	return added, nil
}

func (uc *SickLeaveFormUseCaseImpl) GetBySessionId(ctx context.Context, sessionId int64) (*entity.SickLeaveForm, error) {
	session, err := uc.sessionRepo.FindById(ctx, sessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(&entity.SickLeaveForm{}, "SessionId", sessionId)
		}
		return nil, err
	}

	clientId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	roleId := ctx.Value(appconstant.ContextKeyRoleId).(int64)

	if roleId == appconstant.UserRoleIdDoctor && session.DoctorId != clientId {
		return nil, apperror.ErrForbiddenViewEntity
	}

	if roleId == appconstant.UserRoleIdUser && session.UserId != clientId {
		return nil, apperror.ErrForbiddenViewEntity
	}

	form, err := uc.formRepo.FindBySessionIdDetailed(ctx, sessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(form, "SessionId", sessionId)
		}
		return nil, err
	}

	return form, nil
}

func (uc *SickLeaveFormUseCaseImpl) EditBySessionId(ctx context.Context, sessionId int64, form entity.SickLeaveForm) (*entity.SickLeaveForm, error) {
	session, err := uc.sessionRepo.FindById(ctx, sessionId)
	if err != nil {
		if errors.Is(err, apperror.ErrRecordNotFound) {
			return nil, apperror.NewNotFound(&entity.SickLeaveForm{}, "SessionId", sessionId)
		}
		return nil, err
	}

	doctorId := ctx.Value(appconstant.ContextKeyUserId).(int64)
	if session.DoctorId != doctorId {
		return nil, apperror.ErrForbiddenModifyEntity
	}

	formDb, err := uc.GetBySessionId(ctx, sessionId)
	if err != nil {
		return nil, err
	}

	formDb.SessionId = sessionId
	formDb.StartingDate = form.StartingDate
	formDb.EndingDate = form.EndingDate
	formDb.Description = form.Description

	edited, err := uc.formRepo.UpdateBySessionId(ctx, *formDb)
	if err != nil {
		return nil, err
	}

	if session.ConsultationSessionStatusId == appconstant.ConsultationSessionStatusEnded {
		userId := ctx.Value(appconstant.ContextKeyUserId).(int64)
		roleId := ctx.Value(appconstant.ContextKeyRoleId).(int64)

		if roleId != appconstant.UserRoleIdDoctor {
			return edited, nil
		}

		message := entity.ConsultationMessage{
			SessionId:   appdb.NewSqlNullInt64(formDb.SessionId),
			SenderId:    appdb.NewSqlNullInt64(userId),
			MessageType: appdb.NewSqlNullInt64(appconstant.MessageTypeAlert),
			Message:     appdb.NewSqlNullString(appconstant.MessageDoctorUpdateLeaveSick),
			Attachment:  appdb.NewSqlNullString(""),
		}
		_, _ = uc.messageRepo.Create(ctx, message)
	}

	return edited, nil
}
