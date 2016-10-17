package model

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Teacher struct {
	Id         bson.ObjectId `bson:"_id"`
	CreateTime time.Time     `bson:"create_time"`
	UpdateTime time.Time     `bson:"update_time"`
	Username   string        `bson:"username"` // Indexed
	Password   string        `bson:"password"`
	Fullname   string        `bson:"fullname"`
	Mobile     string        `bson:"mobile"`
	UserType   UserType      `bson:"user_type"`
}

func (m *Model) AddTeacher(username string, password string, fullname string, mobile string) (*Teacher, error) {
	if len(username) == 0 || len(password) == 0 || len(fullname) == 0 || len(mobile) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	newTeacher := &Teacher{
		Id:         bson.NewObjectId(),
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Username:   username,
		Password:   password,
		Fullname:   fullname,
		Mobile:     mobile,
		UserType:   TEACHER,
	}
	if err := collection.Insert(newTeacher); err != nil {
		return nil, err
	}
	return newTeacher, nil
}

func (m *Model) UpsertTeacher(teacher *Teacher) error {
	if teacher == nil || !teacher.Id.Valid() {
		return errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	teacher.UpdateTime = time.Now()
	_, err := collection.UpsertId(teacher.Id, teacher)
	return err
}

func (m *Model) GetTeacherById(teacherId string) (*Teacher, error) {
	if len(teacherId) == 0 || !bson.IsObjectIdHex(teacherId) {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.FindId(bson.ObjectIdHex(teacherId)).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (m *Model) GetTeacherByUsername(username string) (*Teacher, error) {
	if len(username) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.Find(bson.M{"username": username}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (m *Model) GetTeacherByFullname(fullname string) (*Teacher, error) {
	if len(fullname) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.Find(bson.M{"fullname": fullname}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (m *Model) GetTeacherByMobile(mobile string) (*Teacher, error) {
	if len(mobile) == 0 {
		return nil, errors.New("字段不合法")
	}
	collection := m.mongo.C("teacher")
	teacher := &Teacher{}
	if err := collection.Find(bson.M{"mobile": mobile}).One(teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}