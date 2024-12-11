package repository

import (
	"context"
	"scootin/internal/models"
)

type UserRepositoryInterface interface {
	FindByIdentifier(ctx context.Context, userIdentifier string) (*models.User, error)
	Insert(ctx context.Context, user *models.User) error
}

type ScooterRepositoryInterface interface {
	Insert(ctx context.Context, scooter *models.Scooter) error
	FindOneByIdentifier(ctx context.Context, scooterIdentifier string) (*models.Scooter, error)
	UpdateOccupiedByUser(ctx context.Context, scooterIdentifier, userIdentifier string) error
	ReleaseOccupied(ctx context.Context, scooterIdentifier string) error
	FindByCoordinatesAndStatus(ctx context.Context, startLatitude, endLatitude, startLongitude, endLongitude string, includeOccupied bool) ([]*models.Scooter, error)
	Update(ctx context.Context, scooter *models.Scooter) error
}

type ScooterEventRepositoryInterface interface {
	Insert(ctx context.Context, scooterEvent *models.ScooterEvent) error
	FindOneByIdentifier(ctx context.Context, scooterEventIdentifier string) (*models.ScooterEvent, error)
}
