package usecase

import (
	"context"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/appdb"
	"halodeksik-be/app/dto/queryparamdto"
	"halodeksik-be/app/entity"
	"halodeksik-be/app/repository"
)

type ReportUseCase interface {
	GetSellsAllPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetSellsAllAdminPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)

	GetSellsAllPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
	GetSellsAllAdminPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error)
}

type ReportUseCaseImpl struct {
	reportRepository repository.ReportRepository
}

func (uc *ReportUseCaseImpl) GetSellsAllPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	reports, err := uc.reportRepository.FindSalesAllPharmacy(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.reportRepository.CountSalesAllPharmacy(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems, totalPages, int64(len(reports)), int64(*param.PageId), reports,
	)

	return paginatedItems, nil
}

func (uc *ReportUseCaseImpl) GetSellsAllAdminPharmacy(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	pharmacy := new(entity.Pharmacy)
	column := pharmacy.GetSqlColumnFromField("PharmacyAdminId")
	param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, userId))

	reports, err := uc.reportRepository.FindSalesAllPharmacy(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.reportRepository.CountSalesAllPharmacy(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems, totalPages, int64(len(reports)), int64(*param.PageId), reports,
	)

	return paginatedItems, nil

}

func (uc *ReportUseCaseImpl) GetSellsAllPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	reports, err := uc.reportRepository.FindSalesAllPharmacyMonthly(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.reportRepository.CountSalesAllPharmacyMonthly(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems, totalPages, int64(len(reports)), int64(*param.PageId), reports,
	)

	return paginatedItems, nil
}

func (uc *ReportUseCaseImpl) GetSellsAllAdminPharmacyMonthly(ctx context.Context, year int64, param *queryparamdto.GetAllParams) (*entity.PaginatedItems, error) {
	userId := ctx.Value(appconstant.ContextKeyUserId).(int64)

	pharmacy := new(entity.Pharmacy)
	column := pharmacy.GetSqlColumnFromField("PharmacyAdminId")
	param.WhereClauses = append(param.WhereClauses, appdb.NewWhere(column, appdb.EqualTo, userId))

	reports, err := uc.reportRepository.FindSalesAllPharmacyMonthly(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalItems, err := uc.reportRepository.CountSalesAllPharmacyMonthly(ctx, year, param)
	if err != nil {
		return nil, err
	}

	totalPages := totalItems / int64(*param.PageSize)
	if totalItems%int64(*param.PageSize) != 0 || totalPages == 0 {
		totalPages += 1
	}

	paginatedItems := entity.NewPaginationInfo(
		totalItems, totalPages, int64(len(reports)), int64(*param.PageId), reports,
	)

	return paginatedItems, nil
}

func NewReportUseCaseImpl(repo repository.ReportRepository) *ReportUseCaseImpl {
	return &ReportUseCaseImpl{reportRepository: repo}
}
