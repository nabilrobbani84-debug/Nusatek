package usecase

import (
	"context"
	"nusatek-backend/internal/domain"
	"time"
)

type customerUsecase struct {
	customerRepo domain.CustomerRepository
	contextTimeout time.Duration
}

func NewCustomerUsecase(c domain.CustomerRepository, timeout time.Duration) domain.CustomerUsecase {
	return &customerUsecase{
		customerRepo:   c,
		contextTimeout: timeout,
	}
}

func (du *customerUsecase) Fetch(c context.Context, limit int, offset int) ([]domain.Customer, error) {
	ctx, cancel := context.WithTimeout(c, du.contextTimeout)
	defer cancel()
	return du.customerRepo.Fetch(ctx, limit, offset)
}

func (du *customerUsecase) GetByID(c context.Context, id int64) (domain.Customer, error) {
	ctx, cancel := context.WithTimeout(c, du.contextTimeout)
	defer cancel()
	return du.customerRepo.GetByID(ctx, id)
}

func (du *customerUsecase) Store(c context.Context, m *domain.Customer) error {
	ctx, cancel := context.WithTimeout(c, du.contextTimeout)
	defer cancel()
	return du.customerRepo.Store(ctx, m)
}

func (du *customerUsecase) Update(c context.Context, m *domain.Customer) error {
	ctx, cancel := context.WithTimeout(c, du.contextTimeout)
	defer cancel()
	return du.customerRepo.Update(ctx, m)
}

func (du *customerUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, du.contextTimeout)
	defer cancel()
	return du.customerRepo.Delete(ctx, id)
}
