package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// constants
// The first day of each satellite system to count week number
var (
	GPST0  time.Time = time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)
	GST0   time.Time = time.Date(1999, time.August, 22, 0, 0, 0, 0, time.UTC)
	QZSST0 time.Time = time.Date(1980, time.January, 6, 0, 0, 0, 0, time.UTC)
	BDT0   time.Time = time.Date(2006, time.January, 1, 0, 0, 0, 0, time.UTC)
)

// durations
var oneDay time.Duration = time.Duration(time.Hour * 24)
var oneWeek time.Duration = time.Duration(oneDay * 7)

// highlight colors
const (
	H1 = "  \033[7m%2d\033[0m" // reversed color (default)
	H2 = "  \033[4m%2d\033[0m" // underline
)

type gnssCal struct {
	SatSys    string
	Highlight bool
	RefDate   time.Time
	Layout    calLayout
	SysTime0  time.Time
	Today     time.Time
}

type calLayout int

const (
	Layout1Month calLayout = iota
	Layout3Month
	Layout1Year
)

// flags
var (
	flagSatsys      string
	flag3mon        bool
	flagNoHighlight bool
	flagShowHelp    bool
)

func init() {
	flag.StringVar(&flagSatsys, "satsys", "GPS", "satellite system of GNSS week to be shown")
	flag.BoolVar(&flag3mon, "3", false, "three month layout")
	flag.BoolVar(&flagNoHighlight, "n", false, "turns off lighlight of today")

	flag.Usage = func() {
		w := flag.CommandLine.Output()

		fmt.Fprintf(w, "%s\n", helpMsg)
	}
}

const helpMsg = `
gnsscal - displays a GNSS calendar

Usage:
  gnsscal [Flags] [[month] year]

Description:
  The gnsscal displays a calendar similar to 'cal' command except for displaying 
  gnss week and doy. For default, gnsscal displays only the current month.
  If month or year is given, print the specified month / year. In the case only
  the year is specified, a gnss calender for one year period is displayed.

Flags:
  -h        help for gnsscal
  -n        turns off highlight of today [default: highlight on]
  -3        three-month layout that displays previous, current and next months
  -satsys   referenced satellite system of GNSS week to be shown

  Created by Satoshi Kawamoto <satoshi.pes@gmail.com> October 16, 2021
  Inspired by 'gpscal' created by Dr. Yuki Hatanaka
`

func getCalWithOpt() (cal gnssCal, err error) {
	flag.Parse()
	args := flag.Args()

	today := time.Now().Truncate(oneDay)

	// default opt
	cal = gnssCal{
		SatSys:    "GPS",
		Highlight: true,
		RefDate:   today,
		Layout:    Layout1Month,
		SysTime0:  GPST0,
		Today:     today,
	}

	switch len(args) {
	// args [[month] year]
	case 1:
		// 1 year layout
		var year int
		year, err = strconv.Atoi(args[0])

		// check errors
		if err != nil || year < 1980 {
			return cal, fmt.Errorf("invalid year: %s", args[0])
		}

		// set opts
		cal.RefDate = time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
		cal.Layout = Layout1Year
	case 2:
		// one month layout
		var year, month int
		var err error

		// check errors
		if month, err = strconv.Atoi(args[0]); err != nil {
			return cal, fmt.Errorf("invalid month: %s, error: %v", args[0], err)
		}
		if year, err = strconv.Atoi(args[1]); err != nil {
			return cal, fmt.Errorf("invalid year: %s, error: %v", args[1], err)
		}
		if month < 0 || 12 < month {
			return cal, fmt.Errorf("invalid month: %d", month)
		}
		if year < 1980 {
			return cal, fmt.Errorf("invalid year: %d", year)
		}

		// set opts
		cal.Layout = Layout1Month
		if year == today.Year() && time.Month(month) == today.Month() {
			cal.RefDate = today
		} else {
			cal.RefDate = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		}
	}

	// flags
	switch flagSatsys {
	case "GPS":
		cal.SatSys = "GPS"
		cal.SysTime0 = GPST0
	case "QZS":
		cal.SatSys = "QZS"
		cal.SysTime0 = QZSST0
	case "BDS":
		cal.SatSys = "BDS"
		cal.SysTime0 = BDT0
	case "GAL":
		cal.SatSys = "GAL"
		cal.SysTime0 = GST0
	default:
		fmt.Printf("unknown SatSys: '%s'. use GPST instead.\n", flagSatsys)
	}

	if flag3mon {
		cal.Layout = Layout3Month
	}

	if flagNoHighlight {
		cal.Highlight = false
	}

	return cal, nil
}

func main() {
	cal, err := getCalWithOpt()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// print gnss calendar
	fmt.Printf("%s\n", cal.String())
}

func (c gnssCal) String() string {
	var msg []string
	switch c.Layout {
	case Layout1Month:
		msg = c.OneMonthLayout()
	case Layout3Month:
		msg = c.ThreeMonthLayout()
	case Layout1Year:
		msg = c.OneYearLayout()
	}

	return strings.Join(msg, "\n")
}

func (c gnssCal) OneMonthLayout() (msg []string) {
	refDate := c.RefDate
	today := c.Today
	highlight := c.Highlight
	return gnssCalMonth(refDate.Year(), refDate.Month(), today, highlight, c.SysTime0)
}

