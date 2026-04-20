package recurrence

import (
	"context"
	"fmt"
	"time"

	recurrencedomain "example.com/taskservice/internal/domain/recurrence"
	task "example.com/taskservice/internal/usecase/task"
)

type Service struct {
	repo     Repository
	taskRepo task.Repository
	now      func() time.Time
}

func NewService(repo Repository, taskRepo task.Repository) *Service {
	return &Service{
		repo:     repo,
		taskRepo: taskRepo,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input CreateUpdateInput, taskID int64) (*recurrencedomain.RecurrenceRule, error) {
	normalized, err := s.validateInput(input)
	if err != nil {
		return nil, err
	}

	if _, err := s.taskRepo.GetByID(ctx, taskID); err != nil {
		return nil, fmt.Errorf("%w: task not found", ErrInvalidInput)
	}

	recurrenceType := recurrencedomain.RecurrenceType(normalized.RecurrenceType)
	modifiers := make([]recurrencedomain.RecurrenceModifier, len(normalized.RecurrenceModifiers))
	for i, m := range normalized.RecurrenceModifiers {
		modifiers[i] = recurrencedomain.RecurrenceModifier(m)
	}

	model := &recurrencedomain.RecurrenceRule{
		TaskID:              taskID,
		RecurrenceType:      recurrenceType,
		RecurrenceModifiers: modifiers,
		EndDate:             normalized.EndDate,
		MaxOccurrences:      normalized.MaxOccurrences,
		Interval:            normalized.Interval,
		Days:                normalized.Days,
	}

	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.repo.Create(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByTaskID(ctx context.Context, taskID int64) (*recurrencedomain.RecurrenceRule, error) {
	if taskID <= 0 {
		return nil, fmt.Errorf("%w: task_id must be positive", ErrInvalidInput)
	}

	rule, err := s.repo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *Service) Update(ctx context.Context, taskID int64, input CreateUpdateInput) (*recurrencedomain.RecurrenceRule, error) {

	normalized, err := s.validateInput(input)
	if err != nil {
		return nil, err
	}

	modifiers := make([]recurrencedomain.RecurrenceModifier, 0)

	for _, modifier := range normalized.RecurrenceModifiers {
		modifiers = append(modifiers, recurrencedomain.RecurrenceModifier(modifier))
	}

	model := &recurrencedomain.RecurrenceRule{
		TaskID:              taskID,
		RecurrenceType:      recurrencedomain.RecurrenceType(normalized.RecurrenceType),
		RecurrenceModifiers: modifiers,
		EndDate:             normalized.EndDate,
		MaxOccurrences:      normalized.MaxOccurrences,
		Interval:            normalized.Interval,
		Days:                normalized.Days,
		UpdatedAt:           s.now(),
	}

	updated, err := s.repo.Update(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, taskID int64) error {
	if taskID <= 0 {
		return fmt.Errorf("%w: task_id must be positive", ErrInvalidInput)
	}

	if _, err := s.repo.GetByTaskID(ctx, taskID); err != nil {
		return err
	}

	return s.repo.Delete(ctx, taskID)
}

func (s *Service) validateInput(input CreateUpdateInput) (CreateUpdateInput, error) {

	if input.RecurrenceType == "" {
		return CreateUpdateInput{}, fmt.Errorf("%w: %w", ErrEmptyType, ErrInvalidInput)
	}

	if !input.RecurrenceType.Valid() {
		return CreateUpdateInput{}, fmt.Errorf("%w: %w", ErrInvalidType, ErrInvalidInput)
	}

	if input.RecurrenceType == recurrencedomain.RecurrenceWeekly {
		for _, day := range input.Days {
			if day < 0 || day > 6 {
				return CreateUpdateInput{}, fmt.Errorf("%w: %w", ErrWeekly, ErrInvalidInput)
			}
		}
	}

	if input.RecurrenceType == recurrencedomain.RecurrenceMonthly {
		for _, day := range input.Days {
			if day < 1 || day > 31 {
				return CreateUpdateInput{}, fmt.Errorf("%w: %w", ErrMonthly, ErrInvalidInput)
			}
		}
	}

	if containsOppositeModifiers(
		input.RecurrenceModifiers,
		recurrencedomain.ModifierEven,
		recurrencedomain.ModifierOdd,
	) {
		return CreateUpdateInput{}, fmt.Errorf("%w: %w", ErrOppositeModifiers, ErrInvalidInput)
	}

	if containsOppositeModifiers(
		input.RecurrenceModifiers,
		recurrencedomain.ModifierWeekdays,
		recurrencedomain.ModifierWeekends,
	) {
		return CreateUpdateInput{}, fmt.Errorf("%w: %w", ErrOppositeModifiers, ErrInvalidInput)
	}

	return input, nil
}

func containsOppositeModifiers(modifiers []recurrencedomain.RecurrenceModifier, rule1 recurrencedomain.RecurrenceModifier, rule2 recurrencedomain.RecurrenceModifier) bool {
	hasRule1, hasRule2 := false, false

	for _, modifier := range modifiers {
		switch modifier {
		case rule1:
			hasRule1 = true
		case rule2:
			hasRule2 = true
		}

		if hasRule1 && hasRule2 {
			return true
		}
	}
	return false
}
