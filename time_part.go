package gojob

import (
	"errors"
	"slices"
	"time"
)

// TimePart
// len is 1252
type TimePart struct {
	// Possible value is 0-999
	Millisecond []int16 `yaml:"millisecond" json:"millisecond" valid:"range~0:999;"`
	// Possible value is 0-59
	Second []int16 `yaml:"second" json:"second" valid:"range~0:59;"`
	// Possible value is 0-59
	Minute []int16 `yaml:"minute" json:"minute" valid:"range~0:59;"`
	// Possible value is 0-23
	Hour []int16 `yaml:"hour" json:"hour" valid:"range~0:23;"`
	// Possible value is 1-7
	DayOfWeek []int16 `yaml:"dayOfWeek" json:"dayOfWeek" valid:"range~1:7;"`
	// Possible value is 1-31
	DayOfMonth []int16 `yaml:"dayOfMont" json:"dayOfMont" valid:"range~1:31;"`
	// Possible value is 1-5
	WeekOfMonth []int16 `yaml:"weekOfMonth" json:"weekOfMonth" valid:"range~1:5;"`
	// Possible value is 1-53
	WeekOfYear []int16 `yaml:"weekOfYear" json:"weekOfYear" valid:"range~1:53;"`
	// Possible value is 1-12
	Month []int16 `yaml:"month" json:"month" valid:"range~1:12;"`
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
	if len(t.WeekOfMonth) > 0 {
		for i := range t.WeekOfMonth {
			if t.WeekOfMonth[i] > 5 {
				return errors.New("week of year has a range between 1 and 5")
			}
		}
	}
	if len(t.WeekOfYear) > 0 {
		for i := range t.WeekOfYear {
			if t.WeekOfYear[i] > 53 {
				return errors.New("week of year has a range between 1 and 53")
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

// ToCondition transform into condition
func (t TimePart) ToCondition() Condition {
	cond := NewCondition(OperatorAND)
	if len(t.Millisecond) > 0 {
		cond = cond.AddExpression(func() bool {
			return slices.Contains[[]int16, int16](t.Millisecond, int16(time.Now().UnixMilli()%1000))
		})
	}
	if len(t.Second) > 0 {
		cond = cond.AddExpression(func() bool {
			return slices.Contains[[]int16, int16](t.Second, int16(time.Now().Second()))
		})
	}
	if len(t.Minute) > 0 {
		cond = cond.AddExpression(func() bool {
			return slices.Contains[[]int16, int16](t.Minute, int16(time.Now().Minute()))
		})
	}
	if len(t.Hour) > 0 {
		cond = cond.AddExpression(func() bool {
			return slices.Contains[[]int16, int16](t.Hour, int16(time.Now().Hour()))
		})
	}
	if len(t.DayOfWeek) > 0 {
		cond = cond.AddExpression(func() bool {
			return slices.Contains[[]int16, int16](t.DayOfWeek, int16(time.Now().Weekday()+1))
		})
	}
	if len(t.DayOfMonth) > 0 {
		cond = cond.AddExpression(func() bool {
			return slices.Contains[[]int16, int16](t.DayOfMonth, int16(time.Now().Day()))
		})
	}
	if len(t.WeekOfMonth) > 0 {
		cond = cond.AddExpression(func() bool {
			day := time.Now().Day()
			return slices.Contains[[]int16, int16](t.WeekOfMonth, int16(day)/7+1)
		})
	}
	if len(t.WeekOfYear) > 0 {
		cond = cond.AddExpression(func() bool {
			_, week := time.Now().ISOWeek()
			return slices.Contains[[]int16, int16](t.WeekOfYear, int16(week))
		})
	}
	if len(t.Month) > 0 {
		cond = cond.AddExpression(func() bool {
			return slices.Contains[[]int16, int16](t.Month, int16(time.Now().Month()))
		})
	}
	return cond
}
