package models

import "errors"

var ErrFailedToFetchCourses = errors.New("failed to fetch courses")

type CourseScheduleResponse struct {
	DayOfWeek string `json:"day_of_week"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Room      string `json:"room"`
}

type CourseResponse struct {
	ID        uint64                   `json:"id"`
	Code      string                   `json:"code"`
	Name      string                   `json:"name"`
	SKS       int                      `json:"sks"`
	Quota     int                      `json:"quota"`
	Lecturer  string                   `json:"lecturer"`
	Schedules []CourseScheduleResponse `json:"schedules"`
}
