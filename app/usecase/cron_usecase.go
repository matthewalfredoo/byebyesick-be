package usecase

import (
	"github.com/robfig/cron/v3"
	"halodeksik-be/app/appconstant"
	"halodeksik-be/app/repository"
)

type CronUseCase interface {
	StartCron() error
	ValidateTransactions()
	ValidateOrders()
	ValidateOrdersConfirmed()
}

type CronUseCaseImpl struct {
	cronRepo repository.CronRepository
	cronJob  *cron.Cron
}

func (uc CronUseCaseImpl) ValidateTransactions() {
	err := uc.cronRepo.ValidateTransactions()
	if err != nil {
		return
	}
}

func (uc CronUseCaseImpl) ValidateOrders() {
	err := uc.cronRepo.ValidateOrders()
	if err != nil {
		return
	}
}

func (uc CronUseCaseImpl) ValidateOrdersConfirmed() {
	err := uc.cronRepo.ValidateOrdersConfirmed()
	if err != nil {
		return
	}
}

func NewCronUseCase(cronRepo repository.CronRepository) *CronUseCaseImpl {
	return &CronUseCaseImpl{
		cronRepo: cronRepo,
		cronJob:  cron.New(),
	}
}

func (uc CronUseCaseImpl) StartCron() error {
	_, err := uc.cronJob.AddFunc(appconstant.CronDailyTimer, uc.ValidateTransactions)
	if err != nil {
		return err
	}

	_, err = uc.cronJob.AddFunc(appconstant.CronDailyTimer, uc.ValidateOrders)
	if err != nil {
		return err
	}

	_, err = uc.cronJob.AddFunc(appconstant.CronDailyTimer, uc.ValidateOrdersConfirmed)
	if err != nil {
		return err
	}

	uc.cronJob.Start()

	return nil

}
