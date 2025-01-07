package gojob

// Operator Type of operator
type Operator string

const (
	// OperatorAND AND operator
	OperatorAND Operator = "AND"
	// OperatorOR OR operator
	OperatorOR Operator = "OR"
)

type Expression func() bool

// Conditions List of condition
type Conditions []Condition

// Condition job execute condition
type Condition struct {
	// Condition operator
	op Operator
	// Merge Condition operator
	mor Operator
	// list of condition expression
	expressions []Expression
	// list of condition
	conditions Conditions
}

// IsEmpty check if condition is empty
func (c Condition) IsEmpty() bool {
	return c.op == "" && c.mor == "" && len(c.expressions) == 0 && len(c.conditions) == 0
}

// AddExpression add new expressions
func (c Condition) AddExpression(expression ...Expression) Condition {
	c.expressions = append(c.expressions, expression...)
	return c
}

// SetOperator set operator
func (c Condition) SetOperator(operator Operator) Condition {
	c.op = operator
	return c
}

// SetExpression merge with condition
func (c Condition) SetExpression(expression ...Expression) Condition {
	c.expressions = expression
	return c
}

// Merge merge with condition
func (c Condition) Merge(mor Operator, condition ...Condition) Condition {
	c.mor = mor
	c.conditions = condition
	return c
}

// IsTrue check is condition is true
func (c Condition) IsTrue() bool {
	var isTrue bool
	switch c.op {
	case OperatorAND:
		for i := range c.expressions {
			if !c.expressions[i]() {
				return false
			}
		}
		if c.mor == OperatorAND {
			for i := range c.conditions {
				if !c.conditions[i].IsTrue() {
					return false
				}
			}
		}
		isTrue = true
	case OperatorOR:
		for i := range c.expressions {
			if c.expressions[i]() {
				return true
			}
		}
		if c.mor == OperatorOR {
			for i := range c.conditions {
				if c.conditions[i].IsTrue() {
					return true
				}
			}
		}
	}
	return isTrue
}

// NewCondition Init condition
func NewCondition(op Operator, expression ...Expression) Condition {
	return Condition{
		op:          op,
		expressions: expression,
	}
}
