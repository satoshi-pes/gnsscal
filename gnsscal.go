// MIT License
//
// Copyright (c) 2021 Satoshi Kawamoto <satoshi.pes@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// 'gnsscal' - Command similar to 'cal', but also print GNSS week, doy.
// inspired by gpscal created by Dr. Yuki Hatanaka.

package gnsscal

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
	SatSys    SatSys
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

type SatSys string

const (
	SYSGPS SatSys = "GPS"
	SYSGLO SatSys = "GLO"
	SYSGAL SatSys = "GAL"
	SYSQZS SatSys = "QZS"
	SYSBDS SatSys = "BDS"
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
  -satsys   satellite system of GNSS week; 'GPS', 'QZS', 'GAL', 'BDS', or 'GLO' [default: GPS]

  Created by Satoshi Kawamoto <satoshi.pes@gmail.com> October 16, 2021
  Inspired by 'gpscal' created by Dr. Yuki Hatanaka
`

func getCalWithOpt() (cal gnssCal, err error) {
	flag.Parse()
	args := flag.Args()

	today := time.Now().Truncate(oneDay)

	// default opt
	cal = gnssCal{
		SatSys:    SYSGPS,
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
		cal.SatSys = SYSGPS
		cal.SysTime0 = GPST0
	case "QZS":
		cal.SatSys = SYSQZS
		cal.SysTime0 = QZSST0
	case "BDS":
		cal.SatSys = SYSBDS
		cal.SysTime0 = BDT0
	case "GAL":
		cal.SatSys = SYSGAL
		cal.SysTime0 = GST0
	case "GLO":
		cal.SatSys = SYSGLO
		cal.SysTime0 = leapYearDate(cal.RefDate) // Glonass week starts from the first day of leap year
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
	return gnssCalMonth(refDate.Year(), refDate.Month(), c.Today, c.Highlight, c.SysTime0, c.SatSys)
}

func (c gnssCal) OneYearLayout() (msg []string) {
	year := c.RefDate.Year()
	today := c.Today
	refDate1 := time.Date(year, 2, 1, 0, 0, 0, 0, time.UTC)
	refDate2 := time.Date(year, 5, 1, 0, 0, 0, 0, time.UTC)
	refDate3 := time.Date(year, 8, 1, 0, 0, 0, 0, time.UTC)
	refDate4 := time.Date(year, 11, 1, 0, 0, 0, 0, time.UTC)

	// stack 4 rows
	msg = append(msg, threeMonthLayout(refDate1, today, c.Highlight, c.SysTime0, c.SatSys)...)
	msg = append(msg, "")
	msg = append(msg, threeMonthLayout(refDate2, today, c.Highlight, c.SysTime0, c.SatSys)...)
	msg = append(msg, "")
	msg = append(msg, threeMonthLayout(refDate3, today, c.Highlight, c.SysTime0, c.SatSys)...)
	msg = append(msg, "")
	msg = append(msg, threeMonthLayout(refDate4, today, c.Highlight, c.SysTime0, c.SatSys)...)

	return msg
}

func (c gnssCal) ThreeMonthLayout() (msg []string) {
	return threeMonthLayout(c.RefDate, c.Today, c.Highlight, c.SysTime0, c.SatSys)
}

func threeMonthLayout(refDate, today time.Time, highlight bool, initialDate time.Time, sys SatSys) (msg []string) {
	// for three-month layout
	msgc := gnssCalMonth(refDate.Year(), refDate.Month(), today, highlight, initialDate, sys)

	var msgl, msgr []string
	lastmonth := firstDayOfLastMonth(refDate)
	nextmonth := firstDayOfNextMonth(refDate)
	if sys == SYSGLO {
		msgl = gnssCalMonth(lastmonth.Year(), lastmonth.Month(), today, highlight, leapYearDate(lastmonth), sys)
		msgr = gnssCalMonth(nextmonth.Year(), nextmonth.Month(), today, highlight, leapYearDate(nextmonth), sys)
	} else {
		msgl = gnssCalMonth(lastmonth.Year(), lastmonth.Month(), today, highlight, initialDate, sys)
		msgr = gnssCalMonth(nextmonth.Year(), nextmonth.Month(), today, highlight, initialDate, sys)
	}

	// check number of lines
	N := len(msgl)
	if len(msgc) > N {
		N = len(msgc)
	}
	if len(msgr) > N {
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
		if len(msgr) > i {
			buf += fmt.Sprintf("%-34s", msgr[i])
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
//
// Note that the initialDate may not start from Sunday for GLONASS.
// So the week numbers are calculated at first day of the month and
// Sundays, and the same week numbers could be printed.
func gnssCalMonth(year int, month time.Month, today time.Time, highlight bool, initialDate time.Time, sys SatSys) (msg []string) {
	var bufday, bufdoy string

	// prepare
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDayOfNextMonth(firstDay)

	// print header
	head := fmt.Sprintf("%s %4d", month.String(), year)
	msg = append(msg, fmt.Sprintf(fmt.Sprintf("%%s%%%ds", 17+len(head)/2), sys, head)) // centering message
	msg = append(msg, "Week   Sun Mon Tue Wed Thu Fri Sat")

	// print dates
	for date := firstDay; date.Before(lastDay); date = date.Add(oneDay) {
		if date.Equal(firstDay) || date.Weekday() == time.Sunday {
			// calculate GNSS week
			if date.Before(initialDate) {
				bufday += "      "
			} else {
				bufday += fmt.Sprintf("%4d  ", gnssWeek(date, initialDate))
			}
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

func gloWeek(date time.Time) int {
	return gnssWeek(date, leapYearDate(date))
}

func leapYearDate(date time.Time) time.Time {
	year := date.Year()
	leapYear := year - year%4

	return time.Date(leapYear, 1, 1, 0, 0, 0, 0, time.UTC)
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
