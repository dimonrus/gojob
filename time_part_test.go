package gojob

import (
	"context"
	"log"
	"testing"
	"time"
)

func TestTimePart_Validate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e != nil {
			t.Fatal(e)
		}
	})
	t.Run("nok_1", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{9999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Millisecond invalid")
		}
	})
	t.Run("nok_2", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 60},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Second invalid")
		}
	})
	t.Run("nok_3", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 60},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Minute invalid")
		}
	})
	t.Run("nok_4", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 24},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Hour invalid")
		}
	})
	t.Run("nok_5", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 8},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be DayOfWeek invalid")
		}
	})
	t.Run("nok_6", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 32},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be DayOfMonth invalid")
		}
	})
	t.Run("nok_7", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 6},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be WeekOfMonth invalid")
		}
	})
	t.Run("nok_8", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 54},
			Month:       []int16{1, 12},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be WeekOfYear invalid")
		}
	})
	t.Run("nok_9", func(t *testing.T) {
		td := TimePart{
			Millisecond: []int16{999},
			Second:      []int16{1, 59},
			Minute:      []int16{1, 59},
			Hour:        []int16{1, 23},
			DayOfWeek:   []int16{1, 7},
			DayOfMonth:  []int16{1, 31},
			WeekOfMonth: []int16{1, 5},
			WeekOfYear:  []int16{1, 53},
			Month:       []int16{1, 13},
		}
		e := td.Validate()
		if e == nil {
			t.Fatal("must be Month invalid")
		}
	})
}

func TestTimePart_ToCondition(t *testing.T) {
	exp := ScheduleExpression("* * * * * * * * *")
	err := exp.Validate()
	if err != nil {
		t.Fatal(err)
	}
	p := initParser()
	err = p.parse(string(exp))
	if err != nil {
		t.Fatal(err)
	}
	tp := p.toTimePart()
	job := NewJob("every.millisecond.with.500ms.repeat.duration", func(ctx context.Context, args ...any) error {
		t.Log(time.Now().Unix(), time.Now().UnixMilli()%1000)
		return nil
	}, time.Millisecond*500)
	job.SetCondition(tp.ToCondition())

	g := NewGroup(time.Second, GroupModeConsistently)
	g.AddJob(job)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*16)
	ctx = context.WithValue(ctx, "logger", log.Default())
	g.Schedule(ctx)
}
