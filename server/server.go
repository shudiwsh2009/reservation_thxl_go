package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shudiwsh2009/reservation_thxl_go/controllers"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/shudiwsh2009/reservation_thxl_go/workflow"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
	"flag"
)

var needUserPath = regexp.MustCompile("^(/reservation/(student|teacher|admin)$|/(user/logout|(student|teacher|admin)/))")
var redirectStudentPath = regexp.MustCompile("^(/reservation/student$|/student/)")
var redirectTeacherPath = regexp.MustCompile("^(/reservation/teacher$|/teacher/)")
var redirectAdminPath = regexp.MustCompile("^(/reservation/admin|/admin/)")

func handleWithCookie(fn func(http.ResponseWriter, *http.Request, string, models.UserType) interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !needUserPath.MatchString(r.URL.Path) {
			if result := fn(w, r, "", 0); result != nil {
				if data, err := json.Marshal(result); err == nil {
					w.Header().Set("Content-Type", "application/json;charset=UTF-8")
					w.Write(data)
				} else {
					fmt.Errorf("%v", err)
				}
			}
			return
		}
		redirectUrl := "/reservation/entry"
		if redirectStudentPath.MatchString(r.URL.Path) {
			redirectUrl = "/reservation/student/login"
		} else if redirectTeacherPath.MatchString(r.URL.Path) {
			redirectUrl = "/reservation/teacher/login"
		} else if redirectAdminPath.MatchString(r.URL.Path) {
			redirectUrl = "/reservation/admin/login"
		}
		var userId string
		var userType models.UserType
		if cookie, err := r.Cookie("user_id"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		} else {
			userId = cookie.Value
		}
		if _, err := r.Cookie("username"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		}
		if cookie, err := r.Cookie("user_type"); err != nil {
			http.Redirect(w, r, redirectUrl, http.StatusFound)
			return
		} else {
			if ut, err := strconv.Atoi(cookie.Value); err != nil {
				http.Redirect(w, r, redirectUrl, http.StatusFound)
				return
			} else {
				userType = models.UserType(ut)
			}
		}
		if result := fn(w, r, userId, userType); result != nil {
			if data, err := json.Marshal(result); err == nil {
				w.Header().Set("Content-Type", "application/json;charset=UTF-8")
				w.Write(data)
			} else {
				fmt.Errorf("%v", err)
			}
		}
	}
}

