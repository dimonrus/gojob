package gojob

import (
	"strconv"
	"sync"
)

const (
	// TimePartLength Whole parser buffer length
	TimePartLength = 1252

	PositionStartMillisecond = 0
	PositionStartSecond      = 1000
	PositionStartMinute      = 1060
	PositionStartHour        = 1120
	PositionStartDayOfWeek   = 1144
	PositionStartDayOfMonth  = 1151
	PositionStartWeekOfMonth = 1182
	PositionStartWeekOfYear  = 1187
	PositionStartMonth       = 1240
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
func (p *parser) parse(expression string) error {
	p.m.Lock()
	defer p.m.Unlock()
	p.reset()
	// i - expression rune iterator
	// j - sub expression start iterator
	// k - position parts iterator, positions[k]
	var i, j, k int
	var isRange bool
	for i < len(expression) {
		//if i == len(expression)-1 {
		//	i++
		//}
		switch true {
		case expression[i] == ' ':
			i++
			j = i
			continue
		case expression[i] == '*':
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
			k++
		case expression[i] == '-':
			if i == 0 || expression[i-1] == ' ' {
				k++
			}
		default:
			var next bool
			pos := positions[k]
			for expression[i] != ' ' {
				if expression[i] == ',' {
					i++
					j = i
					continue
				} else if expression[i] == '-' {
					isRange = true
					i++
					j = i
					continue
				} else {
					for '0' <= expression[i] && expression[i] <= '9' {
						if i == len(expression)-1 {
							i++
							next = true
							break
						}
						i++
					}
				}
				if j < i {
					num, err := strconv.ParseInt(expression[j:i], 10, 16)
					if err != nil {
						return err
					}
					if isRange {
						for l := p.buf[pos-1] + 1; l <= int16(num); l++ {
							p.buf[pos] = l
							pos++
						}
						isRange = false
					} else {
						p.buf[pos] = int16(num)
						pos++
					}
				}
				if next {
					break
				}
			}
			k++
		}
		i++
		j = i
	}
	return nil
}
