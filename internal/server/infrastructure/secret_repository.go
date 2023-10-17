package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/kripsy/GophKeeper/internal/models"
	"github.com/kripsy/GophKeeper/internal/server/entity"
	"go.uber.org/zap"
)

type SecretRepository interface {
	SaveSecret(ctx context.Context, secret entity.Secret) (int, error)
	GetSecret(ctx context.Context, secretID int) (entity.Secret, error)
	DeleteSecret(ctx context.Context, secretID int) error
}

type secretRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func (r *secretRepository) SaveSecret(ctx context.Context, secret entity.Secret) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("failed to Begin Tx in SaveSecret", zap.Error(err))
		return 0, err
	}
	defer func(tx *sql.Tx) {
		err := tx.Rollback()
		if err != nil {
			r.logger.Error("Error tx.Rollback()", zap.Error(err))
		}
	}(tx)

	queryBuilder := squirrel.
		Insert("secrets").Columns("id", "type", "data", "meta").
		Values(secret.ID, secret.Type, secret.Data, secret.Meta).PlaceholderFormat(squirrel.Dollar)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("failed to build sql in SaveSecret", zap.Error(err))
		return 0, err
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		r.logger.Error("failed to exec sql in SaveSecret", zap.Error(err))
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return secret.ID, nil
}

func (r *secretRepository) GetSecret(ctx context.Context, secretID int) (entity.Secret, error) {
	var secret entity.Secret
	queryBuilder := squirrel.Select("id", "type", "data", "meta").
		From("secrets").
		Where(squirrel.Eq{"id": secretID}).
		PlaceholderFormat(squirrel.Dollar)

	qbsql, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("failed to build sql in GetSecret", zap.Error(err))
		return secret, err
	}

	row := r.db.QueryRowContext(ctx, qbsql, args...)
	err = row.Scan(&secret.ID, &secret.Type, &secret.Data, &secret.Meta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return secret, models.NewSecretNotFoundError(secretID)
		}
		r.logger.Error("failed to scan secret", zap.Error(err))
		return secret, err
	}
	return secret, nil
}

func (r *secretRepository) DeleteSecret(ctx context.Context, secretID int) error {
	queryBuilder := squirrel.Delete("secrets").
		Where(squirrel.Eq{"id": secretID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("failed to build sql in DeleteSecret", zap.Error(err))
		return err
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		r.logger.Error("failed to exec sql in DeleteSecret", zap.Error(err))
		return err
	}
	return nil
}

func NewSecretRepository(repo *repository) (SecretRepository, error) {
	return &secretRepository{
		db:     repo.db,
		logger: repo.logger,
	}, nil
}
