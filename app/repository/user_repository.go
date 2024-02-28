package repository

import (
	"context"
	"database/sql"
	"errors"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) (*entity.User, error)
	CreateAndDoctorProfile(ctx context.Context, user entity.User, profile entity.DoctorProfile) (*entity.User, error)
	CreateAndUserProfile(ctx context.Context, user entity.User, profile entity.UserProfile) (*entity.User, error)
	FindById(ctx context.Context, id int64) (*entity.User, error)
	FindDoctorById(ctx context.Context, id int64) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.User, error)
	FindAllDoctors(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.User, error)
	CountFindAllDoctors(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error)
	Update(ctx context.Context, user entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int64) error
	ChangePassword(ctx context.Context, user entity.User, newPassword string) (*entity.User, error)
}

type UserRepositoryImpl struct {
	db *sql.DB
}

func (repo *UserRepositoryImpl) FindAllDoctors(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.User, error) {
	const getAllDoctors = `SELECT users.id, email, user_role_id, is_verified, doctor_profiles.name AS name, 
	doctor_profiles.profile_photo, doctor_profiles.starting_year, doctor_profiles.doctor_certificate,doctor_profiles.is_online, doctor_specializations.id, doctor_specializations.name FROM users
	INNER JOIN doctor_profiles ON users.id = doctor_profiles.user_id INNER JOIN doctor_specializations ON 
	doctor_profiles.doctor_specialization_id = doctor_specializations.id WHERE user_role_id = 3 AND users.deleted_at IS NULL `

	query, values := buildQuery(getAllDoctors, &entity.User{}, param, true, true)
	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		var profile entity.DoctorProfile
		var profileSpec entity.DoctorSpecialization
		if err := rows.Scan(
			&user.Id, &user.Email, &user.UserRoleId, &user.IsVerified, &profile.Name, &profile.ProfilePhoto, &profile.StartingYear,
			&profile.DoctorCertificate, &profile.IsOnline, &profileSpec.Id, &profileSpec.Name,
		); err != nil {
			return nil, err
		}
		profile.DoctorSpecialization = &profileSpec
		user.DoctorProfile = &profile
		items = append(items, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil

}

func (repo *UserRepositoryImpl) CountFindAllDoctors(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(users.id) FROM users
	INNER JOIN doctor_profiles ON users.id = doctor_profiles.user_id INNER JOIN doctor_specializations ON 
	doctor_profiles.doctor_specialization_id = doctor_specializations.id WHERE user_role_id = 3 AND users.deleted_at IS NULL `

	query, values := buildQuery(initQuery, &entity.User{}, param, false, false)

	var totalItems int64

	row := repo.db.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return totalItems, row.Err()
	}

	if err := row.Scan(
		&totalItems,
	); err != nil {
		return totalItems, err
	}

	return totalItems, nil
}

func (repo *UserRepositoryImpl) FindDoctorById(ctx context.Context, id int64) (*entity.User, error) {
	const getDoctorById = `SELECT users.id, email, user_role_id, is_verified, doctor_profiles.name AS name, doctor_profiles.profile_photo, doctor_profiles.starting_year, doctor_profiles.doctor_certificate, doctor_profiles.is_online,doctor_specializations.id, doctor_specializations.name FROM users
	INNER JOIN doctor_profiles ON users.id = doctor_profiles.user_id INNER JOIN doctor_specializations ON doctor_profiles.doctor_specialization_id = doctor_specializations.id
	WHERE user_role_id = 3 AND users.deleted_at IS NULL AND users.id = $1`

	row := repo.db.QueryRowContext(ctx, getDoctorById,
		id,
	)
	var user entity.User
	var profile entity.DoctorProfile
	var profileSpec entity.DoctorSpecialization
	err := row.Scan(
		&user.Id, &user.Email, &user.UserRoleId, &user.IsVerified, &profile.Name, &profile.ProfilePhoto, &profile.StartingYear,
		&profile.DoctorCertificate, &profile.IsOnline, &profileSpec.Id, &profileSpec.Name,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	profile.DoctorSpecialization = &profileSpec
	user.DoctorProfile = &profile

	return &user, err
}

func (repo *UserRepositoryImpl) ChangePassword(ctx context.Context, user entity.User, newPassword string) (*entity.User, error) {
	const updatePasswordById = `UPDATE users
	SET password = $1, updated_at = now()
	WHERE id = $2
	RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, updatePasswordById,
		newPassword,
		user.Id,
	)

	var updated entity.User
	err := row.Scan(
		&updated.Id,
		&updated.Email,
		&updated.Password,
		&updated.UserRoleId,
		&updated.IsVerified,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	)
	return &updated, err
}

func NewUserRepository(db *sql.DB) UserRepository {
	repo := UserRepositoryImpl{db: db}
	return &repo
}

