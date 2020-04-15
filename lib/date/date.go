package date

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/awoodbeck/strftime"
	"strconv"
	"strings"
	"time"
)

type Date struct {
	Base time.Time
}

func Now() *Date {
	return &Date{
		Base: time.Now(),
	}
}

func (d *Date) Midnight() *Date {
	d.Base = time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), 0, 0, 0, 0, d.Base.Location())
	return d
}

func (d *Date) Calculate(expr string) (*Date, error) {
	expr = strings.ToLower(expr)
	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return d, fmt.Errorf("unable to parse date expression:%s", expr)
	}
	if strings.Contains(expr, "midnight") {
		d.Midnight()
	}

	if fields[0] == "tomorrow" {
		d.Base = d.Base.AddDate(0, 0, 1)
		return d, nil
	} else if fields[0] == "yesterday" {
		d.Base = d.Base.AddDate(0, 0, -1)
		return d, nil
	} else if fields[0] == "today" {
		return d, nil
	}
	if len(fields) < 2 {
		return d, fmt.Errorf("unable to parse date expression:%s", expr)
	}
	var i int
	if fields[0] == "next" {
		i = 1
	} else if fields[0] == "last" {
		i = 2
	} else {
		var err error
		i, err = strconv.Atoi(fields[0])
		if err != nil {
			return d, fmt.Errorf("unable to parse date expression:%s", expr)
		}
		if len(fields) > 2 {
			if fields[2] == "after" {
				if i < 0 {
					i = i * -1
				}
			}
			if fields[2] == "before" {
				if i > 0 {
					i = i * -1
				}
			}
		}
	}

	if strings.HasPrefix(fields[1], "year") {
		d.Base = d.Base.AddDate(i, 0, 0)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), 1, 1, 0, 0, 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "month") {
		d.Base = d.Base.AddDate(0, i, 0)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), d.Base.Month(), 0, 0, 0, 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "day") {
		d.Base = d.Base.AddDate(0, 0, i)
		if strings.Contains(expr, "start") {
			d.Midnight()
		}
	} else if strings.HasPrefix(fields[1], "week") {
		d.Base = d.Base.AddDate(0, 0, i*7)

		if strings.Contains(expr, "start") {
			// Roll back to Monday:
			if wd := d.Base.Weekday(); wd == time.Sunday {
				d.Base = d.Base.AddDate(0, 0, -6)
			} else {
				d.Base = d.Base.AddDate(0, 0, -int(wd)+1)
			}
			d.Midnight()
		}

	} else if strings.HasPrefix(fields[1], "hour") {
		d.Base = d.Base.Add(time.Duration(i) * time.Hour)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), d.Base.Hour(), 0, 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "minute") {
		d.Base = d.Base.Add(time.Duration(i) * time.Minute)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), d.Base.Hour(), d.Base.Minute(), 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "second") {
		d.Base = d.Base.Add(time.Duration(i) * time.Second)
	}

	return d, nil

}

func (d *Date) DiffUnix(t int64) time.Duration {
	return time.Duration(d.Base.Unix()-t) * time.Second
}

func (d *Date) DiffDate(t Date) time.Duration {
	return time.Duration(d.Base.Unix()-t.Unix()) * time.Second
}

func (d *Date) DiffExpr(expr string) (time.Duration, error) {
	t := time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), d.Base.Hour(), d.Base.Minute(), d.Base.Second(), d.Base.Nanosecond(), d.Base.Location())
	_, err := d.Calculate(expr)
	if err != nil {
		return time.Duration(0), err
	}
	return d.DiffTime(t), nil
}

func (d *Date) DiffTime(t time.Time) time.Duration {
	return time.Duration(d.Base.Unix()-t.Unix()) * time.Second
}

func (d *Date) Format(expr string) string {
	return d.Base.Format(expr)
}

func (d *Date) FormatS(expr string) string {
	return strftime.Format(&d.Base, expr)
}

func (d *Date) Unix() int64 {
	return d.Base.Unix()
}

func (d *Date) UnixNano() int64 {
	return d.Base.UnixNano()
}

func FromString(expr string) (*Date, error) {
	t, err := dateparse.ParseLocal(expr)
	if err != nil {
		return nil, err
	}
	return &Date{
		Base: t,
	}, nil
}

func FromTime(t time.Time) *Date {
	return &Date{
		Base: t,
	}
}

func FromUnix(sec int64) *Date {
	t := time.Unix(sec, 0)
	return &Date{
		Base: t,
	}
}
