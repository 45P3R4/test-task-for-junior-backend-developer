package recurrence

import "errors"

var (
	ErrInvalidInput          = errors.New("invalid recurrence input")
	ErrNotFound              = errors.New("Recurrence not found")
	ErrDelete                = errors.New("Failed to delete recurrence")
	ErrInvalidRecurrenceRule = errors.New("Invalid recurrence rule")
	ErrTaskNotRecurring      = errors.New("Task not recurring")

	ErrWeekly            = errors.New("Days must be between 0 and 6 (Sunday=0, Saturday=6)")
	ErrMonthly           = errors.New("days must be between 1 and 31")
	ErrEmptyType         = errors.New("Invalid recurrence_type")
	ErrInvalidType       = errors.New("Recurrence_type is required")
	ErrOppositeModifiers = errors.New("Rule must be without opposing modifiers")
)
