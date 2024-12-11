package repository

import (
	"context"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/mysql"
	"github.com/stephenafamo/bob/dialect/mysql/im"
	"github.com/stephenafamo/bob/dialect/mysql/sm"
	"github.com/stephenafamo/bob/dialect/mysql/um"
	"github.com/stephenafamo/scan"
	"scootin/internal/models"
)

type ScooterRepository struct {
	DB *bob.DB
}

func NewScooterRepository(db *bob.DB) ScooterRepositoryInterface {
	return &ScooterRepository{DB: db}
}

func (r *ScooterRepository) Insert(ctx context.Context, scooter *models.Scooter) error {
	stmt := mysql.Insert(
		im.Into("scooter", "identifier", "last_confirmed_latitude", "last_confirmed_longitude"),
		im.Values(mysql.Arg(scooter.Identifier), mysql.Arg(scooter.LastConfirmedLatitude), mysql.Arg(scooter.LastConfirmedLongitude)),
	)
	_, err := bob.Exec(ctx, r.DB, stmt)
	return err
}

func (r *ScooterRepository) FindOneByIdentifier(ctx context.Context, scooterIdentifier string) (*models.Scooter, error) {
	stmt := mysql.Select(
		sm.From("scooter"),
		sm.Where(mysql.Quote("identifier").EQ(mysql.Arg(scooterIdentifier))),
	)
	return bob.One(ctx, r.DB, stmt, scan.StructMapper[*models.Scooter]())
}

func (r *ScooterRepository) UpdateOccupiedByUser(ctx context.Context, scooterIdentifier, userIdentifier string) error {
	stmt := mysql.Update(
		um.Table("scooter"),
		um.SetCol("occupied_user_identifier").ToArg(userIdentifier),
		um.Where(mysql.Quote("identifier").EQ(mysql.Arg(scooterIdentifier))),
	)
	_, err := bob.Exec(ctx, r.DB, stmt)
	return err
}

func (r *ScooterRepository) ReleaseOccupied(ctx context.Context, scooterIdentifier string) error {
	stmt := mysql.Update(
		um.Table("scooter"),
		um.SetCol("occupied_user_identifier").ToArg(nil),
		um.Where(mysql.Quote("identifier").EQ(mysql.Arg(scooterIdentifier))),
	)
	_, err := bob.Exec(ctx, r.DB, stmt)
	return err
}

func (r *ScooterRepository) FindByCoordinatesAndStatus(ctx context.Context, startLatitude, endLatitude, startLongitude, endLongitude string, includeOccupied bool) ([]*models.Scooter, error) {
	stmt := mysql.Select(
		sm.Columns("identifier", "last_confirmed_latitude", "last_confirmed_longitude"),
		sm.From("scooter"),
		sm.Where(mysql.Quote("last_confirmed_latitude").Between(mysql.Arg(startLatitude), mysql.Arg(endLatitude))),
		sm.Where(mysql.Quote("last_confirmed_longitude").Between(mysql.Arg(startLongitude), mysql.Arg(endLongitude))),
	)

	if !includeOccupied {
		stmt.Expression.Where.AppendWhere(mysql.Quote("occupied_user_identifier").IsNull())
	}

	return bob.All(ctx, r.DB, stmt, scan.StructMapper[*models.Scooter]())
}

func (r *ScooterRepository) Update(ctx context.Context, scooter *models.Scooter) error {
	stmt := mysql.Update(
		um.Table("scooter"),
		um.SetCol("last_confirmed_latitude").ToArg(scooter.LastConfirmedLatitude),
		um.SetCol("last_confirmed_longitude").ToArg(scooter.LastConfirmedLongitude),
		um.SetCol("occupied_user_identifier").ToArg(scooter.OccupiedUserIdentifier),
		um.Where(mysql.Quote("identifier").EQ(mysql.Arg(scooter.Identifier))),
	)
	_, err := bob.Exec(ctx, r.DB, stmt)
	return err
}
