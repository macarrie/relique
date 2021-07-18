package schedule

import (
	"fmt"

	"github.com/macarrie/relique/internal/types/displayable"
)

type ScheduleDisplay struct {
	Name      string `json:"name"`
	Monday    string `json:"monday"`
	Tuesday   string `json:"tuesday"`
	Wednesday string `json:"wednesday"`
	Thursday  string `json:"thursday"`
	Friday    string `json:"friday"`
	Saturday  string `json:"saturday"`
	Sunday    string `json:"sunday"`
}

func (s Schedule) Display() displayable.Struct {
	var d displayable.Struct = ScheduleDisplay{
		Name:      s.Name,
		Monday:    s.Monday.String(),
		Tuesday:   s.Tuesday.String(),
		Wednesday: s.Wednesday.String(),
		Thursday:  s.Thursday.String(),
		Friday:    s.Friday.String(),
		Saturday:  s.Saturday.String(),
		Sunday:    s.Sunday.String(),
	}

	return d
}

func (d ScheduleDisplay) Summary() string {
	return d.Name
}

func (d ScheduleDisplay) Details() string {
	return fmt.Sprintf(
		`SCHEDULE
-------- 

Name: 		%s
Monday: 	%s
Tuesday: 	%s
Wednesday: 	%s
Thursday: 	%s
Friday: 	%s
Saturday: 	%s
Sunday: 	%s`,
		emptyToDash(d.Name),
		emptyToDash(d.Monday),
		emptyToDash(d.Tuesday),
		emptyToDash(d.Wednesday),
		emptyToDash(d.Thursday),
		emptyToDash(d.Friday),
		emptyToDash(d.Saturday),
		emptyToDash(d.Sunday),
	)
}

func (d ScheduleDisplay) TableHeaders() []string {
	return []string{
		"Name",
		"Monday",
		"Tuesday",
		"Wednesday",
		"Thursday",
		"Friday",
		"Saturday",
		"Sunday",
	}
}

func (d ScheduleDisplay) TableRow() []string {
	return []string{
		d.Name,
		emptyToDash(d.Monday),
		emptyToDash(d.Tuesday),
		emptyToDash(d.Wednesday),
		emptyToDash(d.Thursday),
		emptyToDash(d.Friday),
		emptyToDash(d.Saturday),
		emptyToDash(d.Sunday),
	}
}

func emptyToDash(str string) string {
	if str == "" {
		return "---"
	}

	return str
}
