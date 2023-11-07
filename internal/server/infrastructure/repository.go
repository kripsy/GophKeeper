package infrastructure

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	//nolint:revive,nolintlint
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	//nolint:revive,nolintlint
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	//nolint:revive,nolintlint
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const MAXOPENCONN = 25
const MAXIDLECONNS = 10
const CONNMAXLIFETIME = time.Minute * 30

type Repository interface {
	// Close close connection with repo or return error.
	Close() error
}

type repository struct {
	// database instance.
	db     *sql.DB
	logger *zap.Logger
}

//nolint:revive,nolintlint
func InitNewRepository(connString string, logger *zap.Logger) (*repository, error) {
	err := RunMigrations(context.Background(), connString, logger)
	if err != nil {
		logger.Debug("Fail to run migrations", zap.String("msg", err.Error()))

		return nil, err
	}

	db, err := sql.Open("pgx", connString)

	// max open connections.
	db.SetMaxOpenConns(MAXOPENCONN)

	// max unused connections.
	db.SetMaxIdleConns(MAXIDLECONNS)

	// max time lifetime connection.
	db.SetConnMaxLifetime(CONNMAXLIFETIME)
	if err != nil {
		logger.Debug("Fail open db connection", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	logger.Debug("initDB", zap.String("connString", connString))
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		logger.Debug("Fail to ping db", zap.String("msg", err.Error()))

		return nil, fmt.Errorf("%w", err)
	}

	return &repository{
		db:     db,
		logger: logger,
	}, nil
}

//go:embed migrations/*.sql
var fs embed.FS

//nolint:revive,nolintlint
func RunMigrations(ctx context.Context, connString string, logger *zap.Logger) error {
	d, err := iofs.New(fs, "migrations")
	if err != nil {
		logger.Error("error create new iofs", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}
	logger.Debug("Success create iofs")
	m, err := migrate.NewWithSourceInstance("iofs", d, connString)
	if err != nil {
		logger.Error("error create NewWithSourceInstance", zap.String("msg", err.Error()))

		return fmt.Errorf("%w", err)
	}
	logger.Debug("Success create NewWithSourceInstance")

	if err = m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			logger.Error("error run migrations", zap.String("msg", err.Error()))

			return fmt.Errorf("%w", err)
		}
	}

	return nil
}

func (r *repository) Close() error {
	//nolint:wrapcheck
	return r.db.Close()
}
