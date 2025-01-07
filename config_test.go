package gojob

import (
	"testing"
)

func TestTimePart_Validate(t *testing.T) {
	td := TimePart{
		Millisecond: []uint16{9999},
		Second:      []uint16{1, 200},
		Minute:      []uint16{1, 200},
		Hour:        []uint16{1, 200},
		DayOfWeek:   []uint16{1, 200},
		DayOfMonth:  []uint16{1, 200},
		WeekOfMonth: []uint16{1, 200},
		WeekOfYear:  []uint16{1, 200},
		Month:       []uint16{1, 200},
	}
	e := td.Validate()
	if e != nil {
		t.Fatal(e)
	}
}

func TestParseScheduleExpression(t *testing.T) {
	expr := "- /7 * 14-00 2,3,4 - 1,2 - *"
	parts := ParseScheduleExpression(expr)
	t.Log(parts)
}
