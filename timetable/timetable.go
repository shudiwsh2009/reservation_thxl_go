package main

import (
	"flag"
	"fmt"
	"github.com/scorredoira/email"
	"github.com/shudiwsh2009/reservation_thxl_go/config"
	"github.com/shudiwsh2009/reservation_thxl_go/models"
	"github.com/shudiwsh2009/reservation_thxl_go/utils"
	"github.com/shudiwsh2009/reservation_thxl_go/workflow"
	"gopkg.in/mgo.v2"
	"log"
	"net/mail"
	"sort"
	"time"
	"strings"
)

func main() {
	conf := flag.String("conf", "../config/thxl.conf", "conf file path")
	isSmock := flag.Bool("smock", true, "is smock server")
	mailTo := flag.String("mail-to", "shudiwsh2009@gmail.com", "mail to list")
	flag.Parse()
	config.InitWithParams(*conf, *isSmock)
	log.Printf("config loaded: %+v", conf)
	// 数据库连接
	var session *mgo.Session
	var err error
	if config.Instance().IsSmockServer() {
		session, err = mgo.Dial("127.0.0.1:27017")
	} else {
		mongoDbDialInfo := mgo.DialInfo{
			Addrs:    []string{config.Instance().MongoHost},
			Timeout:  60 * time.Second,
			Database: config.Instance().MongoDatabase,
			Username: config.Instance().MongoUser,
			Password: config.Instance().MongoPassword,
		}
		session, err = mgo.DialWithInfo(&mongoDbDialInfo)
	}
	if err != nil {
		log.Fatalf("连接数据库失败：%v", err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	models.Mongo = session.DB("reservation_thxl")
	// 时区
	if utils.Location, err = time.LoadLocation("Asia/Shanghai"); err != nil {
		log.Printf("初始化时区失败：%v", err)
		return
	}
	// timetable
	today := utils.GetToday()
	tomorrow := today.AddDate(0, 0, 1)
	reservations, err := models.GetReservationsBetweenTime(today, tomorrow)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	todayDate := today.Format(utils.DATE_PATTERN)
	if timedReservations, err := models.GetTimedReservationsByWeekday(today.Weekday()); err == nil {
		for _, tr := range timedReservations {
			if !tr.Exceptions[todayDate] && !tr.Timed[todayDate] {
				reservations = append(reservations, tr.ToReservation(today))
			}
		}
	}
	sort.Sort(models.ReservationSlice(reservations))
	filename := "timetable_" + todayDate + utils.CsvSuffix
	if err = workflow.ExportTodayReservationTimetable(reservations, filename); err != nil {
		log.Printf("%v", err)
		return
	}
	// email
	title := fmt.Sprintf("【心理发展中心】%s咨询安排表", todayDate)
	m := email.NewMessage(title, title)
	m.From = mail.Address{Name: "", Address: config.Instance().SMTPUser}
	m.To = strings.Split(*mailTo, ",")
	m.Attach(fmt.Sprintf("%s%s", utils.ExportFolder, filename))
	if err := workflow.SendEmail(m); err != nil {
		log.Printf("发送邮件失败：%v", err)
		return
	}
	log.Printf("发送邮件成功")
}
