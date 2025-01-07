package gojob

import "testing"

func TestScheduleExpression_Validate(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		exp := ScheduleExpression("*/5 - * 1-20 * 3,5 - * */3-6")
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
		exp := ScheduleExpression("/5 - * 1-20 * 3,5 - * */3-6")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'part 1 (/5) can't starts from /'")
		}
	})
	t.Run("wrong_parts_number", func(t *testing.T) {
		exp := ScheduleExpression("/5 - * 1-20 * - * */3-6")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'count of expression parts must be 9. current count is: 8'")
		}
	})
	t.Run("double_special_nok", func(t *testing.T) {
		exp := ScheduleExpression("1 -- * 1-20 * 3,5 - * */3-6")
		err := exp.Validate()
		if err == nil {
			t.Fatal("must be: 'part 2 (--) can't contains double special symbols at positions 0 and 1")
		}
	})
	t.Run("double_special_ok", func(t *testing.T) {
		exp := ScheduleExpression("1 - * 1-20/3 * 3,5 - * */3-6")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("double_special_ok", func(t *testing.T) {
		exp := ScheduleExpression("1 2-4,5-9 * 1-20 * 3,5 - * */3-6")
		err := exp.Validate()
		if err != nil {
			t.Fatal(err)
		}
	})
}
