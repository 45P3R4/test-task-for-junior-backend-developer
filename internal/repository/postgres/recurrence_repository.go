package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	recurrencedomain "example.com/taskservice/internal/domain/recurrence"
	recurrenceusecase "example.com/taskservice/internal/usecase/recurrence"
)

type RecurrenceRepository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *RecurrenceRepository {
	return &RecurrenceRepository{pool: pool}
}

func (r *RecurrenceRepository) Create(ctx context.Context, rule *recurrencedomain.RecurrenceRule) (*recurrencedomain.RecurrenceRule, error) {
	query := `
		INSERT INTO recurrence_rules (task_id, recurrence_type, recurrence_modifiers, end_date, max_occurrences, interval, days, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING task_id, recurrence_type, recurrence_modifiers, end_date, max_occurrences, interval, days, created_at, updated_at`

	row := r.pool.QueryRow(ctx, query,
		rule.TaskID,
		string(rule.RecurrenceType),
		rule.RecurrenceModifiers,
		rule.EndDate,
		rule.MaxOccurrences,
		rule.Interval,
		rule.Days,
		rule.CreatedAt,
		rule.UpdatedAt,
	)

	created, err := scanRecurrenceRule(row)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (r *RecurrenceRepository) GetByTaskID(ctx context.Context, taskID int64) (*recurrencedomain.RecurrenceRule, error) {
	query := `
		SELECT task_id, recurrence_type, recurrence_modifiers, end_date, max_occurrences, interval, days, created_at, updated_at
		FROM recurrence_rules
		WHERE task_id = $1`

	row := r.pool.QueryRow(ctx, query, taskID)
	found, err := scanRecurrenceRule(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, recurrenceusecase.ErrNotFound
		}
		return nil, err
	}

	return found, nil
}

func (r *RecurrenceRepository) Update(ctx context.Context, rule *recurrencedomain.RecurrenceRule) (*recurrencedomain.RecurrenceRule, error) {
	query := `
		UPDATE recurrence_rules
		SET recurrence_type = $1,
			recurrence_modifiers = $2,
			end_date = $3,
			max_occurrences = $4,
			interval = $5,
			days = $6,
			updated_at = $7
		WHERE task_id = $8
		RETURNING task_id, recurrence_type, recurrence_modifiers, end_date, max_occurrences, interval, days, created_at, updated_at`

	row := r.pool.QueryRow(ctx, query,
		string(rule.RecurrenceType),
		rule.RecurrenceModifiers,
		rule.EndDate,
		rule.MaxOccurrences,
		rule.Interval,
		rule.Days,
		rule.UpdatedAt,
		rule.TaskID,
	)

	updated, err := scanRecurrenceRule(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, recurrenceusecase.ErrNotFound
		}
		return nil, err
	}

	return updated, nil
}

func (r *RecurrenceRepository) Delete(ctx context.Context, taskID int64) error {
	query := `DELETE FROM recurrence_rules WHERE task_id = $1`

	result, err := r.pool.Exec(ctx, query, taskID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return recurrenceusecase.ErrDelete
	}

	return nil
}

type recurrenceScanner interface {
	Scan(dest ...any) error
}

func scanRecurrenceRule(scanner recurrenceScanner) (*recurrencedomain.RecurrenceRule, error) {
	var (
		rule           recurrencedomain.RecurrenceRule
		recurrenceType string
		modifiers      []string
		endDate        *time.Time
	)

	if err := scanner.Scan(
		&rule.TaskID,
		&recurrenceType,
		&modifiers,
		&endDate,
		&rule.MaxOccurrences,
		&rule.Interval,
		&rule.Days,
		&rule.CreatedAt,
		&rule.UpdatedAt,
	); err != nil {
		return nil, err
	}

	rule.RecurrenceType = recurrencedomain.RecurrenceType(recurrenceType)
	rule.RecurrenceModifiers = make([]recurrencedomain.RecurrenceModifier, len(modifiers))

	for i, m := range modifiers {
		rule.RecurrenceModifiers[i] = recurrencedomain.RecurrenceModifier(m)
	}

	rule.EndDate = endDate

	return &rule, nil
}
