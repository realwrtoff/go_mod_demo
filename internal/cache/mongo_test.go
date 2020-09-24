package cache

import (
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

func TestMongo_Connect(t *testing.T) {
	addrs := []string{"127.0.0.1:21617"}
	ago := &Mongo{}
	ago.Addrs = addrs
	//ago.Username = "tim"
	//ago.Password = "Spurs21"
	//ago.DbName = "emar"
	//ago.CollectionName = "campaign"
	ago.DbName = "haina"
	ago.CollectionName = "guahao"
	ago.TimeOut = 30
	if err := ago.Connect(); err != nil {
		t.Errorf("connect mongdb [%s] failed [%s]", ago.Addrs, err.Error())
	} else {
		t.Logf("connect cache %s success\n", addrs)
	}
	ago.Close()
}

func TestMongo_Ping(t *testing.T) {
	addrs := []string{"127.0.0.1:21617"}
	ago := &Mongo{}

	ago.Addrs = addrs
	ago.DbName = "haina"
	ago.CollectionName = "guahao"
	ago.TimeOut = 30
	if err := ago.Connect(); err == nil {
		if err = ago.Ping(); err != nil {
			t.Errorf("ping cache %s failed[%s]", addrs, err.Error())
		} else {
			t.Logf("ping cache %s success\n", addrs)
		}
	} else {
		t.Errorf("connect cache %s failed[%s]", addrs, err.Error())
	}
	ago.Close()
}

func TestMongo_Insert(t *testing.T)  {
	addrs := []string{"127.0.0.1:21617"}
	ago := NewMongo("", "", "haina", "ceshi", 30, addrs)
	if err := ago.Connect(); err == nil {
		clickId := bson.NewObjectId()
		type Doc struct {
			ClickId bson.ObjectId `bson:"_id"`
			Name string `bson:"name"`
			CreateTime int64 `bson:"create_time"`
		}
		doc := &Doc{
			ClickId: clickId,
			Name: "qiche",
			CreateTime: time.Now().Unix(),
		}
		if err := ago.Collection.Insert(doc); err != nil {
			t.Error(err.Error())
		} else {
			t.Logf("%v", *doc)
		}
		res := ago.Collection.FindId(clickId)
		rdDoc := &Doc{}
		if err := res.One(rdDoc); err != nil {
			t.Errorf("not found, err %s", err.Error())
		} else {
			t.Log(rdDoc)
			t.Log(*rdDoc)
		}
	} else {
		t.Errorf("connect cache %s failed[%s]", addrs, err.Error())
	}
}

