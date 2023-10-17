package infrastructure

import (
	"database/sql"
	"errors"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kripsy/GophKeeper/internal/entity"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"
)

type UserRepository interface {
	RegisterUser(user *entity.User) (int, error)
	LoginUser() error
}

type userRepository struct {
	// database instance.
	db     *sql.DB
	logger *zap.Logger
}

func (r *userRepository) RegisterUser(user *entity.User) (int, error) {
	return 0, errors.New("not implemented")
	// ctx, canlcel := context.WithTimeout(context.Background(), time.Second)
	// defer canlcel()
	// r.logger.Debug("start RegisterUser")

	// tx, err := r.db.Begin()

	// if err != nil {
	// 	r.logger.Error("failed to Begin Tx in RegisterUser", zap.String("msg", err.Error()))
	// 	return 0, err
	// }
	// defer tx.Rollback()

	// queryBuilder := squirrel.
	// 	Insert("devices").
	// 	Columns("name", "description").
	// 	Values(user.Username, user.Password).
	// 	Suffix("RETURNING id").
	// 	PlaceholderFormat(squirrel.Dollar)

	// sql, args, err := queryBuilder.ToSql()

	// if err != nil {
	// 	r.logger.Error("failed to build sql in RegisterNewDevice", zap.String("msg", err.Error()))
	// 	return 0, err
	// }

	// r.logger.Debug("success build sql", zap.String("msg", sql))

	// var deviceID int
	// row := tx.QueryRowContext(ctx, sql, args...)
	// err = row.Scan(&deviceID)
	// if err != nil {
	// 	r.logger.Error("failed to exec sql in RegisterUser", zap.String("msg", err.Error()))
	// 	return 0, err
	// }

	// err = tx.Commit()
	// if err != nil {
	// 	return 0, err
	// }
	// r.logger.Debug("success commit RegisterNewDevice")
	// return deviceID, nil
}

func (r *userRepository) LoginUser() error {
	return errors.New("not implemented")
}

func NewUserRepository(repo *repository) (UserRepository, error) {
	return &userRepository{
		db:     repo.db,
		logger: repo.logger,
	}, nil
}
