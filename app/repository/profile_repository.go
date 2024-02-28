package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type ProfileRepository interface {
	FindUserProfileByUserId(ctx context.Context, userId int64) (*entity.User, error)
	FindDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.User, error)
	UpdateUserProfileByUserId(ctx context.Context, profile entity.UserProfile) (*entity.UserProfile, error)
	UpdateDoctorProfileByUserId(ctx context.Context, profile entity.DoctorProfile) (*entity.DoctorProfile, error)
}

type ProfileRepositoryImpl struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepositoryImpl {
	repo := ProfileRepositoryImpl{db: db}
	return &repo
}

func (repo *ProfileRepositoryImpl) FindUserProfileByUserId(ctx context.Context, userId int64) (*entity.User, error) {
	const getUserWithProfile = `
	SELECT id, email, user_role_id, is_verified, name, profile_photo, date_of_birth
	FROM users INNER JOIN user_profiles ON users.id = user_profiles.user_id WHERE users.id = $1
	`

	row := repo.db.QueryRowContext(ctx, getUserWithProfile,
		userId,
	)

	var user entity.User
	var (
		profile entity.UserProfile
	)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.UserRoleId,
		&user.IsVerified,
		&profile.Name,
		&profile.ProfilePhoto,
		&profile.DateOfBirth,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	user.UserProfile = &profile
	return &user, nil
}

func (repo *ProfileRepositoryImpl) FindDoctorProfileByUserId(ctx context.Context, userId int64) (*entity.User, error) {
	const getDoctorWithProfile = `
	SELECT u.id, email, user_role_id, is_verified, user_id, dp.name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id, consultation_fee, is_online, ds.name spec, ds.id dsId
	FROM users u INNER JOIN doctor_profiles dp ON u.id = dp.user_id INNER JOIN doctor_specializations ds ON dp.doctor_specialization_id = ds.id WHERE u.id = $1
	`

	row := repo.db.QueryRowContext(ctx, getDoctorWithProfile,
		userId,
	)

	var doctor entity.User
	var (
		profile entity.DoctorProfile
		spec    entity.DoctorSpecialization
	)

	err := row.Scan(
		&doctor.Id,
		&doctor.Email,
		&doctor.UserRoleId,
		&doctor.IsVerified,
		&profile.UserId,
		&profile.Name,
		&profile.ProfilePhoto,
		&profile.StartingYear,
		&profile.DoctorCertificate,
		&profile.DoctorSpecializationId,
		&profile.ConsultationFee,
		&profile.IsOnline,
		&spec.Name,
		&spec.Id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}
	profile.DoctorSpecialization = &spec
	doctor.DoctorProfile = &profile

	return &doctor, nil
}

func (repo *ProfileRepositoryImpl) UpdateUserProfileByUserId(ctx context.Context, profile entity.UserProfile) (*entity.UserProfile, error) {
	const updateUserProfileByUserId = `
	UPDATE user_profiles
	SET name = $1, profile_photo = $2, date_of_birth = $3, updated_at = now() WHERE user_id = $4
	RETURNING user_id, name, profile_photo, date_of_birth, created_at, updated_at, deleted_at
	`

	row := repo.db.QueryRowContext(ctx, updateUserProfileByUserId,
		profile.Name, profile.ProfilePhoto, profile.DateOfBirth, profile.UserId,
	)

	var updatedProfile entity.UserProfile

	err := row.Scan(
		&updatedProfile.UserId,
		&updatedProfile.Name,
		&updatedProfile.ProfilePhoto,
		&updatedProfile.DateOfBirth,
		&updatedProfile.CreatedAt,
		&updatedProfile.UpdatedAt,
		&updatedProfile.DeletedAt,
	)

	return &updatedProfile, err

}

func (repo *ProfileRepositoryImpl) UpdateDoctorProfileByUserId(ctx context.Context, profile entity.DoctorProfile) (*entity.DoctorProfile, error) {
	const updateDoctorProfileByUserId = `
	WITH updated_profile AS (
		UPDATE doctor_profiles
			SET name = $1, profile_photo =  $2, starting_year =  $3, doctor_certificate =  $4, doctor_specialization_id = $5, consultation_fee = $6, is_online = $7, updated_at = now() WHERE user_id = $8
			RETURNING user_id, name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id, consultation_fee, is_online
	) SELECT up.*, ds.name, ds.id FROM updated_profile up INNER JOIN doctor_specializations ds ON up.doctor_specialization_id = ds.id;
	`
	row := repo.db.QueryRowContext(ctx, updateDoctorProfileByUserId,
		profile.Name, profile.ProfilePhoto, profile.StartingYear, profile.DoctorCertificate, profile.DoctorSpecializationId,
		profile.ConsultationFee, profile.IsOnline, profile.UserId,
	)

	var updatedProfile entity.DoctorProfile
	var spec entity.DoctorSpecialization

	err := row.Scan(
		&updatedProfile.UserId,
		&updatedProfile.Name,
		&updatedProfile.ProfilePhoto,
		&updatedProfile.StartingYear,
		&updatedProfile.DoctorCertificate,
		&updatedProfile.DoctorSpecializationId,
		&updatedProfile.ConsultationFee,
		&updatedProfile.IsOnline,
		&spec.Name,
		&spec.Id,
	)

	updatedProfile.DoctorSpecialization = &spec
	return &updatedProfile, err

}
