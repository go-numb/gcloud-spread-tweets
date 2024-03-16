package models

import "time"

// Schedule is a struct for scheduling posts.
type Schedule struct {
	PostID string `csv:"-" dataframe:"post_id" firestore:"post_id,omitempty" json:"post_id,omitempty"`

	IsSchedule bool `csv:"-" dataframe:"is_schedule" firestore:"is_schedule,omitempty" json:"is_schedule,omitempty"`

	// DayOfMonth is the day of the month.
	// set day
	DayOfMonth []int `csv:"-" dataframe:"day_of_month" firestore:"day_of_month,omitempty" json:"day_of_month,omitempty"`
	// DayOfWeek is the day of the week.
	// set weekday
	DayOfWeek []time.Weekday `csv:"-" dataframe:"day_of_week" firestore:"day_of_week,omitempty" json:"day_of_week,omitempty"`

	Hours   []int `csv:"-" dataframe:"hours" firestore:"hours,omitempty" json:"hours,omitempty"`
	Minutes []int `csv:"-" dataframe:"minutes" firestore:"minutes,omitempty" json:"minutes,omitempty"`
}

func (s Schedule) GetIsSchedule(t time.Time) bool {
	if !s.IsSchedule {
		return false
	}
	if len(s.DayOfMonth) == 0 && len(s.DayOfWeek) == 0 {
		return false
	}

	// check hour
	// 引数のDayが、ScheduleのDayに含まれているか
	if len(s.DayOfMonth) > 0 {
		for _, d := range s.DayOfMonth {
			if d == t.Day() {
				return true
			}
		}
	}

	if len(s.DayOfWeek) > 0 {
		for _, d := range s.DayOfWeek {
			if d == t.Weekday() {
				return true
			}
		}
	}

	return false
}
