package mongo

import (
	"gopkg.in/mgo.v2"
	"time"
)

type Mongo struct {
	Addrs          []string
	Username       string
	Password       string
	DbName         string
	CollectionName string
	TimeOut        int
	Session        *mgo.Session
	DB             *mgo.Database
	Collection     *mgo.Collection
}

func NewMongo(userName, password, dbName, collName string, timeout int, addrs []string)(* Mongo)  {
	return &Mongo {
		Addrs: addrs,
		Username: userName,
		Password: password,
		DbName: dbName,
		CollectionName: collName,
		TimeOut: timeout,
	}
}

func (m *Mongo) Connect() (err error) {
	var session *mgo.Session

	dialInfo := &mgo.DialInfo{
		Addrs:     m.Addrs,
		Direct:    true,
		Username:  m.Username,
		Password:  m.Password,
		Database:  m.DbName,
		Timeout:   time.Duration(m.TimeOut) * time.Second,
		PoolLimit: 10,
	}
	if session, err = mgo.DialWithInfo(dialInfo); err != nil {
		return err
	}
	session.SetSocketTimeout(time.Duration(m.TimeOut) * time.Second)
	session.SetSyncTimeout(time.Duration(m.TimeOut) * time.Second)
	session.SetMode(mgo.Monotonic, true)
	if m.Session != nil {
		m.Session.Close()
	}
	m.Session = session
	m.DB = session.DB(m.DbName)
	m.Collection = m.DB.C(m.CollectionName)
	return nil
}

func (m *Mongo) Ping() error {
	return m.Session.Ping()
}

func (m *Mongo) Close() {
	m.Session.Close()
}
