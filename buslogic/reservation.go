package buslogic

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	re "bitbucket.org/shudiwsh2009/reservation_thxl_go/rerror"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/utils"
	"sort"
	"time"
)

const CHECK_FORCE_ERROR = "%CHECK%"

// 学生查看前后一周内的所有咨询
func (w *Workflow) GetReservationsByStudent(userId string, userType int) ([]*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("student not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("user is not student", nil, re.ERROR_NOT_AUTHORIZED)
	}
	student, err := w.mongoClient.GetStudentById(userId)
	if err != nil || student.UserType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("fail to get student", err, re.ERROR_DATABASE)
	}
	from := time.Now().AddDate(0, 0, -7)
	to := time.Now().AddDate(0, 0, 7).Add(-90 * time.Minute)
	reservations, err := w.mongoClient.GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.RESERVATION_STATUS_AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		} else if r.StudentId == student.Id.Hex() {
			if !r.TeacherFeedback.IsEmpty() && r.TeacherFeedback.Participants[0] == 0 {
				// 学生未参与的咨询不展示给学生（家长、老师或者辅导员参加）
				continue
			}
			result = append(result, r)
		} else if student.BindedTeacherId == "" || student.BindedTeacherId == r.TeacherId {
			result = append(result, r)
		}
		//} else if r.TeacherId == student.BindedTeacherId && r.Status == model.AVAILABLE {
		//	result = append(result, r)
		//} else if student.BindedTeacherId == "" && r.Status == model.AVAILABLE {
		//	result = append(result, r)
		//}
	}
	timedReservations, err := w.mongoClient.GetAllTimedReservations()
	if err != nil {
		return result, nil
	}
	today := utils.BeginOfDay(time.Now())
	for _, tr := range timedReservations {
		if tr.Status != model.RESERVATION_STATUS_AVAILABLE {
			continue
		}
		if student.BindedTeacherId != "" && student.BindedTeacherId != tr.TeacherId {
			continue
		}
		minusWeekday := int(tr.Weekday - today.Weekday())
		if minusWeekday < 0 {
			minusWeekday += 7
		}
		date := today.AddDate(0, 0, minusWeekday)
		if utils.ConcatTime(date, tr.StartTime).Before(time.Now()) {
			date = today.AddDate(0, 0, 7)
		}
		if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
			result = append(result, tr.ToReservation(date))
		}
	}
	sort.Sort(ByStartTimeOfReservation(result))
	return result, nil
}

// 咨询师查看负7天之后的所有咨询
func (w *Workflow) GetReservationsByTeacher(userId string, userType int) ([]*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("teacher not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_STUDENT {
		return nil, re.NewRErrorCode("user is not teacher", nil, re.ERROR_NOT_AUTHORIZED)
	}
	teacher, err := w.mongoClient.GetTeacherById(userId)
	if err != nil || teacher.UserType != model.USER_TYPE_TEACHER {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.mongoClient.GetReservationsAfterTime(from)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.RESERVATION_STATUS_AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		} else if r.TeacherId == teacher.Id.Hex() {
			result = append(result, r)
		}
	}
	if timedReservations, err := w.mongoClient.GetTimedReservationsByTeacherId(teacher.Id.Hex()); err == nil {
		today := utils.BeginOfDay(time.Now())
		for _, tr := range timedReservations {
			if tr.Status != model.RESERVATION_STATUS_AVAILABLE {
				continue
			}
			minusWeekday := int(tr.Weekday - today.Weekday())
			if minusWeekday < 0 {
				minusWeekday += 7
			}
			date := today.AddDate(0, 0, minusWeekday)
			if utils.ConcatTime(date, tr.StartTime).Before(time.Now()) {
				date = today.AddDate(0, 0, 7)
			}
			if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
				result = append(result, tr.ToReservation(date))
			}
			for i := 1; i <= 3; i++ {
				// 改变i的上阈值可以改变预设咨询的查看范围
				date = date.AddDate(0, 0, 7)
				if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
					result = append(result, tr.ToReservation(date))
				}
			}
		}
	}
	sort.Sort(ByStartTimeOfReservation(result))
	return result, nil
}

