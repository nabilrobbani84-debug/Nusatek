package postgres

import (
	"context"
	"database/sql"
	"nusatek-backend/internal/domain"
)

type customerRepository struct {
	Conn *sql.DB
}

func NewCustomerRepository(Conn *sql.DB) domain.CustomerRepository {
	return &customerRepository{Conn}
}

func (m *customerRepository) Fetch(ctx context.Context, limit int, offset int) ([]domain.Customer, error) {
	query := `SELECT id, name, email, phone, status, created_at, updated_at FROM customers LIMIT $1 OFFSET $2`
	rows, err := m.Conn.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []domain.Customer
	for rows.Next() {
		var c domain.Customer
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (m *customerRepository) GetByID(ctx context.Context, id int64) (domain.Customer, error) {
	query := `SELECT id, name, email, phone, status, created_at, updated_at FROM customers WHERE id = $1`
	row := m.Conn.QueryRowContext(ctx, query, id)

	var c domain.Customer
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Phone, &c.Status, &c.CreatedAt, &c.UpdatedAt)
	return c, err
}

func (m *customerRepository) Store(ctx context.Context, c *domain.Customer) error {
	query := `INSERT INTO customers (name, email, phone, status, created_at, updated_at) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	return m.Conn.QueryRowContext(ctx, query, c.Name, c.Email, c.Phone, c.Status).Scan(&c.ID)
}

func (m *customerRepository) Update(ctx context.Context, c *domain.Customer) error {
	query := `UPDATE customers SET name=$1, email=$2, phone=$3, status=$4, updated_at=NOW() WHERE id=$5`
	_, err := m.Conn.ExecContext(ctx, query, c.Name, c.Email, c.Phone, c.Status, c.ID)
	return err
}

func (m *customerRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := m.Conn.ExecContext(ctx, query, id)
	return err
}
