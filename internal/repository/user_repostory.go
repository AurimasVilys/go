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

type UserRepository struct {
	DB *bob.DB
}

func NewUserRepository(db *bob.DB) UserRepositoryInterface {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Insert(ctx context.Context, scooter *models.User) error {
	stmt := mysql.Insert(
		im.Into("user", "identifier"),
		im.Values(mysql.Arg(scooter.Identifier)),
	)
	_, err := bob.Exec(ctx, r.DB, stmt)
	return err
}

func (r *UserRepository) FindByIdentifier(ctx context.Context, userIdentifier string) (*models.User, error) {
	stmt := mysql.Select(
		sm.From("user"),
		sm.Where(mysql.Quote("identifier").EQ(mysql.Arg(userIdentifier))),
	)
	return bob.One(ctx, r.DB, stmt, scan.StructMapper[*models.User]())
}
