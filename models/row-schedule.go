package models

import "time"

type TypeSchedule int

const (
	None    TypeSchedule = iota
	Yearly               // 1年ごと
	Monthly              // 1ヶ月ごと
	Weekly               // 1週間ごと
	Daily                // 1日ごと
)

// Schedule is a struct for scheduling posts.
type Schedule struct {
	PostID string `csv:"-" dataframe:"post_id" firestore:"post_id,omitempty" json:"post_id,omitempty"`

	IsSchedule bool `csv:"-" dataframe:"is_schedule" firestore:"is_schedule,omitempty" json:"is_schedule,omitempty"`

	TypeSchedule TypeSchedule `csv:"-" dataframe:"type_schedule" firestore:"type_schedule,omitempty" json:"type_schedule,omitempty"`

	// Month is the day of the month.
	// set day
	Month time.Month `csv:"-" dataframe:"month" firestore:"month,omitempty" json:"month,omitempty"`

	// Day is the day.
	Day int `csv:"-" dataframe:"day" firestore:"day,omitempty" json:"day,omitempty"`

	// Week is the day of the week.
	// set weekday
	Week time.Weekday `csv:"-" dataframe:"week" firestore:"week,omitempty" json:"week,omitempty"`

	Hour   int `csv:"-" dataframe:"hour" firestore:"hour,omitempty" json:"hour,omitempty"`
	Minute int `csv:"-" dataframe:"minute" firestore:"minute,omitempty" json:"minute,omitempty"`

	// Times is setting multiple times.
	Times []time.Time `csv:"-" dataframe:"times" firestore:"times,omitempty" json:"times,omitempty"`
}

func (s Schedule) GetIsSchedule(t time.Time) bool {
	if !s.IsSchedule {
		return false
	}
	if s.Month == 0 {
		return false
	}

	switch s.TypeSchedule {
	case Yearly:
		// check Yearly
		if s.Month == t.Month() && s.Day == t.Day() {
			if s.Hour == t.Hour() && s.Minute == t.Minute() {
				return true
			}
		}
	case Monthly:
		// check Monthly
		if s.Day == t.Day() {
			if s.Hour == t.Hour() && s.Minute == t.Minute() {
				return true
			}
		}
	case Weekly:
		// check Weekly
		if s.Week == t.Weekday() {
			if s.Hour == t.Hour() && s.Minute == t.Minute() {
				return true
			}
		}

	default:
		// check Daily
		// check Daily
		if s.Hour == t.Hour() && s.Minute == t.Minute() {
			return true
		}
	}

	return false
}
