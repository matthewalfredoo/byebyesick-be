package api

import (
	"database/sql"
	"halodeksik-be/app/repository"
)

type AllRepositories struct {
	AddressAreaRepository                 repository.AddressAreaRepository
	CartItemRepository                    repository.CartItemRepository
	CronRepository                        repository.CronRepository
	ConsultationMessageRepository         repository.ConsultationMessageRepository
	ConsultationSessionRepository         repository.ConsultationSessionRepository
	DoctorSpecializationRepository        repository.DoctorSpecializationRepository
	DrugClassificationRepository          repository.DrugClassificationRepository
	ForgotTokenRepository                 repository.ForgotTokenRepository
	ManufacturerRepository                repository.ManufacturerRepository
	OrderRepository                       repository.OrderRepository
	PharmacyRepository                    repository.PharmacyRepository
	PharmacyProductRepository             repository.PharmacyProductRepository
	PrescriptionRepository                repository.PrescriptionRepository
	ProductCategoryRepository             repository.ProductCategoryRepository
	ProductRepository                     repository.ProductRepository
	ProductStockMutationRepository        repository.ProductStockMutationRepository
	ProductStockMutationRequestRepository repository.ProductStockMutationRequestRepository
	ProfileRepository                     repository.ProfileRepository
	RegisterTokenRepository               repository.RegisterTokenRepository
	ReportRepository                      repository.ReportRepository
	TransactionRepository                 repository.TransactionRepository
	ShippingMethodRepository              repository.ShippingMethodRepository
	SickLeaveFormRepository               repository.SickLeaveFormRepository
	UserAddressRepository                 repository.UserAddressRepository
	UserRepository                        repository.UserRepository
}

func InitializeRepositories(db *sql.DB) *AllRepositories {
	return &AllRepositories{
		AddressAreaRepository:                 repository.NewAddressAreaRepositoryImpl(db),
		CartItemRepository:                    repository.NewCartItemRepositoryImpl(db),
		CronRepository:                        repository.NewCronRepoImpl(db),
		ConsultationMessageRepository:         repository.NewConsultationMessageRepositoryImpl(db),
		ConsultationSessionRepository:         repository.NewConsultationSessionRepositoryImpl(db),
		DoctorSpecializationRepository:        repository.NewDoctorSpecializationRepositoryImpl(db),
		DrugClassificationRepository:          repository.NewDrugClassificationRepositoryImpl(db),
		ForgotTokenRepository:                 repository.NewForgotTokenRepository(db),
		ManufacturerRepository:                repository.NewManufacturerRepositoryImpl(db),
		OrderRepository:                       repository.NewOrderRepositoryImpl(db),
		PharmacyRepository:                    repository.NewPharmacyRepository(db),
		PharmacyProductRepository:             repository.NewPharmacyProductRepository(db),
		PrescriptionRepository:                repository.NewPrescriptionRepositoryImpl(db),
		ProductCategoryRepository:             repository.NewProductCategoryRepositoryImpl(db),
		ProductRepository:                     repository.NewProductRepositoryImpl(db),
		ProductStockMutationRepository:        repository.NewProductStockMutationRepositoryImpl(db),
		ProductStockMutationRequestRepository: repository.NewProductStockMutationRequestRepositoryImpl(db),
		ProfileRepository:                     repository.NewProfileRepository(db),
		RegisterTokenRepository:               repository.NewRegisterTokenRepository(db),
		ReportRepository:                      repository.NewReportRepositoryImpl(db),
		TransactionRepository:                 repository.NewTransactionRepositoryImpl(db),
		ShippingMethodRepository:              repository.NewShippingMethodRepositoryImpl(db),
		SickLeaveFormRepository:               repository.NewSickLeaveFormRepositoryImpl(db),
		UserRepository:                        repository.NewUserRepository(db),
		UserAddressRepository:                 repository.NewUserAddressRepositoryImpl(db),
	}
}
