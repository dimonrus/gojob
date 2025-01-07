package gojob

import (
	"testing"
)

func TestTimePart_Validate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("nok_1", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{9999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Millisecond invalid")
		}
	})
	t.Run("nok_2", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 60},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Second invalid")
		}
	})
	t.Run("nok_3", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 60},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Minute invalid")
		}
	})
	t.Run("nok_4", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 24},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Hour invalid")
		}
	})
	t.Run("nok_5", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 8},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be DayOfWeek invalid")
		}
	})
	t.Run("nok_6", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 32},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be DayOfMonth invalid")
		}
	})
	t.Run("nok_7", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 6},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be WeekOfMonth invalid")
		}
	})
	t.Run("nok_8", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 53},
			Month:       []uint16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be WeekOfYear invalid")
		}
	})
	t.Run("nok_9", func(t *testing.T) {
		td := TimePart{
			Millisecond: []uint16{999},
			Second:      []uint16{1, 59},
			Minute:      []uint16{1, 59},
			Hour:        []uint16{1, 23},
			DayOfWeek:   []uint16{1, 7},
			DayOfMonth:  []uint16{1, 31},
			WeekOfMonth: []uint16{1, 5},
			WeekOfYear:  []uint16{1, 52},
			Month:       []uint16{1, 13},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Month invalid")
		}
	})
}
