package postgres

import (
	"context"
	"database/sql"
	"nusatek-backend/internal/domain"
)

type propertyRepository struct {
	Conn *sql.DB
}

func NewPropertyRepository(Conn *sql.DB) domain.PropertyRepository {
	return &propertyRepository{Conn}
}

func (m *propertyRepository) Fetch(ctx context.Context, limit int, offset int) ([]domain.Property, error) {
	query := `SELECT id, title, description, address, price, created_at, updated_at FROM properties LIMIT $1 OFFSET $2`
	rows, err := m.Conn.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var properties []domain.Property
	for rows.Next() {
		var p domain.Property
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Address, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		properties = append(properties, p)
	}
	return properties, nil
}

func (m *propertyRepository) GetByID(ctx context.Context, id int64) (domain.Property, error) {
	query := `SELECT id, title, description, address, price, created_at, updated_at FROM properties WHERE id = $1`
	row := m.Conn.QueryRowContext(ctx, query, id)

	var p domain.Property
	err := row.Scan(&p.ID, &p.Title, &p.Description, &p.Address, &p.Price, &p.CreatedAt, &p.UpdatedAt)
	return p, err
}

func (m *propertyRepository) Store(ctx context.Context, p *domain.Property) error {
	query := `INSERT INTO properties (title, description, address, price, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	return m.Conn.QueryRowContext(ctx, query, p.Title, p.Description, p.Address, p.Price).Scan(&p.ID)
}

func (m *propertyRepository) Update(ctx context.Context, p *domain.Property) error {
	query := `UPDATE properties SET title=$1, description=$2, address=$3, price=$4, updated_at=NOW() WHERE id=$5`
	_, err := m.Conn.ExecContext(ctx, query, p.Title, p.Description, p.Address, p.Price, p.ID)
	return err
}

func (m *propertyRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM properties WHERE id = $1`
	_, err := m.Conn.ExecContext(ctx, query, id)
	return err
}
