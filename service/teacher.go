package service

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"net/http"
	"strconv"
	"time"
)

func (s *Service) ViewReservationsByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	var result = map[string]interface{}{"state": "SUCCESS"}

	teacher, err := s.w.GetTeacherById(userId)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	var teacherJson = make(map[string]interface{})
	teacherJson["teacher_fullname"] = teacher.Fullname
	teacherJson["teacher_mobile"] = teacher.Mobile
	result["teacher_info"] = teacherJson

	reservations, err := s.w.GetReservationsByTeacher(userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	var array = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["reservation_id"] = res.Id.Hex()
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		resJson["source"] = res.Source
		resJson["source_id"] = res.SourceId
		resJson["student_id"] = res.StudentId
		if student, err := s.w.GetStudentById(res.StudentId); err == nil {
			resJson["student_crisis_level"] = student.CrisisLevel
		}
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		array = append(array, resJson)
	}
	result["reservations"] = array

	return result
}

func (s *Service) GetFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var feedback = make(map[string]interface{})
	student, reservation, err := s.w.GetFeedbackByTeacher(reservationId, sourceId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	feedback["category"] = reservation.TeacherFeedback.Category
	if len(reservation.TeacherFeedback.Participants) != len(model.PARTICIPANTS) {
		feedback["participants"] = make([]int, len(model.PARTICIPANTS))
	} else {
		feedback["participants"] = reservation.TeacherFeedback.Participants
	}
	feedback["emphasis"] = reservation.TeacherFeedback.Emphasis
	if len(reservation.TeacherFeedback.Severity) != len(model.SEVERITY) {
		feedback["severity"] = make([]int, len(model.SEVERITY))
	} else {
		feedback["severity"] = reservation.TeacherFeedback.Severity
	}
	if len(reservation.TeacherFeedback.MedicalDiagnosis) != len(model.MEDICAL_DIAGNOSIS) {
		feedback["medical_diagnosis"] = make([]int, len(model.MEDICAL_DIAGNOSIS))
	} else {
		feedback["medical_diagnosis"] = reservation.TeacherFeedback.MedicalDiagnosis
	}
	if len(reservation.TeacherFeedback.Crisis) != len(model.CRISIS) {
		feedback["crisis"] = make([]int, len(model.CRISIS))
	} else {
		feedback["crisis"] = reservation.TeacherFeedback.Crisis
	}
	feedback["record"] = reservation.TeacherFeedback.Record
	feedback["crisis_level"] = student.CrisisLevel
	result["feedback"] = feedback

	return result
}

func (s *Service) SubmitFeedbackByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	reservationId := r.PostFormValue("reservation_id")
	sourceId := r.PostFormValue("source_id")
	category := r.PostFormValue("category")
	r.ParseForm()
	participants := []string(r.Form["participants"])
	emphasis := r.PostFormValue("emphasis")
	severity := []string(r.Form["severity"])
	medicalDiagnosis := []string(r.Form["medical_diagnosis"])
	crisis := []string(r.Form["crisis"])
	record := r.PostFormValue("record")
	crisisLevel := r.PostFormValue("crisis_level")

	var result = map[string]interface{}{"state": "SUCCESS"}

	participantsInt := make([]int, 0)
	for _, p := range participants {
		if pi, err := strconv.Atoi(p); err == nil {
			participantsInt = append(participantsInt, pi)
		}
	}
	severityInt := make([]int, 0)
	for _, s := range severity {
		if si, err := strconv.Atoi(s); err == nil {
			severityInt = append(severityInt, si)
		}
	}
	medicalDiagnosisInt := make([]int, 0)
	for _, m := range medicalDiagnosis {
		if mi, err := strconv.Atoi(m); err == nil {
			medicalDiagnosisInt = append(medicalDiagnosisInt, mi)
		}
	}
	crisisInt := make([]int, 0)
	for _, c := range crisis {
		if ci, err := strconv.Atoi(c); err == nil {
			crisisInt = append(crisisInt, ci)
		}
	}
	_, err := s.w.SubmitFeedbackByTeacher(reservationId, sourceId, category, participantsInt, emphasis, severityInt,
		medicalDiagnosisInt, crisisInt, record, crisisLevel, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}

	return result
}

