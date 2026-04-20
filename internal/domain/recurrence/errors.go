package recurrence

import "errors"

var ErrDay = errors.New("Day must be from 0 (sunday) to 6 (saturday)")
var ErrDayMonth = errors.New("Day must be 1 to 31")
var ErrWeekdaysWeekends = errors.New("You cannot use weekdays and weekends at the same time")
var ErrOddEven = errors.New("You cannot use odd and even at the same time")