// 管理员查看负7天之后的所有咨询
func (w *Workflow) GetReservationsByAdmin(userId string, userType int) ([]*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", nil, re.ERROR_DATABASE)
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.mongoClient.GetReservationsAfterTime(from)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.RESERVATION_STATUS_AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		}
		result = append(result, r)
	}
	if timedReservations, err := w.mongoClient.GetAllTimedReservations(); err == nil {
		today := utils.BeginOfDay(time.Now())
		for _, tr := range timedReservations {
			if tr.Status != model.RESERVATION_STATUS_AVAILABLE {
				continue
			}
			minusWeekday := int(tr.Weekday - today.Weekday())
			if minusWeekday < 0 {
				minusWeekday += 7
			}
			date := today.AddDate(0, 0, minusWeekday)
			if utils.ConcatTime(date, tr.StartTime).Before(time.Now()) {
				date = today.AddDate(0, 0, 7)
			}
			if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
				result = append(result, tr.ToReservation(date))
			}
			for i := 1; i <= 3; i++ {
				// 改变i的上阈值可以改变预设咨询的查看范围
				date = date.AddDate(0, 0, 7)
				if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
					result = append(result, tr.ToReservation(date))
				}
			}
		}
	}
	sort.Sort(ByStartTimeOfReservation(result))
	return result, nil
}

// 管理员查看指定日期的所有咨询
func (w *Workflow) GetReservationsDailyByAdmin(fromDate string, userId string, userType int) ([]*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", nil, re.ERROR_DATABASE)
	}
	from, err := time.ParseInLocation("2006-01-02", fromDate, time.Local)
	if err != nil {
		return nil, re.NewRErrorCodeContext("from date is not valid", err, re.ERROR_INVALID_PARAM, "from_date")
	}
	to := from.AddDate(0, 0, 1)
	reservations, err := w.mongoClient.GetReservationsBetweenTime(from, to)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	if timedReservations, err := w.mongoClient.GetTimedReservationsByWeekday(from.Weekday()); err == nil {
		for _, tr := range timedReservations {
			if tr.Status == model.RESERVATION_STATUS_AVAILABLE && !tr.Exceptions[fromDate] && !tr.Timed[fromDate] {
				reservations = append(reservations, tr.ToReservation(from))
			}
		}
	}
	sort.Sort(ByStartTimeOfReservation(reservations))
	return reservations, nil
}

// 管理员通过咨询师工号查询咨询
func (w *Workflow) GetReservationsWithTeacherUsernameByAdmin(teacherUsername string, userId string, userType int) ([]*model.Reservation, error) {
	if userId == "" {
		return nil, re.NewRErrorCode("admin not login", nil, re.ERROR_NO_LOGIN)
	} else if userType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("user is not admin", nil, re.ERROR_NOT_AUTHORIZED)
	} else if teacherUsername == "" {
		return nil, re.NewRErrorCodeContext("teacher username is empty", nil, re.ERROR_MISSING_PARAM, "teacher_username")
	}
	admin, err := w.mongoClient.GetAdminById(userId)
	if err != nil || admin.UserType != model.USER_TYPE_ADMIN {
		return nil, re.NewRErrorCode("fail to get admin", err, re.ERROR_DATABASE)
	}
	teacher, err := w.mongoClient.GetTeacherByUsername(teacherUsername)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get teacher", err, re.ERROR_DATABASE)
	}
	from := time.Now().AddDate(0, 0, -7)
	reservations, err := w.mongoClient.GetReservationsAfterTime(from)
	if err != nil {
		return nil, re.NewRErrorCode("fail to get reservations", err, re.ERROR_DATABASE)
	}
	var result []*model.Reservation
	for _, r := range reservations {
		if r.Status == model.RESERVATION_STATUS_AVAILABLE && r.StartTime.Before(time.Now()) {
			continue
		}
		if r.TeacherId != teacher.Id.Hex() {
			continue
		}
		result = append(result, r)
	}
	if timedReservations, err := w.mongoClient.GetTimedReservationsByTeacherId(teacher.Id.Hex()); err == nil {
		today := utils.BeginOfDay(time.Now())
		for _, tr := range timedReservations {
			if tr.Status != model.RESERVATION_STATUS_AVAILABLE {
				continue
			}
			minusWeekday := int(tr.Weekday - today.Weekday())
			if minusWeekday < 0 {
				minusWeekday += 7
			}
			date := today.AddDate(0, 0, minusWeekday)
			if utils.ConcatTime(date, tr.StartTime).Before(time.Now()) {
				date = today.AddDate(0, 0, 7)
			}
			if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
				result = append(result, tr.ToReservation(date))
			}
			for i := 1; i <= 3; i++ {
				// 改变i的上阈值可以改变预设咨询的查看范围
				date = date.AddDate(0, 0, 7)
				if !tr.Exceptions[date.Format("2006-01-02")] && !tr.Timed[date.Format("2006-01-02")] {
					result = append(result, tr.ToReservation(date))
				}
			}
		}
	}
	sort.Sort(ByStartTimeOfReservation(result))
	return result, nil
}

