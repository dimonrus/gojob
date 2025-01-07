package gojob

import (
	"testing"
)

func TestCondition(t *testing.T) {
	t.Run("simple_true", func(t *testing.T) {
		cond := NewCondition(OperatorAND, func() bool {
			return true
		})
		if !cond.IsTrue() {
			t.Fatal("must be true")
		}
	})
	t.Run("simple_false", func(t *testing.T) {
		cond := NewCondition(OperatorAND, func() bool {
			return false
		})
		if cond.IsTrue() {
			t.Fatal("must be false")
		}
	})
	t.Run("simple_double_true", func(t *testing.T) {
		cond := NewCondition(OperatorAND,
			func() bool {
				return true
			},
		)
		cond = cond.AddExpression(
			func() bool {
				return true
			})
		if !cond.IsTrue() {
			t.Fatal("must be true")
		}
	})
	t.Run("simple_double_false", func(t *testing.T) {
		cond := NewCondition(OperatorAND)
		cond = cond.SetExpression(
			func() bool {
				return false
			},
			func() bool {
				return false
			})
		if cond.IsTrue() {
			t.Fatal("must be false")
		}
	})
	t.Run("simple_or_true", func(t *testing.T) {
		cond := NewCondition(OperatorOR, func() bool {
			return true
		})
		if !cond.IsTrue() {
			t.Fatal("must be true")
		}
	})
	t.Run("simple_or_false", func(t *testing.T) {
		cond := NewCondition(OperatorOR, func() bool {
			return false
		})
		if cond.IsTrue() {
			t.Fatal("must be false")
		}
	})
	t.Run("simple_or_double_true", func(t *testing.T) {
		cond := NewCondition(OperatorOR,
			func() bool {
				return true
			},
			func() bool {
				return true
			},
		)
		if !cond.IsTrue() {
			t.Fatal("must be true")
		}
	})
	t.Run("simple_or_double_true_2", func(t *testing.T) {
		cond := NewCondition(OperatorOR,
			func() bool {
				return false
			},
			func() bool {
				return true
			},
		)
		if !cond.IsTrue() {
			t.Fatal("must be true")
		}
	})
	t.Run("simple_or_double_false_2", func(t *testing.T) {
		cond := NewCondition(OperatorAND,
			func() bool {
				return false
			},
			func() bool {
				return false
			},
		)
		cond.SetOperator(OperatorOR)
		if cond.IsTrue() {
			t.Fatal("must be false")
		}
	})
}

func TestConditionNested(t *testing.T) {
	t.Run("nested_true", func(t *testing.T) {
		//( true && true ) || (false || true)
		cond := NewCondition(OperatorAND,
			func() bool {
				return true
			},
			func() bool {
				return true
			},
		)
		nested := NewCondition(OperatorOR,
			func() bool {
				return false
			},
			func() bool {
				return true
			})
		cond = cond.Merge(OperatorAND, nested)
		if !cond.IsTrue() {
			t.Fatal("must be true")
		}
	})
	t.Run("nested_complex_true", func(t *testing.T) {
		//( false ) || ((false && false) || (true && (false || true)))
		cond := NewCondition(OperatorOR,
			func() bool {
				return false
			},
		)
		nested1 := NewCondition(OperatorAND,
			func() bool {
				return false
			},
			func() bool {
				return false
			})
		nested2 := NewCondition(OperatorAND,
			func() bool {
				return true
			})
		nested21 := NewCondition(OperatorOR,
			func() bool {
				return false
			},
			func() bool {
				return true
			})
		cond = cond.Merge(OperatorOR, nested1, nested2.Merge(OperatorAND, nested21))
		if !cond.IsTrue() {
			t.Fatal("must be true")
		}
	})
	t.Run("nested_complex_false", func(t *testing.T) {
		//( false ) || ((false && false) || (true && (false && true)))
		cond := NewCondition(OperatorOR,
			func() bool {
				return false
			},
		)
		nested1 := NewCondition(OperatorAND,
			func() bool {
				return false
			},
			func() bool {
				return false
			})
		nested2 := NewCondition(OperatorAND,
			func() bool {
				return true
			})
		nested21 := NewCondition(OperatorAND,
			func() bool {
				return false
			},
			func() bool {
				return true
			})
		cond = cond.Merge(OperatorOR, nested1, nested2.Merge(OperatorAND, nested21))
		if cond.IsTrue() {
			t.Fatal("must be false")
		}
	})
}

// goos: darwin
// goarch: arm64
// pkg: github.com/dimonrus/gojob
// cpu: Apple M2 Max
// BenchmarkCondition
// BenchmarkCondition-12    	19057852	        57.92 ns/op	       0 B/op	       0 allocs/op
func BenchmarkCondition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cond := NewCondition(OperatorOR,
			func() bool {
				return false
			},
		)
		nested1 := NewCondition(OperatorAND,
			func() bool {
				return false
			},
			func() bool {
				return false
			})
		nested2 := NewCondition(OperatorAND,
			func() bool {
				return true
			})
		nested21 := NewCondition(OperatorAND,
			func() bool {
				return false
			},
			func() bool {
				return true
			})
		cond = cond.Merge(OperatorOR, nested1, nested2.Merge(OperatorAND, nested21))
		if cond.IsTrue() {
			b.Fatal("must be false")
		}
	}
	b.ReportAllocs()
}
