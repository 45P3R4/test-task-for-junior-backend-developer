package handlers

import (
	"time"

	recurrencedomain "example.com/taskservice/internal/domain/recurrence"
	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskMutationDTO struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
}

type taskDTO struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      taskdomain.Status `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	return taskDTO{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

type recurrenceMutationDTO struct {
	TaskID              int64                                 `json:"task_id"`
	RecurrenceType      recurrencedomain.RecurrenceType       `json:"recurrence_type"`
	RecurrenceModifiers []recurrencedomain.RecurrenceModifier `json:"recurrence_modifiers,omitempty"`
	EndDate             *time.Time                            `json:"end_date,omitempty"`
	MaxOccurrences      *int                                  `json:"max_occurrences,omitempty"`
	Interval            *int                                  `json:"interval,omitempty"`
	Days                []int                                 `json:"days,omitempty"`
}

type recurrenceRuleDTO struct {
	ID                  int64                                 `json:"id"`
	TaskID              int64                                 `json:"task_id"`
	RecurrenceType      recurrencedomain.RecurrenceType       `json:"recurrence_type"`
	RecurrenceModifiers []recurrencedomain.RecurrenceModifier `json:"recurrence_modifiers,omitempty"`
	EndDate             *time.Time                            `json:"end_date,omitempty"`
	MaxOccurrences      *int                                  `json:"max_occurrences,omitempty"`
	Interval            *int                                  `json:"interval"`
	Days                []int                                 `json:"days,omitempty"`
	CreatedAt           time.Time                             `json:"created_at"`
	UpdatedAt           time.Time                             `json:"updated_at"`
}

func newRecurrenceRuleDTO(rule *recurrencedomain.RecurrenceRule) recurrenceRuleDTO {
	return recurrenceRuleDTO{
		TaskID:         rule.TaskID,
		RecurrenceType: rule.RecurrenceType,
		Interval:       rule.Interval,
		Days:           rule.Days,
	}
}
