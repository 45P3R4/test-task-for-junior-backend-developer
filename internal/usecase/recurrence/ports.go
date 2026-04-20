package recurrence

import (
	"context"
	"time"

	recurrencedomain "example.com/taskservice/internal/domain/recurrence"
)

type Repository interface {
	Create(ctx context.Context, rule *recurrencedomain.RecurrenceRule) (*recurrencedomain.RecurrenceRule, error)
	GetByTaskID(ctx context.Context, id int64) (*recurrencedomain.RecurrenceRule, error)
	Update(ctx context.Context, rule *recurrencedomain.RecurrenceRule) (*recurrencedomain.RecurrenceRule, error)
	Delete(ctx context.Context, taskID int64) error
}

type Usecase interface {
	Create(ctx context.Context, input CreateUpdateInput, taskID int64) (*recurrencedomain.RecurrenceRule, error)
	GetByTaskID(ctx context.Context, taskID int64) (*recurrencedomain.RecurrenceRule, error)
	Update(ctx context.Context, id int64, input CreateUpdateInput) (*recurrencedomain.RecurrenceRule, error)
	Delete(ctx context.Context, id int64) error
}

type RecurrenceInput interface {
	GetRecurrenceType() recurrencedomain.RecurrenceType
	GetDays() []int
}

type CreateUpdateInput struct {
	TaskID              int64
	RecurrenceType      recurrencedomain.RecurrenceType
	RecurrenceModifiers []recurrencedomain.RecurrenceModifier
	EndDate             *time.Time
	MaxOccurrences      *int
	Interval            *int
	Days                []int
}
