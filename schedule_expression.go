package gojob

import (
	"errors"
	"fmt"
	"strings"
)

// ScheduleExpression contains special symbols for cron expression declaration
// * - every number according to expression part range
// - - define range according to expression part range
// / - define each number in part range dimension
// , - specify concrete number in part range dimension
type ScheduleExpression string

// Validate check if expression is incorrect
func (s ScheduleExpression) Validate() error {
	parts := strings.Split(string(s), " ")
	if len(parts) != 9 {
		return errors.New(fmt.Sprintf("count of expression parts must be 9. current count is: %v", len(parts)))
	}
	for i := range parts {
		if len(parts[i]) == 0 {
			return errors.New(fmt.Sprintf("part %v can't have a 0 length string", i+1))
		}
		if parts[i][0] != '*' && parts[i][0] != '-' && (parts[i][0] < '0' || parts[i][0] > '9') {
			return errors.New(fmt.Sprintf("part %v (%s) can't starts from %c", i+1, parts[i], parts[i][0]))
		}
		var specialPos = -1
		for j := range parts[i] {
			if parts[i][j] != '*' && parts[i][j] != '-' && parts[i][j] != ',' && parts[i][j] != '/' && (parts[i][j] < '0' || parts[i][j] > '9') {
				return errors.New(fmt.Sprintf("part %v (%s) can't contains not valid exression symbol '%c'", i+1, parts[i], parts[i][j]))
			}
			if parts[i][j] == '*' || parts[i][j] == '-' || parts[i][j] == ',' || parts[i][j] == '/' {
				if specialPos != -1 && specialPos+1 == j {
					if !((parts[i][j] == '/' && parts[i][j-1] == '*') || (parts[i][j] == '*' && parts[i][j-1] == ',')) {
						return errors.New(fmt.Sprintf("part %v (%s) can't contains double special symbols at positions %v and %v", i+1, parts[i], specialPos, j))
					}
				}
				specialPos = j
			} else {
				specialPos = -1
			}
			var nextSpecial uint8
			if parts[i][j] == '/' {
				for k := j + 1; k < len(parts[i]); k++ {
					if parts[i][k] == ',' && nextSpecial == 0 {
						nextSpecial = parts[i][k]
					}
					if parts[i][k] == '-' && nextSpecial == 0 {
						return errors.New(fmt.Sprintf("in part %v (%s) can't follow special symbol '%c' after '%c'", i+1, parts[i], parts[i][k], parts[i][j]))
					}
				}
			}
		}
	}
	return nil
}

// Parse convert expression to TimePart struct
func (s ScheduleExpression) Parse() (TimePart, error) {
	err := s.Validate()
	if err != nil {
		return TimePart{}, err
	}
	p := initParser()
	err = p.parse(string(s))
	if err != nil {
		return TimePart{}, err
	}
	return p.toTimePart(), nil
}
