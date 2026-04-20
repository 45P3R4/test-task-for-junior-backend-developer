package recurrence

import (
	"time"
)

type RecurrenceType string

const (
	RecurrenceDaily      RecurrenceType = "daily"
	RecurrenceWeekly     RecurrenceType = "weekly"
	RecurrenceMonthly    RecurrenceType = "monthly"
	RecurrenceEndOfMonth RecurrenceType = "monthend"
)

type RecurrenceModifier string

const (
	ModifierOdd      RecurrenceModifier = "oddly"
	ModifierEven     RecurrenceModifier = "evenly"
	ModifierWeekdays RecurrenceModifier = "weekdays"
	ModifierWeekends RecurrenceModifier = "weekends"
)

type RecurrenceRule struct {
	TaskID              int64                `json:"task_id"`
	RecurrenceType      RecurrenceType       `json:"recurrence_type"`
	RecurrenceModifiers []RecurrenceModifier `json:"recurrence_modifiers"`
	EndDate             *time.Time           `json:"end_date,omitempty"`
	MaxOccurrences      *int                 `json:"max_occurrences,omitempty"`
	Interval            *int                 `json:"interval,omitempty"` // Every N days
	Days                []int                `json:"days,omitempty"`     // Month date type = monthly, week day type = weekly
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
}

func (r RecurrenceType) Valid() bool {
	switch r {
	case RecurrenceDaily, RecurrenceWeekly, RecurrenceMonthly, RecurrenceEndOfMonth:
		return true
	default:
		return false
	}
}
