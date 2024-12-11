package repository

import (
	"context"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/im"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/scan"
	"scootin/internal/models"
)

type ScooterEventRepository struct {
	DB *bob.DB
}

func NewScooterEventRepository(db *bob.DB) ScooterEventRepositoryInterface {
	return &ScooterEventRepository{DB: db}
}

func (r ScooterEventRepository) Insert(ctx context.Context, scooterEvent *models.ScooterEvent) error {
	stmt := mysql.Insert(
		im.Into("scooter_event"),
		im.Values(mysql.Arg(scooterEvent.ScooterIdentifier, scooterEvent.Event, scooterEvent.Timestamp, scooterEvent.Latitude, scooterEvent.Longitude)),
	)
	_, err := bob.Exec(ctx, r.DB, stmt)
	return err
}

func (r ScooterEventRepository) FindOneByIdentifier(ctx context.Context, scooterEventIdentifier string) (*models.ScooterEvent, error) {
	stmt := mysql.Select(
		sm.From("scooter_event"),
		sm.Where(mysql.Quote("identifier").EQ(mysql.Arg(scooterEventIdentifier))),
	)
	return bob.One(ctx, r.DB, stmt, scan.StructMapper[*models.ScooterEvent]())
}
