package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type SickLeaveFormRepository interface {
	Create(ctx context.Context, form entity.SickLeaveForm) (*entity.SickLeaveForm, error)
	FindBySessionId(ctx context.Context, sessionId int64) (*entity.SickLeaveForm, error)
	FindBySessionIdDetailed(ctx context.Context, sessionId int64) (*entity.SickLeaveForm, error)
	UpdateBySessionId(ctx context.Context, form entity.SickLeaveForm) (*entity.SickLeaveForm, error)
}

type SickLeaveFormRepositoryImpl struct {
	db *sql.DB
}

func NewSickLeaveFormRepositoryImpl(db *sql.DB) *SickLeaveFormRepositoryImpl {
	return &SickLeaveFormRepositoryImpl{db: db}
}

func (repo *SickLeaveFormRepositoryImpl) Create(ctx context.Context, form entity.SickLeaveForm) (*entity.SickLeaveForm, error) {
	const createSickLeave = `
	INSERT INTO sick_leave_forms(session_id, starting_date, ending_date, description)
	VALUES ($1, $2, $3, $4)
	RETURNING id, session_id, starting_date, ending_date, description, created_at, updated_at`

	row := repo.db.QueryRowContext(ctx, createSickLeave, form.SessionId, form.StartingDate, form.EndingDate, form.Description)
	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrConsultationSessionAlreadyHasSickLeaveForm
		}
		return nil, row.Err()
	}

	var created entity.SickLeaveForm
	err := row.Scan(&created.Id, &created.SessionId, &created.StartingDate, &created.EndingDate, &created.Description, &created.CreatedAt, &created.UpdatedAt)

	return &created, err
}

func (repo *SickLeaveFormRepositoryImpl) FindBySessionId(ctx context.Context, sessionId int64) (*entity.SickLeaveForm, error) {
	const findById = `
	SELECT id, session_id, starting_date, ending_date, description, created_at, updated_at FROM sick_leave_forms
	WHERE session_id = $1`

	row := repo.db.QueryRowContext(ctx, findById, sessionId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var form entity.SickLeaveForm
	err := row.Scan(
		&form.Id, &form.SessionId, &form.StartingDate, &form.EndingDate, &form.Description, &form.CreatedAt, &form.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	return &form, nil
}

func (repo *SickLeaveFormRepositoryImpl) FindBySessionIdDetailed(ctx context.Context, sessionId int64) (*entity.SickLeaveForm, error) {
	const findById = `
	SELECT sick_leave_forms.id, sick_leave_forms.session_id, starting_date, ending_date,
		   description, sick_leave_forms.created_at, sick_leave_forms.updated_at,
		   prescriptions.symptoms, prescriptions.diagnosis,
		   user_profiles.name, user_profiles.date_of_birth, users1.email,
		   doctor_profiles.name, doctor_specializations.name, users2.email
	FROM sick_leave_forms
			 INNER JOIN consultation_sessions ON sick_leave_forms.session_id = consultation_sessions.id
			 INNER JOIN prescriptions ON consultation_sessions.id = prescriptions.session_id
			 INNER JOIN user_profiles ON consultation_sessions.user_id = user_profiles.user_id
			 INNER JOIN users AS users1 ON user_profiles.user_id = users1.id
			 INNER JOIN doctor_profiles ON consultation_sessions.doctor_id = doctor_profiles.user_id
			 INNER JOIN users AS users2 ON doctor_profiles.user_id = users2.id
			 INNER JOIN doctor_specializations ON doctor_profiles.doctor_specialization_id = doctor_specializations.id
	WHERE sick_leave_forms.session_id = $1;`

	row := repo.db.QueryRowContext(ctx, findById, sessionId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var (
		form                 entity.SickLeaveForm
		prescription         entity.Prescription
		user                 entity.User
		userProfile          entity.UserProfile
		doctor               entity.User
		doctorProfile        entity.DoctorProfile
		doctorSpecialization entity.DoctorSpecialization
	)
	err := row.Scan(
		&form.Id, &form.SessionId, &form.StartingDate, &form.EndingDate, &form.Description, &form.CreatedAt, &form.UpdatedAt,
		&prescription.Symptoms, &prescription.Diagnosis,
		&userProfile.Name, &userProfile.DateOfBirth, &user.Email,
		&doctorProfile.Name, &doctorSpecialization.Name, &doctor.Email,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	user.UserProfile = &userProfile
	doctorProfile.DoctorSpecialization = &doctorSpecialization
	doctor.DoctorProfile = &doctorProfile

	form.User = &user
	form.Doctor = &doctor
	form.Prescription = &prescription

	return &form, nil
}

func (repo *SickLeaveFormRepositoryImpl) UpdateBySessionId(ctx context.Context, form entity.SickLeaveForm) (*entity.SickLeaveForm, error) {
	const updateSickLeave = `
	UPDATE sick_leave_forms
	SET starting_date = $1, ending_date = $2, description = $3, updated_at = now()
	WHERE session_id = $4
	RETURNING id, session_id, starting_date, ending_date, description, created_at, updated_at`

	row := repo.db.QueryRowContext(ctx, updateSickLeave, form.StartingDate, form.EndingDate, form.Description, form.SessionId)
	var updated entity.SickLeaveForm
	err := row.Scan(&updated.Id, &updated.SessionId, &updated.StartingDate, &updated.EndingDate, &updated.Description, &updated.CreatedAt, &updated.UpdatedAt)

	return &updated, err
}
