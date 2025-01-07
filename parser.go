package gojob

import (
	"strconv"
	"sync"
)

const (
	// TimePartLength Whole parser buffer length
	TimePartLength = 1251

	PositionStartMillisecond = 0
	PositionStartSecond      = 1000
	PositionStartMinute      = 1060
	PositionStartHour        = 1120
	PositionStartDayOfWeek   = 1144
	PositionStartDayOfMonth  = 1151
	PositionStartWeekOfMonth = 1182
	PositionStartWeekOfYear  = 1187
	PositionStartMonth       = 1239
)

// parser struct
type parser struct {
	// buffer
	buf []int16
	// parser must be thread safe
	m sync.RWMutex
}

var (
	// for parser
	parserBuf = initParser()

	// format positions
	positions = []int{
		PositionStartMillisecond, PositionStartSecond, PositionStartMinute,
		PositionStartHour, PositionStartDayOfWeek, PositionStartDayOfMonth,
		PositionStartWeekOfMonth, PositionStartWeekOfYear, PositionStartMonth,
		TimePartLength,
	}
)

// init parser on application start
func initParser() parser {
	return parser{
		buf: make([]int16, TimePartLength),
	}
}

// reset all parser buff values to -1
func (p *parser) reset() {
	for i := range p.buf {
		p.buf[i] = -1
	}
}

// extract values form start to end
func (p *parser) extract(start, length int) []uint16 {
	var result = make([]uint16, length)
	var j = 0
	for i := start; i < start+length; i++ {
		if p.buf[i] != -1 {
			result[j] = uint16(p.buf[i])
			j++
		}
	}
	return result[:j]
}

// transform parser buffer to time part struct
func (p *parser) toTimePart() TimePart {
	return TimePart{
		Millisecond: p.extract(PositionStartMillisecond, PositionStartSecond-PositionStartMillisecond),
		Second:      p.extract(PositionStartSecond, PositionStartMinute-PositionStartSecond),
		Minute:      p.extract(PositionStartMinute, PositionStartHour-PositionStartMinute),
		Hour:        p.extract(PositionStartHour, PositionStartDayOfWeek-PositionStartHour),
		DayOfWeek:   p.extract(PositionStartDayOfWeek, PositionStartDayOfMonth-PositionStartDayOfWeek),
		DayOfMonth:  p.extract(PositionStartDayOfMonth, PositionStartWeekOfMonth-PositionStartDayOfMonth),
		WeekOfMonth: p.extract(PositionStartWeekOfMonth, PositionStartWeekOfYear-PositionStartWeekOfMonth),
		WeekOfYear:  p.extract(PositionStartWeekOfYear, PositionStartMonth-PositionStartWeekOfYear),
		Month:       p.extract(PositionStartMonth, TimePartLength-PositionStartMonth),
	}
}

// parse cron expression
func (p *parser) parse(expression string) TimePart {
	p.m.Lock()
	defer p.m.Unlock()
	p.reset()
	// i - expression rune iterator
	// j - sub expression start iterator
	// k - position parts iterator, positions[k]
	var i, j, k int
	for i < len(expression) {
		if i == len(expression)-1 {
			i++
		} else if expression[i] != ' ' {
			i++
			continue
		}
		switch expression[j:i] {
		case "*":
			start := positions[k]
			end := positions[k+1]
			value := int16(0)
			// see TimePart, until day of week start value is 0
			if k > 3 {
				value = 1
			}
			for m := start; m < end; m++ {
				p.buf[m] = value
				value++
			}
			fallthrough
		case "-":
			k++
		default:
			pos := positions[k]
			expr := expression[j:i]
			if len(expr) > 1 && expr[0:2] == "*/" {

			} else {
				n, m := 0, 0
				var isRange bool
				for m < len(expr) {
					if expr[m] == ',' {
						if isRange {
							num, err := strconv.ParseInt(expr[n:m], 10, 16)
							if err == nil {
								for l := p.buf[pos-1] + 1; l <= int16(num); l++ {
									p.buf[pos] = l
									pos++
								}
							}
							isRange = false
						} else {
							num, err := strconv.ParseInt(expr[n:m], 10, 16)
							if err == nil {
								p.buf[pos] = int16(num)
								pos++
							}
						}
						n = m + 1
					} else if expr[m] == '-' {
						isRange = true
						num, err := strconv.ParseInt(expr[n:m], 10, 16)
						if err == nil {
							p.buf[pos] = int16(num)
							pos++
						}
						n = m + 1
					}
					m++
				}
				if n < m {
					if isRange {
						num, err := strconv.ParseInt(expr[n:m], 10, 16)
						if err == nil {
							for l := p.buf[pos-1] + 1; l <= int16(num); l++ {
								p.buf[pos] = l
								pos++
							}
						}
						isRange = false
					} else {
						num, err := strconv.ParseInt(expr[n:m], 10, 16)
						if err == nil {
							p.buf[pos] = int16(num)
							pos++
						}
					}
				}
			}
			k++
		}
		i++
		j = i
	}
	// TODO perform optimisation
	return p.toTimePart()
}

// var i, j int
//	var numStart, numEnd int
//	var isRange, isEach bool
//	// > - */7,8,9-11,16-50,55 * 14-00 2,3,4 - 1,2 - *
//	for i < len(expression) {
//		switch true {
//		case expression[i] == ' ':
//			if isRange {
//				isRange = false
//			}
//			if isEach {
//				isEach = false
//			}
//			j++
//		case expression[i] == '-':
//			isRange = true
//		case expression[i] == '/':
//			isEach = true
//			numStart = i
//		case expression[i] == ',':
//		case expression[i] == '*':
//			start := positions[j]
//			end := positions[j+1]
//			value := int16(0)
//			// see TimePart, until day of week start value is 0
//			if j > 3 {
//				value = 1
//			}
//			for k := start; k < end; k++ {
//				p.buf[k] = value
//				value++
//			}
//		case '0' <= expression[i] && expression[i] <= '9':
//			numEnd = i
//		}
//		i++
//	}
