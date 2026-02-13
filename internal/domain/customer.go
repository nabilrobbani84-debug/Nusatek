package domain

import (
	"context"
	"time"
)

type Customer struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"` // Active, Inactive
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerRepository interface {
	Fetch(ctx context.Context, limit int, offset int) ([]Customer, error)
	GetByID(ctx context.Context, id int64) (Customer, error)
	Store(ctx context.Context, c *Customer) error
	Update(ctx context.Context, c *Customer) error
	Delete(ctx context.Context, id int64) error
}

type CustomerUsecase interface {
	Fetch(ctx context.Context, limit int, offset int) ([]Customer, error)
	GetByID(ctx context.Context, id int64) (Customer, error)
	Store(ctx context.Context, c *Customer) error
	Update(ctx context.Context, c *Customer) error
	Delete(ctx context.Context, id int64) error
}
