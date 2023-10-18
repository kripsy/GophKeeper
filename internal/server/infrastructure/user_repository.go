package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"github.com/kripsy/GophKeeper/internal/utils"

	_ "github.com/jackc/pgx/v5/stdlib"

	"go.uber.org/zap"
)

type userRepository struct {
	// database instance.
	db     *sql.DB
	logger *zap.Logger
}

func (r *userRepository) RegisterUser(ctx context.Context, user entity.User) (int, error) {

	userExist, err := r.isUserExists(ctx, user.Username)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	if userExist {
		return 0, models.NewUserExistsError(user.Username)
	}

	tx, err := r.db.Begin()
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("Error tx.Rollback()", zap.Error(err))
		}
	}(tx)

	if err != nil {
		r.logger.Error("failed to Begin Tx in RegisterUser", zap.Error(err))

		return 0, fmt.Errorf("%w", err)
	}

	id, err := r.getNextRowID(ctx, "users")
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	hashPassword, err := utils.GetHash(ctx, user.Password, r.logger)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	quryBuilder := squirrel.
		Insert("users").Columns("id", "username", "password").
		Values(id, user.Username, hashPassword).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := quryBuilder.ToSql()

	if err != nil {
		r.logger.Error("failed to build sql in RegisterUser", zap.String("msg", err.Error()))

		return 0, fmt.Errorf("%w", err)
	}

	r.logger.Debug("success build sql", zap.String("msg", sql))

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		r.logger.Error("failed to exec sql in RegisterUser", zap.String("msg", err.Error()))

		return 0, fmt.Errorf("%w", err)
	}

	err = tx.Commit()
	if err != nil {

		return 0, fmt.Errorf("%w", err)
	}
	r.logger.Debug("success commit RegisterUser")

	return id, nil
}

func (r *userRepository) LoginUser(ctx context.Context, user entity.User) (int, error) {
	userExist, err := r.isUserExists(ctx, user.Username)
	if err != nil {
		return 0, fmt.Errorf("%w", err)
	}

	if !userExist {
		return 0, models.NewUserExistsError(user.Username)
	}

	userID, hashPassword, err := r.getUserHashPassword(ctx, user.Username)
	if err != nil {
		r.logger.Error("error getUserHashPassword in LoginUser", zap.Error(err))

		return 0, fmt.Errorf("%w", err)
	}

	err = utils.IsPasswordCorrect(ctx, []byte(user.Password), []byte(hashPassword), r.logger)
	if err != nil {
		r.logger.Error("password incorrect", zap.Error(err))

		return 0, models.NewUserLoginError(user.Username)
	}

	return userID, nil
}

func (r *userRepository) isUserExists(ctx context.Context, username string) (bool, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	r.logger.Debug("start IsUserExists")

	tx, err := r.db.Begin()
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("Error tx.Rollback()", zap.Error(err))
		}
	}(tx)

	if err != nil {
		r.logger.Error("failed to Begin Tx in IsUserExists", zap.Error(err))

		return false, err
	}

	var userExists bool
	queryBuilder := squirrel.Select("1").
		Prefix("SELECT EXISTS (").
		From("users").
		Where(squirrel.Eq{"username": username}).
		Suffix(")").
		PlaceholderFormat(squirrel.Dollar)
	sql, args, err := queryBuilder.ToSql()

	if err != nil {
		r.logger.Error("failed to build sql in IsUserExists", zap.Error(err))

		return false, err
	}

	r.logger.Debug("success build sql", zap.String("msg", sql))

	row := tx.QueryRowContext(ctx, sql, args...)

	err = row.Scan(&userExists)

	if err != nil {
		r.logger.Error("failed scan userExists", zap.Error(err))

		return false, err
	}
	r.logger.Debug("success scan userExists, value ->", zap.Bool("msg", userExists))
	err = tx.Commit()
	if err != nil {
		r.logger.Error("Error tx.Commit()", zap.Error(err))

		return false, err
	}

	return userExists, nil
}

func (r *userRepository) getNextRowID(ctx context.Context, tableName string) (int, error) {

	ctx, canlcel := context.WithTimeout(ctx, time.Second)
	defer canlcel()
	r.logger.Debug("start getNextRowID")

	tx, err := r.db.Begin()
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
	}(tx)

	if err != nil {
		r.logger.Error("failed to Begin Tx in getNextRowID", zap.String("msg", err.Error()))

		return 0, fmt.Errorf("%w", err)
	}

	queryBuilder := squirrel.
		Select("MAX(id)+1").
		From(tableName)

	qbsql, _, err := queryBuilder.ToSql()

	if err != nil {
		r.logger.Error("failed to build sql in getNextRowID", zap.String("msg", err.Error()))

		return 0, fmt.Errorf("%w", err)
	}

	r.logger.Debug("success build sql", zap.String("msg", qbsql))

	stmt, err := tx.PrepareContext(ctx, qbsql)
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			r.logger.Error("Unable PrepareContext: ", zap.String("msg", err.Error()))
		}
	}(stmt)
	if err != nil {
		r.logger.Error("failed to PrepareContext stmt in getNextRowID", zap.String("msg", err.Error()))

		return 0, fmt.Errorf("%w", err)
	}

	row := stmt.QueryRowContext(ctx)
	var nextID sql.NullInt32
	err = row.Scan(&nextID)
	if err != nil {
		r.logger.Error("failed to scan getNextRowID", zap.String("msg", err.Error()))

		return 0, fmt.Errorf("%w", err)
	}

	err = tx.Commit()
	if err != nil {
		r.logger.Error("Error tx.Rollback()", zap.String("msg", err.Error()))

		return 0, fmt.Errorf("%w", err)
	}
	r.logger.Debug("success commit getNextRowID")
	if !nextID.Valid {
		return 1, nil
	}

	return int(nextID.Int32), nil
}

func (r *userRepository) getUserHashPassword(ctx context.Context, username string) (int, string, error) {
	r.logger.Debug("start GetUserHashPassword")

	tx, err := r.db.Begin()
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("Error tx.Rollback()", zap.String("msg", err.Error()))
		}
	}(tx)

	if err != nil {
		r.logger.Error("failed to Begin Tx in GetUserHashPassword", zap.String("msg", err.Error()))

		return 0, "", fmt.Errorf("%w", err)
	}

	var userID int
	var hashPassword string
	queryBuilder := squirrel.Select("id, password").
		From("users").
		Where(squirrel.Eq{"username": username}).
		PlaceholderFormat(squirrel.Dollar)
	bsql, args, err := queryBuilder.ToSql()

	if err != nil {
		r.logger.Error("failed to build sql in GetUserHashPassword", zap.String("msg", err.Error()))

		return 0, "", fmt.Errorf("%w", err)
	}

	r.logger.Debug("success build sql in GetUserHashPassword", zap.String("msg", bsql))

	row := tx.QueryRowContext(ctx, bsql, args...)

	err = row.Scan(&userID, &hashPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("error compare username and pwd", zap.String("msg", username))

			return 0, "", models.NewUserLoginError(username)
		}
		r.logger.Error("failed scan userExists", zap.String("msg", err.Error()))

		return 0, "", fmt.Errorf("%w", err)
	}

	r.logger.Debug("success get hash password ->", zap.String("msg", hashPassword))
	err = tx.Commit()
	if err != nil {

		return 0, "", fmt.Errorf("%w", err)
	}

	return userID, hashPassword, nil
}

func NewUserRepository(repo *repository) (*userRepository, error) {
	return &userRepository{
		db:     repo.db,
		logger: repo.logger,
	}, nil
}
