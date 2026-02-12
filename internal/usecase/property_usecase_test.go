package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"nusatek-backend/internal/domain"
	"nusatek-backend/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockPropertyRepo struct {
	mock.Mock
}

func (m *MockPropertyRepo) Fetch(ctx context.Context, limit, offset int) ([]domain.Property, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]domain.Property), args.Error(1)
}
func (m *MockPropertyRepo) GetByID(ctx context.Context, id int64) (domain.Property, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Property), args.Error(1)
}
func (m *MockPropertyRepo) Store(ctx context.Context, p *domain.Property) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockPropertyRepo) Update(ctx context.Context, p *domain.Property) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}
func (m *MockPropertyRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockCacheRepo struct {
	mock.Mock
}

func (m *MockCacheRepo) Get(ctx context.Context, key string) (*domain.Property, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Property), args.Error(1)
}
func (m *MockCacheRepo) Set(ctx context.Context, key string, p *domain.Property, ttl time.Duration) error {
	args := m.Called(ctx, key, p, ttl)
	return args.Error(0)
}
func (m *MockCacheRepo) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func TestGetByID(t *testing.T) {
	mockRepo := new(MockPropertyRepo)
	mockCache := new(MockCacheRepo)
	// Passing nil for amqp channel since we are not testing Store here, or we can mock it if needed but it's a struct pointer in implementation, strict dependency injection would be better with interface.
	// For this test we only test GetByID which doesn't use RabbitMQ.
	u := usecase.NewPropertyUsecase(mockRepo, mockCache, nil, 2*time.Second)

	t.Run("success from cache", func(t *testing.T) {
		mockProp := &domain.Property{ID: 1, Title: "Test Property"}
		mockCache.On("Get", mock.Anything, "property:1").Return(mockProp, nil).Once()

		res, err := u.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, "Test Property", res.Title)
		mockCache.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("success from db", func(t *testing.T) {
		mockProp := domain.Property{ID: 2, Title: "DB Property"}
		mockCache.On("Get", mock.Anything, "property:2").Return(nil, errors.New("not found")).Once()
		mockRepo.On("GetByID", mock.Anything, int64(2)).Return(mockProp, nil).Once()
		mockCache.On("Set", mock.Anything, "property:2", &mockProp, mock.Anything).Return(nil).Once()

		res, err := u.GetByID(context.Background(), 2)

		assert.NoError(t, err)
		assert.Equal(t, "DB Property", res.Title)
		mockRepo.AssertExpectations(t)
		mockCache.AssertExpectations(t)
	})
}
