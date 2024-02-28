package api

import (
	"halodeksik-be/app/appcloud"
	"halodeksik-be/app/usecase"
)

type AllUseCases struct {
	AddressAreaUseCase          usecase.AddressAreaUseCase
	AuthUseCase                 usecase.AuthUsecase
	CartItemUseCase             usecase.CartItemUseCase
	ConsultationMessageUseCase  usecase.ConsultationMessageUseCase
	ConsultationSessionUseCase  usecase.ConsultationSessionUseCase
	CronUseCase                 usecase.CronUseCase
	DoctorSpecializationUseCase usecase.DoctorSpecializationUseCase
	DrugClassificationUseCase   usecase.DrugClassificationUseCase
	ForgotTokenUseCase          usecase.ForgotTokenUseCase
	ManufacturerUseCase         usecase.ManufacturerUseCase
	OrderUseCase                usecase.OrderUseCase
	PharmacyUseCase             usecase.PharmacyUseCase
	PharmacyProductUseCase      usecase.PharmacyProductUseCase
	PrescriptionUseCase         usecase.PrescriptionUseCase
	ProductCategoryUseCase      usecase.ProductCategoryUseCase
	ProductStockMutation        usecase.ProductStockMutationUseCase
	ProductStockMutationRequest usecase.ProductStockMutationRequestUseCase
	ProductUseCase              usecase.ProductUseCase
	ProfileUseCase              usecase.ProfileUseCase
	RegisterTokenUseCase        usecase.RegisterTokenUseCase
	ReportUseCase               usecase.ReportUseCase
	TransactionUseCase          usecase.TransactionUseCase
	ShippingMethodUseCase       usecase.ShippingMethodUseCase
	SickLeaveFormUseCase        usecase.SickLeaveFormUseCase
	UserAddressUseCase          usecase.AddressUseCase
	UserUseCase                 usecase.UserUseCase
}

func InitializeUseCases(allRepo *AllRepositories, allUtil *AllUtil) *AllUseCases {

	forgotTokenUseCase := usecase.NewForgotTokenUsecase(allRepo.UserRepository, allRepo.ForgotTokenRepository, allUtil.AuthUtil, allUtil.MailUtil)
	registerTokenUseCase := usecase.NewRegisterTokenUseCase(allRepo.UserRepository, allRepo.RegisterTokenRepository, allUtil.AuthUtil, allUtil.MailUtil)
	authRepos := usecase.AuthRepos{
		UserRepo:      allRepo.UserRepository,
		TForgotRepo:   allRepo.ForgotTokenRepository,
		TRegisterRepo: allRepo.RegisterTokenRepository,
		ProfileRepo:   allRepo.ProfileRepository,
	}
	authCases := usecase.AuthUseCases{TForgotUseCase: forgotTokenUseCase, TRegisterUseCase: registerTokenUseCase}

	return &AllUseCases{
		AddressAreaUseCase:          usecase.NewAddressAreaUseCaseImpl(allRepo.AddressAreaRepository, allUtil.LocUtil),
		AuthUseCase:                 usecase.NewAuthUsecase(authRepos, allUtil.AuthUtil, appcloud.AppFileUploader, authCases),
		CartItemUseCase:             usecase.NewCartItemUseCaseImpl(allRepo.CartItemRepository, allRepo.ProductRepository, allRepo.PharmacyProductRepository),
		CronUseCase:                 usecase.NewCronUseCase(allRepo.CronRepository),
		ConsultationMessageUseCase:  usecase.NewConsultationMessageUseCaseImpl(allRepo.ConsultationMessageRepository),
		ConsultationSessionUseCase:  usecase.NewConsultationSessionUseCaseImpl(allRepo.ConsultationSessionRepository, allRepo.PrescriptionRepository, allRepo.SickLeaveFormRepository, allRepo.UserRepository),
		DrugClassificationUseCase:   usecase.NewDrugClassificationUseCaseImpl(allRepo.DrugClassificationRepository),
		DoctorSpecializationUseCase: usecase.NewDoctorSpecializationUseCaseImpl(allRepo.DoctorSpecializationRepository, appcloud.AppFileUploader),
		ForgotTokenUseCase:          forgotTokenUseCase,
		ManufacturerUseCase:         usecase.NewManufacturerUseCaseImpl(allRepo.ManufacturerRepository, appcloud.AppFileUploader),
		OrderUseCase:                usecase.NewOrderUseCaseImpl(allRepo.OrderRepository),
		PharmacyUseCase:             usecase.NewPharmacyUseCaseImpl(allRepo.PharmacyRepository, allRepo.AddressAreaRepository),
		PharmacyProductUseCase:      usecase.NewPharmacyProductUseCaseImpl(allRepo.PharmacyProductRepository, allRepo.PharmacyRepository, allRepo.ProductRepository),
		PrescriptionUseCase:         usecase.NewPrescriptionUseCaseImpl(allRepo.PrescriptionRepository, allRepo.ConsultationSessionRepository),
		ProductCategoryUseCase:      usecase.NewProductCategoryUseCaseImpl(allRepo.ProductCategoryRepository),
		ProductUseCase:              usecase.NewProductUseCaseImpl(allRepo.ProductRepository, allRepo.PharmacyRepository, appcloud.AppFileUploader),
		ProductStockMutation:        usecase.NewProductStockMutationUseCaseImpl(allRepo.ProductStockMutationRepository, allRepo.PharmacyProductRepository, allRepo.PharmacyRepository),
		ProductStockMutationRequest: usecase.NewProductStockMutationRequestUseCaseImpl(allRepo.ProductStockMutationRequestRepository, allRepo.PharmacyProductRepository, allRepo.PharmacyRepository),
		ProfileUseCase:              usecase.NewProfileUseCaseImpl(allRepo.ProfileRepository, appcloud.AppFileUploader),
		ShippingMethodUseCase:       usecase.NewShippingMethodUseCaseImpl(allRepo.ShippingMethodRepository, allRepo.UserAddressRepository, allRepo.AddressAreaRepository, allRepo.PharmacyProductRepository, allUtil.OngkirUtil),
		SickLeaveFormUseCase:        usecase.NewSickLeaveFormUseCaseImpl(allRepo.SickLeaveFormRepository, allRepo.ConsultationSessionRepository, allRepo.PrescriptionRepository, allRepo.ConsultationMessageRepository),
		RegisterTokenUseCase:        registerTokenUseCase,
		ReportUseCase:               usecase.NewReportUseCaseImpl(allRepo.ReportRepository),
		TransactionUseCase:          usecase.NewTransactionUseCaseImpl(allRepo.TransactionRepository, allRepo.UserAddressRepository, allRepo.PharmacyProductRepository, appcloud.AppFileUploader),
		UserUseCase:                 usecase.NewUserUseCaseImpl(allRepo.UserRepository, allRepo.PharmacyRepository, allUtil.AuthUtil),
		UserAddressUseCase:          usecase.NewAddressUseCaseImpl(allRepo.UserAddressRepository, allRepo.AddressAreaRepository, allUtil.LocUtil),
	}
}
