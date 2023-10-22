package infrastructure

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"

// 	"github.com/Masterminds/squirrel"
// 	"github.com/kripsy/GophKeeper/internal/models"
// 	"github.com/kripsy/GophKeeper/internal/server/entity"
// 	"go.uber.org/zap"
// )

// type secretRepository struct {
// 	db     *sql.DB
// 	logger *zap.Logger
// }

// func (r *secretRepository) SaveSecret(ctx context.Context, secret entity.Secret) (int, error) {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		r.logger.Error("failed to Begin Tx in SaveSecret", zap.Error(err))
// 		return 0, err
// 	}
// 	defer func(tx *sql.Tx) {
// 		err := tx.Rollback()
// 		if err != nil {
// 			r.logger.Error("Error tx.Rollback()", zap.Error(err))
// 		}
// 	}(tx)

// 	queryBuilder := squirrel.
// 		Insert("secrets").
// 		Columns("type", "data", "meta", "chunk_num", "total_chunks", "user_id").
// 		Values(secret.Type, secret.Data, secret.Meta, secret.ChunkNum, secret.TotalChunks, secret.UserID).
// 		Suffix("RETURNING id").
// 		PlaceholderFormat(squirrel.Dollar)

// 	sql, args, err := queryBuilder.ToSql()
// 	if err != nil {
// 		r.logger.Error("failed to build sql in SaveSecret", zap.Error(err))
// 		return 0, err
// 	}

// 	row := tx.QueryRowContext(ctx, sql, args...)
// 	var id int
// 	err = row.Scan(&id)
// 	if err != nil {
// 		r.logger.Error("failed to scan id in SaveSecret", zap.Error(err))
// 		return 0, err
// 	}

// 	err = tx.Commit()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return id, nil
// }

// func (r *secretRepository) GetSecretsByUserID(ctx context.Context, userID int, limit int, offset int) ([]entity.Secret, error) {
// 	var secrets []entity.Secret

// 	queryBuilder := squirrel.
// 		Select("id", "type", "meta", "chunk_num", "total_chunks", "user_id").
// 		From("secrets").
// 		Where(squirrel.Eq{"user_id": userID}).
// 		Limit(uint64(limit)).
// 		Offset(uint64(offset)).
// 		OrderBy("id DESC").
// 		PlaceholderFormat(squirrel.Dollar)

// 	sql, args, err := queryBuilder.ToSql()
// 	if err != nil {
// 		r.logger.Error("failed to build sql in GetSecretsByUserID", zap.Error(err))

// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	rows, err := r.db.QueryContext(ctx, sql, args...)
// 	if err != nil {
// 		r.logger.Error("failed to query secrets in GetSecretsByUserID", zap.Error(err))

// 		return nil, fmt.Errorf("%w", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var secret entity.Secret
// 		err := rows.Scan(&secret.ID, &secret.Type, &secret.Meta, &secret.ChunkNum, &secret.TotalChunks, &secret.UserID)
// 		if err != nil {
// 			r.logger.Error("failed to scan secret in GetSecretsByUserID", zap.Error(err))

// 			return nil, fmt.Errorf("%w", err)
// 		}
// 		secrets = append(secrets, secret)
// 	}

// 	if err = rows.Err(); err != nil {
// 		r.logger.Error("rows error in GetSecretsByUserID", zap.Error(err))

// 		return nil, fmt.Errorf("%w", err)
// 	}

// 	return secrets, nil
// }

// func (r *secretRepository) GetSecretByID(ctx context.Context, secretID, userID int) (entity.Secret, error) {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		r.logger.Error("Failed to begin transaction", zap.Error(err))

// 		return entity.Secret{}, err
// 	}
// 	defer func() {
// 		if err := tx.Rollback(); err != nil {
// 			r.logger.Error("Failed to rollback transaction", zap.Error(err))
// 		}
// 	}()

// 	queryBuilder := squirrel.
// 		Select("id", "type", "data", "meta", "chunk_num", "total_chunks", "user_id").
// 		From("secrets").
// 		Where(squirrel.Eq{"id": secretID, "user_id": userID}).
// 		RunWith(tx).
// 		PlaceholderFormat(squirrel.Dollar)

// 	qbsql, args, err := queryBuilder.ToSql()
// 	if err != nil {
// 		r.logger.Error("Failed to build SQL query", zap.Error(err))

// 		return entity.Secret{}, err
// 	}

// 	row := tx.QueryRow(qbsql, args...)

// 	var secret entity.Secret
// 	err = row.Scan(&secret.ID, &secret.Type, &secret.Data, &secret.Meta, &secret.ChunkNum, &secret.TotalChunks, &secret.UserID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return entity.Secret{}, models.NewSecretNotFoundError(secretID)
// 		}
// 		r.logger.Error("Failed to retrieve secret from DB", zap.Error(err))

// 		return entity.Secret{}, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		r.logger.Error("Failed to commit transaction", zap.Error(err))

// 		return entity.Secret{}, err
// 	}

// 	return secret, nil
// }

// func (r *secretRepository) DeleteSecret(ctx context.Context, secretID int) error {
// 	queryBuilder := squirrel.Delete("secrets").
// 		Where(squirrel.Eq{"id": secretID}).
// 		PlaceholderFormat(squirrel.Dollar)

// 	sql, args, err := queryBuilder.ToSql()
// 	if err != nil {
// 		r.logger.Error("failed to build sql in DeleteSecret", zap.Error(err))
// 		return err
// 	}

// 	_, err = r.db.ExecContext(ctx, sql, args...)
// 	if err != nil {
// 		r.logger.Error("failed to exec sql in DeleteSecret", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }

// func NewSecretRepository(repo *repository) (*secretRepository, error) {
// 	return &secretRepository{
// 		db:     repo.db,
// 		logger: repo.logger,
// 	}, nil
// }
