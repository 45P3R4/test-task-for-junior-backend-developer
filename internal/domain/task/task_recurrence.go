package task

import (
	"fmt"
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

// Указатели означают, что поля с ними опциональные.
// При инициализации структуры полям без указателей задается значение по умолчанию.
// Например все int поля будут нулями.
// Задача может содержать повторение для нескольких числел\дней недели, поэтому используется массив
type RecurrenceRule struct {
	RecurrenceType     RecurrenceType       `json:"recurrence_type"`
	RecurrenceModifier []RecurrenceModifier `json:"recurrence_modifier"`
	EndDate            *time.Time           `json:"end_date,omitempty"`        // когда остановить
	MaxOccurrences     *int                 `json:"max_occurrences,omitempty"` // макс. количество повторений
	Interval           *int                 `json:"interval,omitempty"`        // каждый N-й день, 1 = каждый день
	Days               []int                `json:"days,omitempty"`            // каждое число месяца если type = monthly, недели если type = weekly
}

// Паттерн Builder позволит комбинировать правила
// Например каждый четный будний день
type RecurrenceRuleBuilder struct {
	rule *RecurrenceRule
}

func NewRecurrenceBuilder(recurrenceType RecurrenceType) *RecurrenceRuleBuilder {
	return &RecurrenceRuleBuilder{
		rule: &RecurrenceRule{
			RecurrenceType: recurrenceType,
		},
	}
}

func (b *RecurrenceRuleBuilder) MaxOccurrences(max int) (*RecurrenceRuleBuilder, error) {
	if max <= 0 {
		return nil, fmt.Errorf("Количество повторений должно быть больше нуля")
	}
	b.rule.MaxOccurrences = &max
	return b, nil
}

func (b *RecurrenceRuleBuilder) EndDate(endDate time.Time) (*RecurrenceRuleBuilder, error) {
	if endDate.Before(time.Now()) {
		return nil, fmt.Errorf("дата окончания должна быть в будущем")
	}
	b.rule.EndDate = &endDate
	return b, nil
}

func (b *RecurrenceRuleBuilder) Interval(interval int) (*RecurrenceRuleBuilder, error) {
	if interval <= 0 {
		return nil, fmt.Errorf("Интервал должен быть больше нуля")
	}
	b.rule.Interval = &interval
	return b, nil
}

func (b *RecurrenceRuleBuilder) Days(days ...int) (*RecurrenceRuleBuilder, error) {
	for _, day := range days {
		if day < 1 || day > 31 {
			return nil, fmt.Errorf("Число должно быть от 1 до 31")
		}
	}
	b.rule.Days = days
	return b, nil
}

func (b *RecurrenceRuleBuilder) Weekdays() *RecurrenceRuleBuilder {
	b.rule.RecurrenceModifier = append(b.rule.RecurrenceModifier, ModifierWeekdays)
	return b
}

func (b *RecurrenceRuleBuilder) Weekends() *RecurrenceRuleBuilder {
	b.rule.RecurrenceModifier = append(b.rule.RecurrenceModifier, ModifierWeekends)
	return b
}

func (b *RecurrenceRuleBuilder) EvenDays() *RecurrenceRuleBuilder {
	b.rule.RecurrenceModifier = append(b.rule.RecurrenceModifier, ModifierEven)
	return b
}

func (b *RecurrenceRuleBuilder) OddDays() *RecurrenceRuleBuilder {
	b.rule.RecurrenceModifier = append(b.rule.RecurrenceModifier, ModifierOdd)
	return b
}

func (b *RecurrenceRuleBuilder) Build() (*RecurrenceRule, error) {
	if b.hasModifier(ModifierOdd) && b.hasModifier(ModifierEven) {
		return nil, fmt.Errorf("нельзя одновременно использовать четные и нечетные дни")
	}

	if b.hasModifier(ModifierWeekdays) && b.hasModifier(ModifierWeekends) {
		return nil, fmt.Errorf("нельзя одновременно использовать будни и выходные")
	}

	if b.rule.RecurrenceType == RecurrenceWeekly {
		for _, day := range b.rule.Days {
			if day < 0 || day > 6 {
				return nil, fmt.Errorf("день недели должен быть от 0 (воскресенье) до 6 (суббота)")
			}
		}
	}

	if b.rule.RecurrenceType == RecurrenceMonthly {
		for _, day := range b.rule.Days {
			if day > 31 || day < 1 {
				return nil, fmt.Errorf("день должен быть от 1 до 31")
			}
		}
	}

	return b.rule, nil
}

func (b *RecurrenceRuleBuilder) hasModifier(mod RecurrenceModifier) bool {
	for _, m := range b.rule.RecurrenceModifier {
		if m == mod {
			return true
		}
	}
	return false
}
