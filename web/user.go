package web

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/model"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"github.com/mijia/sweb/render"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
	"time"
)

type UserController struct {
	BaseMuxController
}

const (
	kUserApiBaseUrl = "/api/user"
)

func (uc *UserController) MuxHandlers(m JsonMuxer) {
	m.Get("/m", "EntryPage", uc.getEntryPage)

	baseUrl := kUserApiBaseUrl
	m.PostJson(baseUrl+"/student/login", "StudentLogin", uc.studentLogin)
	m.PostJson(baseUrl+"/student/register", "StudentRegister", uc.studentRegister)
	m.PostJson(baseUrl+"/teacher/login", "TeacherLogin", uc.teacherLogin)
	m.PostJson(baseUrl+"/admin/login", "AdminLogin", uc.adminLogin)
	m.GetJson(baseUrl+"/logout", "Logout", RoleCookieInjection(uc.logout))
}

func (uc *UserController) GetTemplates() []*render.TemplateSet {
	return []*render.TemplateSet{
		render.NewTemplateSet("entry", "desktop.html", "reservation/entry.html", "layout/desktop.html"),
	}
}

func (uc *UserController) getEntryPage(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	params := map[string]interface{}{}
	uc.RenderHtmlOr500(w, http.StatusOK, "entry", params)
	return ctx
}

func (uc *UserController) studentRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = make(map[string]interface{})

	student, err := service.Workflow().StudentRegister(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    student.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    student.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(student.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	result["url"] = "/reservation/student"

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) studentLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = make(map[string]interface{})

	student, err := service.Workflow().StudentLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    student.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    student.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(student.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	result["url"] = "/reservation/student"

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) teacherLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = make(map[string]interface{})

	teacher, err := service.Workflow().TeacherLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    teacher.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    teacher.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(teacher.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	switch teacher.UserType {
	case model.USER_TYPE_TEACHER:
		result["url"] = "/reservation/teacher"
	default:
		result["url"] = "/reservation/entry"
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) adminLogin(ctx context.Context, w http.ResponseWriter, r *http.Request) (int, interface{}) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	var result = make(map[string]interface{})

	admin, err := service.Workflow().AdminLogin(username, password)
	if err != nil {
		return http.StatusOK, wrapJsonError(err.Error())
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    admin.Id.Hex(),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    admin.Username,
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "user_type",
		Value:    strconv.Itoa(admin.UserType),
		Path:     "/",
		Expires:  time.Now().Local().AddDate(1, 0, 0),
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
	})
	switch admin.UserType {
	case model.USER_TYPE_ADMIN:
		result["url"] = "/reservation/admin"
	default:
		result["url"] = "/reservation/entry"
	}

	return http.StatusOK, wrapJsonOk(result)
}

func (uc *UserController) logout(w http.ResponseWriter, r *http.Request, userId string, userType int) (int, interface{}) {
	var result = make(map[string]interface{})

	switch userType {
	case model.USER_TYPE_ADMIN:
		result["url"] = "/reservation/admin"
	case model.USER_TYPE_TEACHER:
		result["url"] = "/reservation/teacher"
	case model.USER_TYPE_STUDENT:
		result["url"] = "/reservation/student"
	default:
		result["url"] = "/reservation/entry"
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "user_id",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "username",
		Path:   "/",
		MaxAge: -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:   "user_type",
		Path:   "/",
		MaxAge: -1,
	})

	return http.StatusOK, wrapJsonOk(result)
}
