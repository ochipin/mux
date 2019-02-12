package helpers

import (
	"fmt"
	"strings"
	"time"
)

// DateTime : 日付管理構造体
type DateTime struct {
	datetime time.Time
	format   string
}

// Second : 秒を返却する
func (date DateTime) Second() int { return date.datetime.Second() }

// Minute : 分を返却する
func (date DateTime) Minute() int { return date.datetime.Minute() }

// Hour : 時を返却する
func (date DateTime) Hour() int { return date.datetime.Hour() }

// Day : 日を返却する
func (date DateTime) Day() int { return date.datetime.Day() }

// Month : 月をtime.Month型で返却する
func (date DateTime) Month() time.Month { return date.datetime.Month() }

// Year : 年を返却する
func (date DateTime) Year() int { return date.datetime.Year() }

// YearDay : 年中の通算日を返却する
func (date DateTime) YearDay() int { return date.datetime.YearDay() }

// Weekday : 週を返却する
func (date DateTime) Weekday() time.Weekday { return date.datetime.Weekday() }

// Time : time.Time型へ変換する
func (date DateTime) Time() time.Time {
	return date.datetime
}

// Format : フォーマットを指定して日時を表示する
func (date DateTime) Format(format string) string {
	d := date.datetime
	var ampm, AMPM string
	if d.Hour() >= 12 {
		ampm = "pm"
		AMPM = "PM"
	} else {
		ampm = "am"
		AMPM = "AM"
	}
	zone, _ := d.Zone()

	var count = 0
	var weekday = int(time.Date(d.Year(), 1, 1, 0, 0, 0, 0, time.Local).Weekday())
	if weekday == 0 {
		count = 6
	}
	var U = fmt.Sprint((d.YearDay() + count) / 7)
	count = 0
	if weekday == 1 {
		count = 6
	}
	var W = fmt.Sprint((d.YearDay() + count) / 7)
	rep := strings.NewReplacer(
		"%A", d.Weekday().String(),
		"%a", d.Weekday().String()[:3],
		"%B", d.Month().String(),
		"%b", d.Month().String()[:3],
		"%c", fmt.Sprintf("%s %s %02d %02d:%02d:%02d %d",
			d.Weekday().String()[:3],
			d.Month().String()[:3],
			d.Day(), d.Hour(), d.Minute(), d.Second(), d.Year()),
		"%d", fmt.Sprintf("%02d", d.Day()),
		"%H", fmt.Sprintf("%02d", d.Hour()),
		"%I", fmt.Sprintf("%02d", d.Hour()%12),
		"%j", fmt.Sprintf("%03d", d.YearDay()),
		"%M", fmt.Sprintf("%02d", d.Minute()),
		"%m", fmt.Sprintf("%02d", d.Month()),
		"%P", ampm,
		"%p", AMPM,
		"%S", fmt.Sprintf("%02d", d.Second()),
		"%U", U,
		"%W", W,
		"%w", fmt.Sprintf("%d", d.Weekday()),
		"%X", fmt.Sprintf("%02d:%02d:%02d", d.Hour(), d.Minute(), d.Second()),
		"%x", fmt.Sprintf("%02d/%02d/%s", d.Month(), d.Day(), fmt.Sprint(d.Year())[2:]),
		"%Y", fmt.Sprint(d.Year()),
		"%y", fmt.Sprint(d.Year())[2:],
		"%Z", zone,
	)

	return rep.Replace(format)
}

func (date DateTime) String() string {
	return date.Format(date.format)
}
