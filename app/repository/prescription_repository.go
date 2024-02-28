package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"halodeksik-be/app/apperror"
	"halodeksik-be/app/entity"
)

type PrescriptionRepository interface {
	Create(ctx context.Context, prescription entity.Prescription) (*entity.Prescription, error)
	FindBySessionId(ctx context.Context, sessionId int64) (*entity.Prescription, error)
	FindBySessionIdDetailed(ctx context.Context, sessionId int64) (*entity.Prescription, error)
	UpdateBySessionId(ctx context.Context, prescription entity.Prescription) (*entity.Prescription, error)
}

type PrescriptionRepositoryImpl struct {
	db *sql.DB
}

func NewPrescriptionRepositoryImpl(db *sql.DB) *PrescriptionRepositoryImpl {
	return &PrescriptionRepositoryImpl{db: db}
}

func (repo *PrescriptionRepositoryImpl) Create(ctx context.Context, prescription entity.Prescription) (*entity.Prescription, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	const createPrescription = `
	INSERT INTO prescriptions(session_id, symptoms, diagnosis)
	VALUES ($1, $2, $3)
	RETURNING id, session_id, symptoms, diagnosis, created_at, updated_at`

	row := tx.QueryRowContext(ctx, createPrescription, prescription.SessionId, prescription.Symptoms, prescription.Diagnosis)
	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrConsultationSessionAlreadyHasPrescription
		}
		return nil, row.Err()
	}

	var createdPrescription entity.Prescription
	err = row.Scan(
		&createdPrescription.Id, &createdPrescription.SessionId,
		&createdPrescription.Symptoms, &createdPrescription.Diagnosis,
		&createdPrescription.CreatedAt, &createdPrescription.UpdatedAt,
	)

	createPrescriptionProduct := `INSERT INTO prescription_products(prescription_id, product_id, note) VALUES `
	values := make([]interface{}, 0)

	indexPreparedStatement := 0
	for index, prescriptionProduct := range prescription.PrescriptionProducts {
		createPrescriptionProduct += fmt.Sprintf("($%d, $%d, $%d)", indexPreparedStatement+1, indexPreparedStatement+2, indexPreparedStatement+3)
		indexPreparedStatement += 3
		if index != len(prescription.PrescriptionProducts)-1 {
			createPrescriptionProduct += ", "
		}

		values = append(values,
			createdPrescription.Id, prescriptionProduct.ProductId, prescriptionProduct.Note,
		)
	}
	createPrescriptionProduct += " RETURNING id, prescription_id, product_id, note, created_at, updated_at"

	rows, err := tx.QueryContext(ctx, createPrescriptionProduct, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prescriptionProducts := make([]*entity.PrescriptionProduct, 0)
	for rows.Next() {
		var prescriptionProduct entity.PrescriptionProduct
		if err = rows.Scan(
			&prescriptionProduct.Id, &prescriptionProduct.PrescriptionId, &prescriptionProduct.ProductId,
			&prescriptionProduct.Note, &prescriptionProduct.CreatedAt, &prescriptionProduct.UpdatedAt,
		); err != nil {
			return nil, err
		}
		prescriptionProducts = append(prescriptionProducts, &prescriptionProduct)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	createdPrescription.PrescriptionProducts = prescriptionProducts

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &createdPrescription, err
}

func (repo *PrescriptionRepositoryImpl) FindBySessionId(ctx context.Context, sessionId int64) (*entity.Prescription, error) {
	query := `
	SELECT prescriptions.id, session_id, symptoms, diagnosis, prescriptions.created_at, prescriptions.updated_at,
		   cm.prescription_product_id, cm.prescription_product_product_id, cm.note,  cm.created_at, cm.updated_at,
		   cm.product_id, cm.product_name, cm.product_generic_name, cm.product_content, cm.product_image,
		   cm.manufacturer_name
	FROM  prescriptions
		LEFT JOIN LATERAL (
			SELECT prescription_products.id AS prescription_product_id, product_id AS prescription_product_product_id, note, prescription_products.created_at, prescription_products.updated_at,
				   products.id AS product_id, products.name AS product_name, products.generic_name AS product_generic_name, products.content AS product_content, products.image AS product_image,
				   manufacturers.name AS manufacturer_name
			FROM prescription_products
			INNER JOIN products ON prescription_products.product_id = products.id
			INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
			WHERE prescription_products.prescription_id = prescriptions.id
			ORDER BY prescription_products.id ASC
		) cm ON true
	WHERE prescriptions.deleted_at IS NULL AND prescriptions.session_id = $1;`

	rows, err := repo.db.QueryContext(ctx, query, sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prescription entity.Prescription
	prescriptionProducts := make([]*entity.PrescriptionProduct, 0)
	for rows.Next() {
		var (
			prescriptionProduct entity.PrescriptionProduct
			product             entity.Product
			manufacturer        entity.Manufacturer
		)
		if err = rows.Scan(
			&prescription.Id, &prescription.SessionId, &prescription.Symptoms, &prescription.Diagnosis, &prescription.CreatedAt, &prescription.UpdatedAt,
			&prescriptionProduct.Id, &prescriptionProduct.ProductId, &prescriptionProduct.Note, &prescriptionProduct.CreatedAt, &prescriptionProduct.UpdatedAt,
			&product.Id, &product.Name, &product.GenericName, &product.Content, &product.Image,
			&manufacturer.Name,
		); err != nil {
			return nil, err
		}
		product.Manufacturer = &manufacturer
		prescriptionProduct.Product = &product
		prescriptionProducts = append(prescriptionProducts, &prescriptionProduct)
	}
	prescription.PrescriptionProducts = prescriptionProducts

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if prescription.Id == 0 {
		return nil, apperror.ErrRecordNotFound
	}

	return &prescription, nil
}

func (repo *PrescriptionRepositoryImpl) FindBySessionIdDetailed(ctx context.Context, sessionId int64) (*entity.Prescription, error) {
	query := `
	SELECT prescriptions.id, session_id, symptoms, diagnosis, prescriptions.created_at, prescriptions.updated_at,
		   user_profiles.name, user_profiles.date_of_birth, users.email,
		   doctor_profiles.name, doctor_specializations.name, doctors.email,
		   cm.prescription_product_id, cm.prescription_product_product_id, cm.note,  cm.created_at, cm.updated_at,
		   cm.product_id, cm.product_name, cm.product_generic_name, cm.product_content, cm.product_image,
		   cm.manufacturer_name
	FROM  prescriptions
		INNER JOIN consultation_sessions ON prescriptions.session_id = consultation_sessions.id
		INNER JOIN user_profiles ON consultation_sessions.user_id = user_profiles.user_id
		INNER JOIN users ON user_profiles.user_id = users.id
		INNER JOIN doctor_profiles ON consultation_sessions.doctor_id = doctor_profiles.user_id
		INNER JOIN users AS doctors ON doctor_profiles.user_id = doctors.id
		INNER JOIN doctor_specializations ON doctor_profiles.doctor_specialization_id = doctor_specializations.id
		LEFT JOIN LATERAL (
		SELECT prescription_products.id AS prescription_product_id, product_id AS prescription_product_product_id, note, prescription_products.created_at, prescription_products.updated_at,
			   products.id AS product_id, products.name AS product_name, products.generic_name AS product_generic_name, products.content AS product_content, products.image AS product_image,
			   manufacturers.name AS manufacturer_name
		FROM prescription_products
				 INNER JOIN products ON prescription_products.product_id = products.id
				 INNER JOIN manufacturers ON products.manufacturer_id = manufacturers.id
		WHERE prescription_products.prescription_id = prescriptions.id
		ORDER BY prescription_products.id ASC
		) cm ON true
	WHERE prescriptions.deleted_at IS NULL AND prescriptions.session_id = $1`

	rows, err := repo.db.QueryContext(ctx, query, sessionId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		prescription         entity.Prescription
		user                 entity.User
		userProfile          entity.UserProfile
		doctor               entity.User
		doctorProfile        entity.DoctorProfile
		doctorSpecialization entity.DoctorSpecialization
	)
	prescriptionProducts := make([]*entity.PrescriptionProduct, 0)
	for rows.Next() {
		var (
			prescriptionProduct entity.PrescriptionProduct
			product             entity.Product
			manufacturer        entity.Manufacturer
		)
		if err = rows.Scan(
			&prescription.Id, &prescription.SessionId, &prescription.Symptoms, &prescription.Diagnosis, &prescription.CreatedAt, &prescription.UpdatedAt,
			&userProfile.Name, &userProfile.DateOfBirth, &user.Email,
			&doctorProfile.Name, &doctorSpecialization.Name, &doctor.Email,
			&prescriptionProduct.Id, &prescriptionProduct.ProductId, &prescriptionProduct.Note, &prescriptionProduct.CreatedAt, &prescriptionProduct.UpdatedAt,
			&product.Id, &product.Name, &product.GenericName, &product.Content, &product.Image,
			&manufacturer.Name,
		); err != nil {
			return nil, err
		}
		product.Manufacturer = &manufacturer
		prescriptionProduct.Product = &product
		prescriptionProducts = append(prescriptionProducts, &prescriptionProduct)
	}
	prescription.PrescriptionProducts = prescriptionProducts

	user.UserProfile = &userProfile
	doctorProfile.DoctorSpecialization = &doctorSpecialization
	doctor.DoctorProfile = &doctorProfile

	prescription.User = &user
	prescription.Doctor = &doctor

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if prescription.Id == 0 {
		return nil, apperror.ErrRecordNotFound
	}

	return &prescription, nil
}

func (repo *PrescriptionRepositoryImpl) UpdateBySessionId(ctx context.Context, prescription entity.Prescription) (*entity.Prescription, error) {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {
			return
		}
	}(tx)

	const getPrescriptionBySessionId = `SELECT id FROM prescriptions WHERE session_id = $1`

	row := repo.db.QueryRowContext(ctx, getPrescriptionBySessionId, prescription.SessionId)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var prescriptionId int64
	err = row.Scan(&prescriptionId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrRecordNotFound
		}
		return nil, err
	}

	const deletePrescriptionProduct = `DELETE FROM prescription_products WHERE prescription_id = $1`

	_, err = tx.ExecContext(ctx, deletePrescriptionProduct, prescriptionId)
	if err != nil {
		return nil, err
	}

	const deletePrescription = `DELETE FROM prescriptions WHERE session_id = $1`

	_, err = tx.ExecContext(ctx, deletePrescription, prescription.SessionId)
	if err != nil {
		return nil, err
	}

	const createPrescription = `
	INSERT INTO prescriptions(session_id, symptoms, diagnosis)
	VALUES ($1, $2, $3)
	RETURNING id, session_id, symptoms, diagnosis, created_at, updated_at`

	row = tx.QueryRowContext(ctx, createPrescription, prescription.SessionId, prescription.Symptoms, prescription.Diagnosis)
	if row.Err() != nil {
		var errPgConn *pgconn.PgError
		if errors.As(row.Err(), &errPgConn) && errPgConn.Code == apperror.PgconnErrCodeUniqueConstraintViolation {
			return nil, apperror.ErrConsultationSessionAlreadyHasPrescription
		}
		return nil, row.Err()
	}

	var createdPrescription entity.Prescription
	err = row.Scan(
		&createdPrescription.Id, &createdPrescription.SessionId,
		&createdPrescription.Symptoms, &createdPrescription.Diagnosis,
		&createdPrescription.CreatedAt, &createdPrescription.UpdatedAt,
	)

	createPrescriptionProduct := `INSERT INTO prescription_products(prescription_id, product_id, note) VALUES `
	values := make([]interface{}, 0)

	indexPreparedStatement := 0
	for index, prescriptionProduct := range prescription.PrescriptionProducts {
		createPrescriptionProduct += fmt.Sprintf("($%d, $%d, $%d)", indexPreparedStatement+1, indexPreparedStatement+2, indexPreparedStatement+3)
		indexPreparedStatement += 3
		if index != len(prescription.PrescriptionProducts)-1 {
			createPrescriptionProduct += ", "
		}

		values = append(values,
			createdPrescription.Id, prescriptionProduct.ProductId, prescriptionProduct.Note,
		)
	}
	createPrescriptionProduct += " RETURNING id, prescription_id, product_id, note, created_at, updated_at"

	rows, err := tx.QueryContext(ctx, createPrescriptionProduct, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prescriptionProducts := make([]*entity.PrescriptionProduct, 0)
	for rows.Next() {
		var prescriptionProduct entity.PrescriptionProduct
		if err = rows.Scan(
			&prescriptionProduct.Id, &prescriptionProduct.PrescriptionId, &prescriptionProduct.ProductId,
			&prescriptionProduct.Note, &prescriptionProduct.CreatedAt, &prescriptionProduct.UpdatedAt,
		); err != nil {
			return nil, err
		}
		prescriptionProducts = append(prescriptionProducts, &prescriptionProduct)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	createdPrescription.PrescriptionProducts = prescriptionProducts

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &createdPrescription, err
}
