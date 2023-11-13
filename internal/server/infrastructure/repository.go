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

// Maximum number of open connections to the database.
const MAXOPENCONN = 25

// Maximum number of idle connections in the pool.
const MAXIDLECONNS = 10

// Maximum lifetime of a connection.
const CONNMAXLIFETIME = time.Minute * 30

// Repository interface defines the standard contract for database operations.
// It includes a method to close the database connection.
type Repository interface {
	// Close terminates the connection with the database and returns an error if unsuccessful.
	Close() error
}

type repository struct {
	// database instance.
	db     *sql.DB
	logger *zap.Logger
}

// InitNewRepository initializes a new repository with the given database connection string.
// It runs necessary database migrations and sets up the database connection.
// connString: Database connection string.
// logger: Logger for logging purposes.
// Returns a pointer to the repository or an error if initialization fails.
//
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

// RunMigrations applies database migrations using embedded SQL files.
// It ensures the database schema is up to date.
// ctx: Context for execution.
// connString: Database connection string.
// logger: Logger for logging migration process.
// Returns an error if migration fails.
//
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

// Close gracefully terminates the database connection.
// Returns an error if closing the connection fails.
func (r *repository) Close() error {
	//nolint:wrapcheck
	return r.db.Close()
}