func (repo *UserRepositoryImpl) CreateAndDoctorProfile(ctx context.Context, user entity.User, profile entity.DoctorProfile) (*entity.User, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	const createUser = `INSERT INTO users(email, password, user_role_id, is_verified)
	VALUES ($1, $2, $3, $4)
	RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at`

	row1 := tx.QueryRowContext(ctx, createUser,
		user.Email,
		user.Password,
		user.UserRoleId,
		user.IsVerified,
	)

	var createdUser entity.User

	err = row1.Scan(
		&createdUser.Id,
		&createdUser.Email,
		&createdUser.Password,
		&createdUser.UserRoleId,
		&createdUser.IsVerified,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		&createdUser.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	const createDoctorProfile = `
	INSERT INTO doctor_profiles(user_id, name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id, consultation_fee, is_online)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING user_id, name, profile_photo, starting_year, doctor_certificate, doctor_specialization_id, consultation_fee, is_online, created_at, updated_at, deleted_at
	`

	row2 := tx.QueryRowContext(ctx, createDoctorProfile,
		createdUser.Id,
		profile.Name,
		profile.ProfilePhoto,
		profile.StartingYear,
		profile.DoctorCertificate,
		profile.DoctorSpecializationId,
		profile.ConsultationFee.String(),
		profile.IsOnline,
	)

	var createdProfile entity.DoctorProfile

	err = row2.Scan(
		&createdProfile.UserId,
		&createdProfile.Name,
		&createdProfile.ProfilePhoto,
		&createdProfile.StartingYear,
		&createdProfile.DoctorCertificate,
		&createdProfile.DoctorSpecializationId,
		&createdProfile.ConsultationFee,
		&createdProfile.IsOnline,
		&createdProfile.CreatedAt,
		&createdProfile.UpdatedAt,
		&createdProfile.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &createdUser, nil

}

func (repo *UserRepositoryImpl) CreateAndUserProfile(ctx context.Context, user entity.User, profile entity.UserProfile) (*entity.User, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	const createUser = `INSERT INTO users(email, password, user_role_id, is_verified)
	VALUES ($1, $2, $3, $4)
	RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at`

	row1 := tx.QueryRowContext(ctx, createUser,
		user.Email,
		user.Password,
		user.UserRoleId,
		user.IsVerified,
	)

	var createdUser entity.User

	err = row1.Scan(
		&createdUser.Id,
		&createdUser.Email,
		&createdUser.Password,
		&createdUser.UserRoleId,
		&createdUser.IsVerified,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		&createdUser.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	const createUserProfile = `
	INSERT INTO user_profiles(user_id, name, profile_photo, date_of_birth)
	VALUES ($1, $2, $3, $4)
	RETURNING user_id, name, profile_photo, date_of_birth, created_at, updated_at, deleted_at
	`

	var createdProfile entity.UserProfile

	row2 := tx.QueryRowContext(ctx, createUserProfile,
		createdUser.Id,
		profile.Name,
		profile.ProfilePhoto,
		profile.DateOfBirth,
	)

	err = row2.Scan(
		&createdProfile.UserId,
		&createdProfile.Name,
		&createdProfile.ProfilePhoto,
		&createdProfile.DateOfBirth,
		&createdProfile.CreatedAt,
		&createdProfile.UpdatedAt,
		&createdProfile.DeletedAt,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &createdUser, nil

}

func (repo *UserRepositoryImpl) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	const create = `INSERT INTO users(email, password, user_role_id, is_verified)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at`

	row := repo.db.QueryRowContext(ctx, create,
		user.Email,
		user.Password,
		user.UserRoleId,
		user.IsVerified,
	)

	var createdUser entity.User

	err := row.Scan(
		&createdUser.Id,
		&createdUser.Email,
		&createdUser.Password,
		&createdUser.UserRoleId,
		&createdUser.IsVerified,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
		&createdUser.DeletedAt,
	)

	return &createdUser, err
}

func (repo *UserRepositoryImpl) FindById(ctx context.Context, id int64) (*entity.User, error) {
	const getById = `SELECT id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at
	FROM users WHERE id = $1
`

	row := repo.db.QueryRowContext(ctx, getById, id)
	var user entity.User
	err := row.Scan(
		&user.Id, &user.Email, &user.Password, &user.UserRoleId,
		&user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (repo *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	const getUserByEmail = `SELECT id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at FROM users
WHERE email = $1
`

	row := repo.db.QueryRowContext(ctx, getUserByEmail, email)
	var user entity.User
	err := row.Scan(
		&user.Id, &user.Email, &user.Password, &user.UserRoleId,
		&user.IsVerified, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, apperror.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (repo *UserRepositoryImpl) FindAll(ctx context.Context, param *queryparamdto.GetAllParams) ([]*entity.User, error) {
	initQuery := `SELECT id, email, user_role_id, is_verified FROM users WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.User{}, param, true, true)

	rows, err := repo.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*entity.User, 0)
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.Id, &user.Email, &user.UserRoleId, &user.IsVerified,
		); err != nil {
			return nil, err
		}
		items = append(items, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (repo *UserRepositoryImpl) CountFindAll(ctx context.Context, param *queryparamdto.GetAllParams) (int64, error) {
	initQuery := `SELECT count(id) FROM users WHERE deleted_at IS NULL `
	query, values := buildQuery(initQuery, &entity.User{}, param, false, false)

	var totalItems int64

	row := repo.db.QueryRowContext(ctx, query, values...)
	if row.Err() != nil {
		return totalItems, row.Err()
	}

	if err := row.Scan(
		&totalItems,
	); err != nil {
		return totalItems, err
	}

	return totalItems, nil
}

func (repo *UserRepositoryImpl) Update(ctx context.Context, user entity.User) (*entity.User, error) {
	const updateById = `UPDATE users
SET email = $1, password = $2, user_role_id = $3, is_verified = $4, updated_at = now()
WHERE id = $5
RETURNING id, email, password, user_role_id, is_verified, created_at, updated_at, deleted_at
`

	row := repo.db.QueryRowContext(ctx, updateById,
		user.Email,
		user.Password,
		user.UserRoleId,
		user.IsVerified,
		user.Id,
	)
	var updated entity.User
	err := row.Scan(
		&updated.Id,
		&updated.Email,
		&updated.Password,
		&updated.UserRoleId,
		&updated.IsVerified,
		&updated.CreatedAt,
		&updated.UpdatedAt,
		&updated.DeletedAt,
	)
	return &updated, err
}

func (repo *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	const deleteById = `UPDATE users SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`

	_, err := repo.db.ExecContext(ctx, deleteById, id)
	return err
}
