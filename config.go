package gojob

import (
	"errors"
)

type Config struct {
}

// TimePart
// len is 1251
type TimePart struct {
	// Possible value is 0-999
	Millisecond []uint16 `yaml:"millisecond" json:"millisecond" valid:"range~0:999;"`
	// Possible value is 0-59
	Second []uint16 `yaml:"second" json:"second" valid:"range~0:59;"`
	// Possible value is 0-59
	Minute []uint16 `yaml:"minute" json:"minute" valid:"range~0:59;"`
	// Possible value is 0-23
	Hour []uint16 `yaml:"hour" json:"hour" valid:"range~0:23;"`
	// Possible value is 1-7
	DayOfWeek []uint16 `yaml:"dayOfWeek" json:"dayOfWeek" valid:"range~1:7;"`
	// Possible value is 1-31
	DayOfMonth []uint16 `yaml:"dayOfMont" json:"dayOfMont" valid:"range~1:31;"`
	// Possible value is 1-5
	WeekOfMonth []uint16 `yaml:"weekOfMonth" json:"weekOfMonth" valid:"range~1:5;"`
	// Possible value is 1-52
	WeekOfYear []uint16 `yaml:"weekOfYear" json:"weekOfYear" valid:"range~1:52;"`
	// Possible value is 1-12
	Month []uint16 `yaml:"month" json:"month" valid:"range~1:12;"`
}

// Validate check if values is incorrect
func (t TimePart) Validate() error {
	if len(t.Millisecond) > 0 {
		for i := range t.Millisecond {
			if t.Millisecond[i] > 999 {
				return errors.New("millisecond has a range between 0 and 999")
			}
		}
	}
	if len(t.Second) > 0 {
		for i := range t.Second {
			if t.Second[i] > 59 {
				return errors.New("second has range a between 0 and 59")
			}
		}
	}
	if len(t.Minute) > 0 {
		for i := range t.Minute {
			if t.Minute[i] > 59 {
				return errors.New("minute has range a between 0 and 59")
			}
		}
	}
	if len(t.Hour) > 0 {
		for i := range t.Hour {
			if t.Hour[i] > 23 {
				return errors.New("hour has range a between 0 and 23")
			}
		}
	}
	if len(t.DayOfWeek) > 0 {
		for i := range t.DayOfWeek {
			if t.DayOfWeek[i] > 7 {
				return errors.New("day of week has a range between 1 and 7")
			}
		}
	}
	if len(t.DayOfMonth) > 0 {
		for i := range t.DayOfMonth {
			if t.DayOfMonth[i] > 31 {
				return errors.New("day of month has a range between 1 and 31")
			}
		}
	}
	if len(t.WeekOfYear) > 0 {
		for i := range t.WeekOfYear {
			if t.WeekOfYear[i] > 52 {
				return errors.New("week of year has a range between 1 and 52")
			}
		}
	}
	if len(t.Month) > 0 {
		for i := range t.Month {
			if t.Month[i] > 12 {
				return errors.New("month has range a between 1 and 12")
			}
		}
	}
	return nil
}

// Эта история про повторять
// минимальное число - это и есть период?
// Во вторник с 14:00:45 по четверг до 00:00:59, в первые 2 недели каждого месяца, повторять каждые 7 секунд
// > - /7 * 14-00 2,3,4 - 1,2 - *
// * * * * * * * * *
// 1 2 3 4 5 6 7 8 9

// ParseScheduleExpression prepare time struct
func ParseScheduleExpression(expression string) TimePart {
	var i, j int
	var t TimePart
	for i < len(expression) {
		switch true {
		case expression[i] == '-':
		case expression[i] == ' ':

			j++
		case expression[i] == '/':
		case expression[i] == ',':
		case expression[i] == '*':
		case '0' <= expression[i] && expression[i] <= '9':
			_ = expression[i]
		}
		i++
	}
	return t
}
