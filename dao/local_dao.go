package dao

import (
	"log"

	. "darkseid/models"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type LocalsDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "locals"
)

// Establish a connection to database
func (m *LocalsDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *LocalsDAO) FindAll() ([]Local, error) {
	var locals []Local
	err := db.C(COLLECTION).Find(bson.M{}).All(&locals)
	return locals, err
}

func (m *LocalsDAO) FindById(id string) (Local, error) {
	var local Local
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&local)
	return local, err
}

func (m *LocalsDAO) Insert(local Local) error {
	err := db.C(COLLECTION).Insert(&local)
	return err
}

func (m *LocalsDAO) Delete(local Local) error {
	err := db.C(COLLECTION).Remove(&local)
	return err
}

func (m *LocalsDAO) Update(local Local) error {
	error := db.C(COLLECTION).UpdateId(local.ID, &local)
	return error
}
