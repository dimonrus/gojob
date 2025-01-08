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

// transform parser buffer to time part struct
func (p *parser) toTimePart() TimePart {
	tp := TimePart{}
	buf := make([]int16, TimePartLength)
	copy(buf, p.buf)
	var j int
	// Millisecond
	for i := PositionStartMillisecond; i < PositionStartSecond; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.Millisecond = buf[PositionStartMillisecond : PositionStartMillisecond+j]
	j = 0
	// Second
	for i := PositionStartSecond; i < PositionStartMinute; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.Second = buf[PositionStartSecond : PositionStartSecond+j]
	j = 0
	// Minute
	for i := PositionStartMinute; i < PositionStartHour; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.Minute = buf[PositionStartMinute : PositionStartMinute+j]
	j = 0
	// Hour
	for i := PositionStartHour; i < PositionStartDayOfWeek; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.Hour = buf[PositionStartHour : PositionStartHour+j]
	j = 0
	// DayOfWeek
	for i := PositionStartDayOfWeek; i < PositionStartDayOfMonth; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.DayOfWeek = buf[PositionStartDayOfWeek : PositionStartDayOfWeek+j]
	j = 0
	// DayOfMonth
	for i := PositionStartDayOfMonth; i < PositionStartWeekOfMonth; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.DayOfMonth = buf[PositionStartDayOfMonth : PositionStartDayOfMonth+j]
	j = 0
	// WeekOfMonth
	for i := PositionStartWeekOfMonth; i < PositionStartWeekOfYear; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.WeekOfMonth = buf[PositionStartWeekOfMonth : PositionStartWeekOfMonth+j]
	j = 0
	// WeekOfYear
	for i := PositionStartWeekOfYear; i < PositionStartMonth; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.WeekOfYear = buf[PositionStartWeekOfYear : PositionStartWeekOfYear+j]
	j = 0
	// Month
	for i := PositionStartMonth; i < TimePartLength; i++ {
		if buf[i] > -1 {
			j++
		}
	}
	tp.Month = buf[PositionStartMonth : PositionStartMonth+j]
	return tp
}

// parse cron expression
func (p *parser) parse(expression string) error {
	p.m.Lock()
	defer p.m.Unlock()
	p.reset()
	var i, j, k int
	var n, m, d = -1, -1, -1
	var N, M, D = int64(-1), int64(-1), int64(-1)
	var isRange bool
	var err error
	for i < len(expression) {
		switch true {
		case expression[i] == '*':
			if k > 3 {
				N = 1
				M = int64(positions[k+1] - positions[k])
			} else {
				N = 0
				M = int64(positions[k+1] - positions[k] - 1)
			}
			i++
		case expression[i] == '-':
			if n == -1 {
				i++
				continue
			}
			N, err = strconv.ParseInt(expression[n:i], 10, 16)
			if err != nil {
				return err
			}
			isRange = true
			i++
			m = i
		case expression[i] == '/':
			if isRange {
				M, err = strconv.ParseInt(expression[m:i], 10, 16)
				if err != nil {
					return err
				}
				isRange = false
			} else if n > -1 {
				N, err = strconv.ParseInt(expression[n:i], 10, 16)
				if err != nil {
					return err
				}
			}
			i++
			d = i
		case '0' <= expression[i] && expression[i] <= '9':
			if n == -1 {
				n = i
			}
			i++
		case expression[i] == ',' || expression[i] == ' ':
			if N == -1 && n > -1 {
				N, err = strconv.ParseInt(expression[n:i], 10, 16)
				if err != nil {
					return err
				}
			}
			if M == -1 && m > -1 {
				M, err = strconv.ParseInt(expression[m:i], 10, 16)
				if err != nil {
					return err
				}
			}
			if d > -1 {
				D, err = strconv.ParseInt(expression[d:i], 10, 16)
				if err != nil {
					return err
				}
			}
			if j == 0 {
				j = positions[k]
			}
			if N > -1 || M > -1 {
				if D == -1 {
					D = 1
				}
				if M == -1 {
					if D == 1 {
						M = N
					} else {
						M = int64(positions[k+1] - positions[k] - 1)
					}
				}
				for N <= M {
					p.buf[j] = int16(N)
					N += D
					j++
				}
			}
			if expression[i] == ' ' {
				k++
				j = 0
			}
			i++
			isRange = false
			n, m, d = -1, -1, -1
			N, M, D = -1, -1, -1
		}
	}
	if N == -1 && n > -1 {
		N, err = strconv.ParseInt(expression[n:i], 10, 16)
		if err != nil {
			return err
		}
	}
	if M == -1 && m > -1 {
		M, err = strconv.ParseInt(expression[m:i], 10, 16)
		if err != nil {
			return err
		}
	}
	if d > -1 {
		D, err = strconv.ParseInt(expression[d:i], 10, 16)
		if err != nil {
			return err
		}
	}
	if j == 0 {
		j = positions[k]
	}
	if N > -1 || M > -1 {
		if D == -1 {
			D = 1
		}
		if M == -1 {
			if D == 1 {
				M = N
			} else {
				M = int64(positions[k+1] - positions[k] - 1)
			}
		}
		for N <= M {
			p.buf[j] = int16(N)
			N += D
			j++
		}
	}
	return nil
}