//func (c gnssCal) oneYearLayout(refDate, today time.Time, highlight bool) (msg []string) {
func (c gnssCal) OneYearLayout() (msg []string) {
	year := c.RefDate.Year()
	today := c.Today
	refDate1, hl1 := time.Date(year, 2, 1, 0, 0, 0, 0, time.UTC), false
	refDate2, hl2 := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC), false
	refDate3, hl3 := time.Date(year, 8, 1, 0, 0, 0, 0, time.UTC), false
	refDate4, hl4 := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC), false

	// highlight opt
	if c.Highlight && year == today.Year() {
		if int(today.Month()) < 4 {
			refDate1, hl1 = c.Today, true
		} else if int(c.Today.Month()) < 7 {
			refDate2, hl2 = today, true
		} else if int(today.Month()) < 10 {
			refDate3, hl3 = today, true
		} else {
			refDate4, hl4 = today, true
		}
	}

	// stack 4 rows
	msg = append(msg, threeMonthLayout(refDate1, today, hl1, c.SysTime0)...)
	msg = append(msg, "")
	msg = append(msg, threeMonthLayout(refDate2, today, hl2, c.SysTime0)...)
	msg = append(msg, "")
	msg = append(msg, threeMonthLayout(refDate3, today, hl3, c.SysTime0)...)
	msg = append(msg, "")
	msg = append(msg, threeMonthLayout(refDate4, today, hl4, c.SysTime0)...)

	return msg
}

func (c gnssCal) ThreeMonthLayout() (msg []string) {
	refDate := c.RefDate
	today := c.Today
	highlight := c.Highlight
	return threeMonthLayout(refDate, today, highlight, c.SysTime0)
}

func threeMonthLayout(refDate, today time.Time, highlight bool, initialDate time.Time) (msg []string) {
	// for three-month layout
	msgc := gnssCalMonth(refDate.Year(), refDate.Month(), today, highlight, initialDate)
	lastmonth := firstDayOfLastMonth(refDate)
	nextmonth := firstDayOfNextMonth(refDate)
	msgl := gnssCalMonth(lastmonth.Year(), lastmonth.Month(), today, highlight, initialDate)
	msgn := gnssCalMonth(nextmonth.Year(), nextmonth.Month(), today, highlight, initialDate)

	// check number of lines
	N := len(msgl)
	if len(msgc) > N {
		N = len(msgc)
	}
	if len(msgn) > N {
		N = len(msgc)
	}

	var buf string
	for i := 0; i < N; i++ {
		// leftside
		if len(msgl) > i {
			buf += fmt.Sprintf("%-34s", msgl[i])
		} else {
			buf += fmt.Sprintf("%34s", "")
		}
		buf += fmt.Sprintf("    ")

		// center
		if len(msgc) > i {
			buf += fmt.Sprintf("%-34s", msgc[i])
		} else {
			buf += fmt.Sprintf("%34s", "")
		}
		buf += fmt.Sprintf("    ")

		// right side
		if len(msgn) > i {
			buf += fmt.Sprintf("%-34s", msgn[i])
		} else {
			buf += fmt.Sprintf("%34s", "")
		}
		msg = append(msg, buf)
		buf = ""
	}

	return
}

// gnssCalMonth returns calendar msg for a month.
//
// 'year', 'month' specify the month to be shown.
// If 'highlight' is true, 'today' is highlighted.
// GNSS week is calculated based on the 'initialDate'.
func gnssCalMonth(year int, month time.Month, today time.Time, highlight bool, initialDate time.Time) (msg []string) {
	var bufday, bufdoy string

	// prepare
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDayOfNextMonth(firstDay)

	// print header
	head := fmt.Sprintf("%s %4d", month.String(), year)
	msg = append(msg, fmt.Sprintf(fmt.Sprintf("%%%ds", 20+len(head)/2), head)) // centering message
	msg = append(msg, "Week   Sun Mon Tue Wed Thu Fri Sat")

	// print dates
	for date := firstDay; date.Before(lastDay); date = date.Add(oneDay) {
		if date.Equal(firstDay) || date.Weekday() == time.Sunday {
			// calculate GNSS week
			//bufday += fmt.Sprintf("%4d  ", gnssWeek(date, GPST0))
			bufday += fmt.Sprintf("%4d  ", gnssWeek(date, initialDate))
			bufdoy += "      "
			for i := 0; i < int(date.Weekday()); i++ {
				bufday += "    "
				bufdoy += "    "
			}
		}

		if date.Equal(today) && highlight {
			bufday += fmt.Sprintf(H1, date.Day()) // reversed color
		} else {
			bufday += fmt.Sprintf("  %2d", date.Day())
		}
		bufdoy += fmt.Sprintf(" %3d", doy(date))

		if date.Weekday() == time.Saturday {
			msg = append(msg, bufday)
			msg = append(msg, bufdoy)
			bufday = ""
			bufdoy = ""
		}
	}

	if lastDay.Weekday() != time.Sunday {
		msg = append(msg, bufday)
		msg = append(msg, bufdoy)
	}

	return
}

func doy(date time.Time) int {
	newYearDay := time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	return int(date.Sub(newYearDay).Seconds()/oneDay.Seconds()) + 1
}

func gnssWeek(date time.Time, initialDate time.Time) int {
	return int(date.Sub(initialDate).Seconds() / oneWeek.Seconds())
}

func firstDayOfNextMonth(date time.Time) time.Time {
	if date.Month() == time.December {
		return time.Date(date.Year()+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	} else {
		return time.Date(date.Year(), date.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	}
}

func firstDayOfLastMonth(date time.Time) time.Time {
	if date.Month() == time.January {
		return time.Date(date.Year()-1, time.December, 1, 0, 0, 0, 0, time.UTC)
	} else {
		return time.Date(date.Year(), date.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	}
}