func (s *Service) GetStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentId := r.PostFormValue("student_id")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := s.w.GetStudentInfoByTeacher(studentId, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	studentJson["student_id"] = student.Id.Hex()
	studentJson["student_username"] = student.Username
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_archive_category"] = student.ArchiveCategory
	studentJson["student_archive_number"] = student.ArchiveNumber
	studentJson["student_crisis_level"] = student.CrisisLevel
	studentJson["student_gender"] = student.Gender
	studentJson["student_email"] = student.Email
	studentJson["student_school"] = student.School
	studentJson["student_grade"] = student.Grade
	studentJson["student_current_address"] = student.CurrentAddress
	studentJson["student_mobile"] = student.Mobile
	studentJson["student_birthday"] = student.Birthday
	studentJson["student_family_address"] = student.FamilyAddress
	if !student.Experience.IsEmpty() {
		studentJson["student_experience_time"] = student.Experience.Time
		studentJson["student_experience_location"] = student.Experience.Location
		studentJson["student_experience_teacher"] = student.Experience.Teacher
	}
	studentJson["student_father_age"] = student.FatherAge
	studentJson["student_father_job"] = student.FatherJob
	studentJson["student_father_edu"] = student.FatherEdu
	studentJson["student_mother_age"] = student.MotherAge
	studentJson["student_mother_job"] = student.MotherJob
	studentJson["student_mother_edu"] = student.MotherEdu
	studentJson["student_parent_marriage"] = student.ParentMarriage
	studentJson["student_significant"] = student.Significant
	studentJson["student_problem"] = student.Problem
	if len(student.BindedTeacherId) != 0 {
		teacher, err := s.w.GetTeacherById(student.BindedTeacherId)
		if err != nil {
			studentJson["student_binded_teacher_username"] = "无"
			studentJson["student_binded_teacher_fullname"] = ""
		}
		studentJson["student_binded_teacher_username"] = teacher.Username
		studentJson["student_binded_teacher_fullname"] = teacher.Fullname
	} else {
		studentJson["student_binded_teacher_username"] = "无"
		studentJson["student_binded_teacher_fullname"] = ""
	}
	result["student_info"] = studentJson

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		resJson["student_feedback"] = res.StudentFeedback.ToJson()
		resJson["teacher_feedback"] = res.TeacherFeedback.ToJson()
		reservationJson = append(reservationJson, resJson)
	}
	result["reservations"] = reservationJson

	return result
}

func (s *Service) QueryStudentInfoByTeacher(w http.ResponseWriter, r *http.Request, userId string, userType model.UserType) interface{} {
	studentUsername := r.PostFormValue("student_username")

	var result = map[string]interface{}{"state": "SUCCESS"}

	var studentJson = make(map[string]interface{})
	student, reservations, err := s.w.QueryStudentInfoByTeacher(studentUsername, userId, userType)
	if err != nil {
		s.ErrorHandler(w, r, err)
		return nil
	}
	studentJson["student_id"] = student.Id.Hex()
	studentJson["student_username"] = student.Username
	studentJson["student_fullname"] = student.Fullname
	studentJson["student_archive_category"] = student.ArchiveCategory
	studentJson["student_archive_number"] = student.ArchiveNumber
	studentJson["student_crisis_level"] = student.CrisisLevel
	studentJson["student_key_case"] = student.KeyCase
	studentJson["student_medical_diagnosis"] = student.MedicalDiagnosis
	studentJson["student_gender"] = student.Gender
	studentJson["student_email"] = student.Email
	studentJson["student_school"] = student.School
	studentJson["student_grade"] = student.Grade
	studentJson["student_current_address"] = student.CurrentAddress
	studentJson["student_mobile"] = student.Mobile
	studentJson["student_birthday"] = student.Birthday
	studentJson["student_family_address"] = student.FamilyAddress
	if !student.Experience.IsEmpty() {
		studentJson["student_experience_time"] = student.Experience.Time
		studentJson["student_experience_location"] = student.Experience.Location
		studentJson["student_experience_teacher"] = student.Experience.Teacher
	}
	studentJson["student_father_age"] = student.FatherAge
	studentJson["student_father_job"] = student.FatherJob
	studentJson["student_father_edu"] = student.FatherEdu
	studentJson["student_mother_age"] = student.MotherAge
	studentJson["student_mother_job"] = student.MotherJob
	studentJson["student_mother_edu"] = student.MotherEdu
	studentJson["student_parent_marriage"] = student.ParentMarriage
	studentJson["student_significant"] = student.Significant
	studentJson["student_problem"] = student.Problem
	if len(student.BindedTeacherId) != 0 {
		teacher, err := s.w.GetTeacherById(student.BindedTeacherId)
		if err != nil {
			studentJson["student_binded_teacher_username"] = "无"
			studentJson["student_binded_teacher_fullname"] = ""
		}
		studentJson["student_binded_teacher_username"] = teacher.Username
		studentJson["student_binded_teacher_fullname"] = teacher.Fullname
	} else {
		studentJson["student_binded_teacher_username"] = "无"
		studentJson["student_binded_teacher_fullname"] = ""
	}
	result["student_info"] = studentJson

	var reservationJson = make([]interface{}, 0)
	for _, res := range reservations {
		resJson := make(map[string]interface{})
		resJson["start_time"] = res.StartTime.Format("2006-01-02 15:04")
		resJson["end_time"] = res.EndTime.Format("2006-01-02 15:04")
		if res.Status == model.AVAILABLE {
			resJson["status"] = model.AVAILABLE.String()
		} else if res.Status == model.RESERVATED && res.StartTime.Before(time.Now()) {
			resJson["status"] = model.FEEDBACK.String()
		} else {
			resJson["status"] = model.RESERVATED.String()
		}
		resJson["student_id"] = res.StudentId
		resJson["teacher_id"] = res.TeacherId
		if teacher, err := s.w.GetTeacherById(res.TeacherId); err == nil {
			resJson["teacher_username"] = teacher.Username
			resJson["teacher_fullname"] = teacher.Fullname
			resJson["teacher_mobile"] = teacher.Mobile
		}
		resJson["student_feedback"] = res.StudentFeedback.ToJson()
		resJson["teacher_feedback"] = res.TeacherFeedback.ToJson()
		reservationJson = append(reservationJson, resJson)
	}
	result["reservations"] = reservationJson

	return result
}