package mysql

import (
	"context"
	"sistemkontrakmatkul/backend/internal/domain/models"
)

func (r *EnrollmentRepository) GetStudentSchedule(ctx context.Context, userID uint64) ([]models.StudentScheduleResponse, error) {
	const query = `
		SELECT 
			c.code, 
			c.name, 
			s.day_of_week, 
			TIME_FORMAT(s.start_time, '%H:%i:%s') AS start_time, 
			TIME_FORMAT(s.end_time, '%H:%i:%s') AS end_time, 
			s.room, 
			c.lecturer
		FROM enrollments e
		JOIN courses c ON e.course_id = c.id
		JOIN schedules s ON s.course_id = c.id
		WHERE e.user_id = ?
		ORDER BY s.day_of_week, s.start_time
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.StudentScheduleResponse
	for rows.Next() {
		var res models.StudentScheduleResponse
		if err := rows.Scan(&res.CourseCode, &res.CourseName, &res.Day, &res.StartTime, &res.EndTime, &res.Room, &res.Lecturer); err != nil {
			return nil, err
		}
		schedules = append(schedules, res)
	}
	return schedules, nil
}