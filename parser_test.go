package gojob

import (
	"testing"
)

func Test_parser_parse(t *testing.T) {
	t.Run("all_stars", func(t *testing.T) {
		expression := "* * * * * * * * *"
		p := initParser()
		err := p.parse(expression)
		if err != nil {
			t.Fatal(err)
		}
		tp := p.toTimePart()
		err = tp.Validate()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 1000 {
			t.Fatal("len of Millisecond must be 1000")
		}
		if len(tp.Second) != 60 {
			t.Fatal("len of Second must be 60")
		}
		if len(tp.Minute) != 60 {
			t.Fatal("len of Minute must be 60")
		}
		if len(tp.Hour) != 24 {
			t.Fatal("len of Hour must be 24")
		}
		if len(tp.DayOfWeek) != 7 {
			t.Fatal("len of DayOfWeek must be 7")
		}
		if len(tp.DayOfMonth) != 31 {
			t.Fatal("len of DayOfMonth must be 31")
		}
		if len(tp.WeekOfMonth) != 5 {
			t.Fatal("len of WeekOfMonth must be 5")
		}
		if len(tp.WeekOfYear) != 53 {
			t.Fatal("len of WeekOfYear must be 53")
		}
		if len(tp.Month) != 12 {
			t.Fatal("len of Month must be 12")
		}
	})
	t.Run("all_minus", func(t *testing.T) {
		expression := "- - - - - - - - -"
		p := initParser()
		err := p.parse(expression)
		if err != nil {
			t.Fatal(err)
		}
		tp := p.toTimePart()
		err = tp.Validate()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 0 {
			t.Fatal("len of Millisecond must be 0")
		}
		if len(tp.Second) != 0 {
			t.Fatal("len of Second must be 0")
		}
		if len(tp.Minute) != 0 {
			t.Fatal("len of Minute must be 0")
		}
		if len(tp.Hour) != 0 {
			t.Fatal("len of Hour must be 0")
		}
		if len(tp.DayOfWeek) != 0 {
			t.Fatal("len of DayOfWeek must be 0")
		}
		if len(tp.DayOfMonth) != 0 {
			t.Fatal("len of DayOfMonth must be 0")
		}
		if len(tp.WeekOfMonth) != 0 {
			t.Fatal("len of WeekOfMonth must be 0")
		}
		if len(tp.WeekOfYear) != 0 {
			t.Fatal("len of WeekOfYear must be 0")
		}
		if len(tp.Month) != 0 {
			t.Fatal("len of Month must be 0")
		}
	})
	t.Run("only_numbers", func(t *testing.T) {
		expression := "11 22 33 14 5 26 2 8 9-12"
		p := initParser()
		err := p.parse(expression)
		if err != nil {
			t.Fatal(err)
		}
		tp := p.toTimePart()
		err = tp.Validate()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 1 || tp.Millisecond[0] != 11 {
			t.Fatal("len of Millisecond must be 1 and value 11")
		}
		if len(tp.Second) != 1 || tp.Second[0] != 22 {
			t.Fatal("len of Second must be 1 and value 22")
		}
		if len(tp.Minute) != 1 || tp.Minute[0] != 33 {
			t.Fatal("len of Minute must be 1 and value 33")
		}
		if len(tp.Hour) != 1 || tp.Hour[0] != 14 {
			t.Fatal("len of Hour must be 1 and value 14")
		}
		if len(tp.DayOfWeek) != 1 || tp.DayOfWeek[0] != 5 {
			t.Fatal("len of DayOfWeek must be 1 and value 5")
		}
		if len(tp.DayOfMonth) != 1 || tp.DayOfMonth[0] != 26 {
			t.Fatal("len of DayOfMonth must be 1 and value 26")
		}
		if len(tp.WeekOfMonth) != 1 || tp.WeekOfMonth[0] != 2 {
			t.Fatal("len of WeekOfMonth must be 1 and value 2")
		}
		if len(tp.WeekOfYear) != 1 || tp.WeekOfYear[0] != 8 {
			t.Fatal("len of WeekOfYear must be 1 and value 8")
		}
		if len(tp.Month) != 4 || tp.Month[0] != 9 {
			t.Fatal("len of Month must be 4 and value 9")
		}
	})
	t.Run("numbers_with_coma_stars_minus", func(t *testing.T) {
		expression := "- * 10,20,30,40,50,0 - * - 1,2,3 - *"
		p := initParser()
		err := p.parse(expression)
		if err != nil {
			t.Fatal(err)
		}
		tp := p.toTimePart()
		err = tp.Validate()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 0 {
			t.Fatal("len of Millisecond must be 0")
		}
		if len(tp.Second) != 60 {
			t.Fatal("len of Second must be 60")
		}
		if len(tp.Minute) != 6 {
			t.Fatal("len of Minute must be 6")
		}
		if len(tp.Hour) != 0 {
			t.Fatal("len of Hour must be 0")
		}
		if len(tp.DayOfWeek) != 7 {
			t.Fatal("len of DayOfWeek must be 7")
		}
		if len(tp.DayOfMonth) != 0 {
			t.Fatal("len of DayOfMonth must be 0")
		}
		if len(tp.WeekOfMonth) != 3 {
			t.Fatal("len of WeekOfMonth must be 3")
		}
		if len(tp.WeekOfYear) != 0 {
			t.Fatal("len of WeekOfYear must be 0")
		}
		if len(tp.Month) != 12 {
			t.Fatal("len of Month must be 12")
		}
	})
	t.Run("range", func(t *testing.T) {
		expression := "650-700 * 30-40,45-59 - - - - - 3,4-6,12"
		p := initParser()
		err := p.parse(expression)
		if err != nil {
			t.Fatal(err)
		}
		tp := p.toTimePart()
		err = tp.Validate()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 51 {
			t.Fatal("len of Millisecond must be 51")
		}
		if len(tp.Second) != 60 {
			t.Fatal("len of Second must be 60")
		}
		if len(tp.Minute) != 26 {
			t.Fatal("len of Minute must be 26")
		}
		if len(tp.Hour) != 0 {
			t.Fatal("len of Hour must be 0")
		}
		if len(tp.DayOfWeek) != 0 {
			t.Fatal("len of DayOfWeek must be 0")
		}
		if len(tp.DayOfMonth) != 0 {
			t.Fatal("len of DayOfMonth must be 0")
		}
		if len(tp.WeekOfMonth) != 0 {
			t.Fatal("len of WeekOfMonth must be 0")
		}
		if len(tp.WeekOfYear) != 0 {
			t.Fatal("len of WeekOfYear must be 0")
		}
		if len(tp.Month) != 5 {
			t.Fatal("len of Month must be 5")
		}
	})
	t.Run("every.second.and.ever.5.minutes", func(t *testing.T) {
		expression := "- * */5,*/6 - - - - - -"
		p := initParser()
		err := p.parse(expression)
		if err != nil {
			t.Fatal(err)
		}
		tp := p.toTimePart()
		err = tp.Validate()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 0 {
			t.Fatal("len of Millisecond must be 0")
		}
		if len(tp.Second) != 60 {
			t.Fatal("len of Second must be 60")
		}
		if len(tp.Minute) != 22 {
			t.Fatal("len of Minute must be 26")
		}
		if len(tp.Hour) != 0 {
			t.Fatal("len of Hour must be 0")
		}
		if len(tp.DayOfWeek) != 0 {
			t.Fatal("len of DayOfWeek must be 0")
		}
		if len(tp.DayOfMonth) != 0 {
			t.Fatal("len of DayOfMonth must be 0")
		}
		if len(tp.WeekOfMonth) != 0 {
			t.Fatal("len of WeekOfMonth must be 0")
		}
		if len(tp.WeekOfYear) != 0 {
			t.Fatal("len of WeekOfYear must be 0")
		}
		if len(tp.Month) != 0 {
			t.Fatal("len of Month must be 0")
		}
	})
	t.Run("complex.every", func(t *testing.T) {
		expression := "200-300/10,400/50,800 30/5 */5,6,7-10 - - - - 10/5 1-6/2"
		p := initParser()
		err := p.parse(expression)
		if err != nil {
			t.Fatal(err)
		}
		tp := p.toTimePart()
		err = tp.Validate()
		if err != nil {
			t.Fatal(err)
		}
		if len(tp.Millisecond) != 24 {
			t.Fatal("len of Millisecond must be 24")
		}
		if len(tp.Second) != 6 {
			t.Fatal("len of Second must be 6")
		}
		if len(tp.Minute) != 17 {
			t.Fatal("len of Minute must be 17")
		}
		if len(tp.Hour) != 0 {
			t.Fatal("len of Hour must be 0")
		}
		if len(tp.DayOfWeek) != 0 {
			t.Fatal("len of DayOfWeek must be 0")
		}
		if len(tp.DayOfMonth) != 0 {
			t.Fatal("len of DayOfMonth must be 0")
		}
		if len(tp.WeekOfMonth) != 0 {
			t.Fatal("len of WeekOfMonth must be 0")
		}
		if len(tp.WeekOfYear) != 9 {
			t.Fatal("len of WeekOfYear must be 9")
		}
		if len(tp.Month) != 3 {
			t.Fatal("len of Month must be 3")
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gojob
// cpu: Apple M2 Max
// BenchmarkParser
// BenchmarkParser-12    	 2432192	       483.3 ns/op	       0 B/op	       0 allocs/op
func BenchmarkParser(b *testing.B) {
	p := initParser()
	expression := "* * * * * * * * *"
	for i := 0; i < b.N; i++ {
		p.parse(expression)
	}
	b.ReportAllocs()
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gojob
// cpu: Apple M2 Max
// BenchmarkParser2
// BenchmarkParser2-12    	 3826430	       307.6 ns/op	       0 B/op	       0 allocs/op
func BenchmarkParser2(b *testing.B) {
	p := initParser()
	expression := "100-200 */5 */10,5,6-9 - * 1-10/2 * 1-5 1,2,6,10"
	for i := 0; i < b.N; i++ {
		p.parse(expression)
	}
	b.ReportAllocs()
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gojob
// cpu: Apple M2 Max
// BenchmarkToTimePart
// BenchmarkToTimePart-12    	 1833747	       638.0 ns/op	    2688 B/op	       1 allocs/op
func BenchmarkToTimePart(b *testing.B) {
	exp := ScheduleExpression("* * * * * * * * *")
	err := exp.Validate()
	if err != nil {
		b.Fatal(err)
	}
	p := initParser()
	err = p.parse(string(exp))
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < b.N; i++ {
		tp := p.toTimePart()
		_ = tp
	}
	b.ReportAllocs()
}
