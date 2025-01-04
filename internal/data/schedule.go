package data

import (
	"fmt"
	"time"
)

type Schedule struct {
	Minute int
	Hour   int
	Day    time.Weekday
}

func (s Schedule) String() string {
	return fmt.Sprintf("M%d H%d D%s", s.Minute, s.Hour, s.Day)
}

func (s Schedule) Match(time time.Time) bool {
	return (s.Day == -1 || s.Day == time.Weekday()) &&
		(s.Hour == -1 || s.Hour == time.Hour()) &&
		(s.Minute == -1 || s.Minute == time.Minute())
}
