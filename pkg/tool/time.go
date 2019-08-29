package tool

import (
	"strconv"
	"time"

	log "gopkg.in/clog.v1"
)

// Seconds-based time units
const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func HumanTimeSince(then time.Time) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()
	log.Trace("HumanTimeSince[now:%v, then:%v, diff:%d]", now, then, diff)

	switch {
	case diff <= 0:
		return "Now"
	case diff <= 2:
		return "1 Second"
	case diff < 1*Minute:
		return strconv.FormatInt(diff, 10) + " Seconds"

	case diff < 2*Minute:
		return "1 Minute"
	case diff < 1*Hour:
		return strconv.FormatInt(diff/Minute, 10) + " Minutes"

	case diff < 2*Hour:
		return "1 Hour"
	case diff < 1*Day:
		return strconv.FormatInt(diff/Hour, 10) + " Hours"

	case diff < 2*Day:
		return "1 Day"
	case diff < 1*Week:
		return strconv.FormatInt(diff/Day, 10) + " Days"

	case diff < 2*Week:
		return "1 Week"
	case diff < 1*Month:
		return strconv.FormatInt(diff/Week, 10) + " Weeks"

	case diff < 2*Month:
		return "1 Month"
	case diff < 1*Year:
		return strconv.FormatInt(diff/Month, 10) + " Months"

	case diff < 2*Year:
		return "1 Year"
	default:
		return strconv.FormatInt(diff/Year, 10) + " Years"
	}
}
