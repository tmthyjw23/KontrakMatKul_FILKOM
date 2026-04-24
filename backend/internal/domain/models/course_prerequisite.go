package models

import "errors"

var ErrPrerequisiteNotMet = errors.New("prerequisite course not passed")

type CoursePrerequisite struct {
	ID                    uint64 `json:"id"`
	CourseID              uint64 `json:"course_id"`
	PrerequisiteCourseID  uint64 `json:"prerequisite_course_id"`
}

type CoursePrerequisiteRequest struct {
	CourseID             uint64 `json:"course_id" binding:"required"`
	PrerequisiteCourseID uint64 `json:"prerequisite_course_id" binding:"required"`
}