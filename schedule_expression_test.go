package gojob

import "testing"

func TestScheduleExpression_Validate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		exp := ScheduleExpression("*/5 - * 1-20 * 3,5 - * */6")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("ok_1", func(t *testing.T) {
		exp := ScheduleExpression("* * * * * * * * *")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("nok", func(t *testing.T) {
		exp := ScheduleExpression("*/* * * * * * * * *")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'part 1 (*/*) can't contains double special symbols at positions 1 and 2")
		}
	})
	t.Run("bad_start", func(t *testing.T) {
		exp := ScheduleExpression("/5 - * 1-20 * 3,5 - * */3,6")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'part 1 (/5) can't starts from /'")
		}
	})
	t.Run("wrong_parts_number", func(t *testing.T) {
		exp := ScheduleExpression("/5 - * 1-20 * - * */3,6")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'count of expression parts must be 9. current count is: 8'")
		}
	})
	t.Run("zero_length_parts", func(t *testing.T) {
		exp := ScheduleExpression("        ")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'part 1 can't have a 0 length string'")
		}
	})
	t.Run("double_special_nok", func(t *testing.T) {
		exp := ScheduleExpression("1 -- * 1-20 * 3,5 - * */3,6")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'part 2 (--) can't contains double special symbols at positions 0 and 1")
		}
	})
	t.Run("double_nok_every", func(t *testing.T) {
		exp := ScheduleExpression("*/3-6 * * * * * * * *")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: part 1 (*/3-6) can't follow special symbol '-' after '/'")
		}
	})
	t.Run("double_nok_every", func(t *testing.T) {
		exp := ScheduleExpression("*/3,4-6 * * * * * * * *")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("double_special_ok", func(t *testing.T) {
		exp := ScheduleExpression("1 - * 1-20/3 * 3,5 - * */3,4-6")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("double_special_ok", func(t *testing.T) {
		exp := ScheduleExpression("1 2-4,5-9 * 1-20 * 3,5 - * */3,6")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("multiple_correct_sequence", func(t *testing.T) {
		exp := ScheduleExpression("*/5,*/10,15,16,20/4 * * * * * * * *")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("multiple_incorrect_sequence", func(t *testing.T) {
		exp := ScheduleExpression("*/5,*/10-20,15,16,20/4 * * * * * * * *")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: in part 1 (*/5,*/10-20,15,16,20/4) can't follow special symbol '-' after '/'")
		}
	})
}

func TestScheduleExpression_Parse(t *testing.T) {
	t.Run("complex", func(t *testing.T) {
		exp := ScheduleExpression("1 2-4,5-9 * 1-20 * 3,5 - * */3,6")
		tp, err := exp.Parse()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 1 {
			t.Fatal("len of Millisecond must be 1")
		}
		if len(tp.Second) != 8 {
			t.Fatal("len of Second must be 8")
		}
		if len(tp.Minute) != 60 {
			t.Fatal("len of Minute must be 60")
		}
		if len(tp.Hour) != 20 {
			t.Fatal("len of Hour must be 20")
		}
		if len(tp.DayOfWeek) != 7 {
			t.Fatal("len of DayOfWeek must be 7")
		}
		if len(tp.DayOfMonth) != 2 {
			t.Fatal("len of DayOfMonth must be 2")
		}
		if len(tp.WeekOfMonth) != 0 {
			t.Fatal("len of WeekOfMonth must be 0")
		}
		if len(tp.WeekOfYear) != 53 {
			t.Fatal("len of WeekOfYear must be 53")
		}
		if len(tp.Month) != 5 {
			t.Fatal("len of Month must be 5")
		}
	})
	t.Run("validate_error", func(t *testing.T) {
		exp := ScheduleExpression("1 //2-4,5-9 * 1-20 * 3,5 - * */3,6")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be validate error")
		}
	})
	t.Run("parse_error", func(t *testing.T) {
		exp := ScheduleExpression("*/1a4 * * * * * * * *")
		_, err := exp.Parse()
		if err == nil {
			t.Fatal("must be parse error")
		}
	})
}
