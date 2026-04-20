package recurrence

import "errors"

var (
	ErrInvalidInput          = errors.New("invalid recurrence input")
	ErrInvalidId             = errors.New("task_id must be positive")
	ErrNotFound              = errors.New("recurrence not found")
	ErrDelete                = errors.New("failed to delete recurrence")
	ErrInvalidRecurrenceRule = errors.New("invalid recurrence rule")
	ErrTaskNotRecurring      = errors.New("task not recurring")

	ErrWeekly            = errors.New("days must be between 0 and 6 (Sunday=0, Saturday=6)")
	ErrMonthly           = errors.New("days must be between 1 and 31")
	ErrEmptyType         = errors.New("invalid recurrence_type")
	ErrInvalidType       = errors.New("recurrence_type is required")
	ErrOppositeModifiers = errors.New("rule must be without opposing modifiers")
)