func main() {
	appEnv := flag.String("app-env", "STAGING", "app environment")
	smsUid := flag.String("sms-uid", "", "sms uid")
	smsKey := flag.String("sms-key", "", "sms key")
	mailSmtp := flag.String("mail-smtp", "", "mail smtp")
	mailUsername := flag.String("mail-username", "", "mail username")
	mailPassword := flag.String("mail-password", "", "mail password")
	flag.Parse()
	utils.APP_ENV = *appEnv
	utils.SMS_UID = *smsUid
	utils.SMS_KEY = *smsKey
	utils.MAIL_SMTP = *mailSmtp
	utils.MAIL_USERNAME = *mailUsername
	utils.MAIL_PASSWORD = *mailPassword
	fmt.Println(utils.APP_ENV, utils.SMS_UID, utils.SMS_KEY, utils.MAIL_SMTP, utils.MAIL_USERNAME, utils.MAIL_PASSWORD)
	// 数据库连接
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		fmt.Errorf("连接数据库失败：%v", err)
		return
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("reservation_thxl")
	// 时区
	if utils.Location, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		fmt.Errorf("初始化时区失败：%v", err)
		return
	}
	// 初始化档案
	if err := workflow.ImportArchiveFromCSVFile(); err != nil {
		fmt.Errorf("初始化档案失败：%v", err)
	}

	// TODO: Remove the following test codes

	// mux
	router := mux.NewRouter()
	// 加载页面处理器
	pageRouter := router.PathPrefix("/reservation").Methods("GET").Subrouter()
	pageRouter.HandleFunc("/", handleWithCookie(controllers.EntryPage))
	pageRouter.HandleFunc("/entry", handleWithCookie(controllers.EntryPage))
	pageRouter.HandleFunc("/student", handleWithCookie(controllers.StudentPage))
	pageRouter.HandleFunc("/student/login", handleWithCookie(controllers.StudentLoginPage))
	pageRouter.HandleFunc("/student/register", handleWithCookie(controllers.StudentRegisterPage))
	pageRouter.HandleFunc("/teacher", handleWithCookie(controllers.TeacherPage))
	pageRouter.HandleFunc("/teacher/login", handleWithCookie(controllers.TeacherLoginPage))
	pageRouter.HandleFunc("/admin", handleWithCookie(controllers.AdminPage))
	pageRouter.HandleFunc("/admin/login", handleWithCookie(controllers.AdminLoginPage))
	pageRouter.HandleFunc("/admin/timetable", handleWithCookie(controllers.AdminTimetablePage))
	pageRouter.HandleFunc("/protocol", handleWithCookie(controllers.ProtocolPage))
	// 加载动态处理器
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/student/login", handleWithCookie(controllers.StudentLogin)).Methods("POST")
	userRouter.HandleFunc("/student/register", handleWithCookie(controllers.StudentRegister)).Methods("POST")
	userRouter.HandleFunc("/teacher/login", handleWithCookie(controllers.TeacherLogin)).Methods("POST")
	userRouter.HandleFunc("/admin/login", handleWithCookie(controllers.AdminLogin)).Methods("POST")
	userRouter.HandleFunc("/logout", handleWithCookie(controllers.Logout)).Methods("GET")
	studentRouter := router.PathPrefix("/student").Subrouter()
	studentRouter.HandleFunc("/reservation/view", handleWithCookie(controllers.ViewReservationsByStudent)).Methods("GET")
	studentRouter.HandleFunc("/reservation/make", handleWithCookie(controllers.MakeReservationByStudent)).Methods("POST")
	studentRouter.HandleFunc("/reservation/feedback/get", handleWithCookie(controllers.GetFeedbackByStudent)).Methods("POST")
	studentRouter.HandleFunc("/reservation/feedback/submit", handleWithCookie(controllers.SubmitFeedbackByStudent)).Methods("POST")
	teacherRouter := router.PathPrefix("/teacher").Subrouter()
	teacherRouter.HandleFunc("/reservation/view", handleWithCookie(controllers.ViewReservationsByTeacher)).Methods("GET")
	teacherRouter.HandleFunc("/reservation/feedback/get", handleWithCookie(controllers.GetFeedbackByTeacher)).Methods("POST")
	teacherRouter.HandleFunc("/reservation/feedback/submit", handleWithCookie(controllers.SubmitFeedbackByTeacher)).Methods("POST")
	teacherRouter.HandleFunc("/student/get", handleWithCookie(controllers.GetStudentInfoByTeacher)).Methods("POST")
	teacherRouter.HandleFunc("/student/query", handleWithCookie(controllers.QueryStudentInfoByTeacher)).Methods("POST")
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/timetable/view", handleWithCookie(controllers.ViewTimedReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/timetable/add", handleWithCookie(controllers.AddTimedReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/edit", handleWithCookie(controllers.EditTimedReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/remove", handleWithCookie(controllers.RemoveTimedReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/open", handleWithCookie(controllers.OpenTimedReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/timetable/close", handleWithCookie(controllers.CloseTimedReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/view", handleWithCookie(controllers.ViewReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/view/daily", handleWithCookie(controllers.ViewDailyReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/export/today", handleWithCookie(controllers.ExportTodayReservationsByAdmin)).Methods("GET")
	adminRouter.HandleFunc("/reservation/export/report", handleWithCookie(controllers.ExportReportFormByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/export/report/monthly", handleWithCookie(controllers.ExportReportMonthlyByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/add", handleWithCookie(controllers.AddReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/edit", handleWithCookie(controllers.EditReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/remove", handleWithCookie(controllers.RemoveReservationsByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/cancel", handleWithCookie(controllers.CancelReservationByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/feedback/get", handleWithCookie(controllers.GetFeedbackByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/feedback/submit", handleWithCookie(controllers.SubmitFeedbackByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/reservation/student/set", handleWithCookie(controllers.SetStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/get", handleWithCookie(controllers.GetStudentInfoByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/crisis/update", handleWithCookie(controllers.UpdateStudentCrisisLevelByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/archive/update", handleWithCookie(controllers.UpdateStudentArchiveNumberByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/password/reset", handleWithCookie(controllers.ResetStudentPasswordByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/account/delete", handleWithCookie(controllers.DeleteStudentAccountByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/export", handleWithCookie(controllers.ExportStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/unbind", handleWithCookie(controllers.UnbindStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/bind", handleWithCookie(controllers.BindStudentByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/student/query", handleWithCookie(controllers.QueryStudentInfoByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/teacher/search", handleWithCookie(controllers.SearchTeacherByAdmin)).Methods("POST")
	adminRouter.HandleFunc("/teacher/workload", handleWithCookie(controllers.GetTeacherWorkloadByAdmin)).Methods("POST")
	categoryRouter := router.PathPrefix("/category").Subrouter()
	categoryRouter.HandleFunc("/feedback", handleWithCookie(controllers.GetFeedbackCategories)).Methods("GET")
	// http加载处理器
	http.Handle("/", router)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}