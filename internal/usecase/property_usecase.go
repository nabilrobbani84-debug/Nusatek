package usecase

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"nusatek-backend/internal/domain"

	"github.com/streadway/amqp"
)

type propertyUsecase struct {
	propertyRepo domain.PropertyRepository
	cacheRepo    domain.PropertyCacheRepository
	mqChannel    *amqp.Channel
	timeout      time.Duration
}

func NewPropertyUsecase(a domain.PropertyRepository, c domain.PropertyCacheRepository, mq *amqp.Channel, timeout time.Duration) domain.PropertyUsecase {
	return &propertyUsecase{
		propertyRepo: a,
		cacheRepo:    c,
		mqChannel:    mq,
		timeout:      timeout,
	}
}

func (a *propertyUsecase) Fetch(c context.Context, limit int, offset int) ([]domain.Property, error) {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()
	return a.propertyRepo.Fetch(ctx, limit, offset)
}

func (a *propertyUsecase) GetByID(c context.Context, id int64) (domain.Property, error) {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()

	// 1. Try Cache
	cacheKey := "property:" + strconv.FormatInt(id, 10)
	if cachedProp, err := a.cacheRepo.Get(ctx, cacheKey); err == nil && cachedProp != nil {
		return *cachedProp, nil
	}

	// 2. Fetch from DB
	res, err := a.propertyRepo.GetByID(ctx, id)
	if err != nil {
		return domain.Property{}, err
	}

	// 3. Set Cache (Async or blocking, here blocking for simplicity)
	_ = a.cacheRepo.Set(ctx, cacheKey, &res, 5*time.Minute)

	return res, nil
}

func (a *propertyUsecase) Store(c context.Context, p *domain.Property) error {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()

	// 1. Store in DB
	if err := a.propertyRepo.Store(ctx, p); err != nil {
		return err
	}

	// 2. Publish Event to RabbitMQ
	// We do this asynchronously or synchronously depending on consistency requirements.
	// For this demo, we ignore errors here to not block the response, but in prod we'd handle them.
	body := []byte(fmt.Sprintf(`{"event": "property_created", "id": %d, "title": "%s"}`, p.ID, p.Title))
	_ = a.mqChannel.Publish(
		"",                // exchange
		"property_events", // routing key (queue name)
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	return nil
}

func (a *propertyUsecase) Update(c context.Context, p *domain.Property) error {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()
	return a.propertyRepo.Update(ctx, p)
}

func (a *propertyUsecase) Delete(c context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(c, a.timeout)
	defer cancel()
	return a.propertyRepo.Delete(ctx, id)
}
