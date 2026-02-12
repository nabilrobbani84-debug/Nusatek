package domain

import (
	"context"
	"time"
)

// Property represents a real estate property
type Property struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Address     string    `json:"address"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PropertyRepository defines the interface for database operations
type PropertyRepository interface {
	Fetch(ctx context.Context, limit int, offset int) ([]Property, error)
	GetByID(ctx context.Context, id int64) (Property, error)
	Store(ctx context.Context, p *Property) error
	Update(ctx context.Context, p *Property) error
	Delete(ctx context.Context, id int64) error
}

// PropertyCacheRepository defines the interface for caching operations
type PropertyCacheRepository interface {
	Get(ctx context.Context, key string) (*Property, error)
	Set(ctx context.Context, key string, p *Property, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// PropertyUsecase defines the interface for business logic
type PropertyUsecase interface {
	Fetch(ctx context.Context, limit int, offset int) ([]Property, error)
	GetByID(ctx context.Context, id int64) (Property, error)
	Store(ctx context.Context, p *Property) error
	Update(ctx context.Context, p *Property) error
	Delete(ctx context.Context, id int64) error
}