type ByStartTimeOfReservation []*model.Reservation

func (rs ByStartTimeOfReservation) Len() int {
	return len(rs)
}

func (rs ByStartTimeOfReservation) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs ByStartTimeOfReservation) Less(i, j int) bool {
	if rs[i].StartTime.Equal(rs[j].StartTime) {
		return rs[i].TeacherId < rs[j].TeacherId
	}
	return rs[i].StartTime.Before(rs[j].StartTime)
}

func (w *Workflow) WrapSimpleReservation(reservation *model.Reservation) map[string]interface{} {
	var result = make(map[string]interface{})
	if reservation == nil {
		return result
	}
	result["id"] = reservation.Id.Hex()
	result["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	result["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	result["status"] = reservation.Status
	if reservation.Status == model.RESERVATION_STATUS_RESERVATED && reservation.StartTime.Before(time.Now()) {
		result["status"] = model.RESERVATION_STATUS_FEEDBACK
	}
	result["source"] = reservation.Source
	result["source_id"] = reservation.SourceId
	result["teacher_id"] = reservation.TeacherId
	if teacher, err := w.mongoClient.GetTeacherById(reservation.TeacherId); err == nil {
		result["teacher_fullname"] = teacher.Fullname
	}
	return result
}

func (w *Workflow) WrapReservation(reservation *model.Reservation) map[string]interface{} {
	var result = make(map[string]interface{})
	if reservation == nil {
		return result
	}
	result["id"] = reservation.Id.Hex()
	result["start_time"] = reservation.StartTime.Format("2006-01-02 15:04")
	result["end_time"] = reservation.EndTime.Format("2006-01-02 15:04")
	result["status"] = reservation.Status
	if reservation.Status == model.RESERVATION_STATUS_RESERVATED && reservation.StartTime.Before(time.Now()) {
		result["status"] = model.RESERVATION_STATUS_FEEDBACK
	}
	result["source"] = reservation.Source
	result["source_id"] = reservation.SourceId
	result["teacher_id"] = reservation.TeacherId
	if teacher, err := w.mongoClient.GetTeacherById(reservation.TeacherId); err == nil {
		result["teacher_username"] = teacher.Username
		result["teacher_fullname"] = teacher.Fullname
		result["teacher_mobile"] = teacher.Mobile
	}
	result["student_id"] = reservation.StudentId
	if student, err := w.mongoClient.GetStudentById(reservation.StudentId); err == nil {
		result["student_username"] = student.Username
		result["student_fullname"] = student.Fullname
		result["student_crisis_level"] = student.CrisisLevel
	}
	result["student_feedback"] = reservation.StudentFeedback.ToStringJson()
	result["teacher_feedback"] = reservation.TeacherFeedback.ToStringJson()
	return result
}
