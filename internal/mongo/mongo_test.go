package cache

import (
	"testing"
)

func TestMongo_Connect(t *testing.T) {
	addrs := []string{"127.0.0.1:27017"}
	ago := &Mongo{}
	ago.Addrs = addrs
	ago.Username = "tim"
	ago.Password = "Spurs21"
	ago.DbName = "emar"
	ago.CollectionName = "campaign"
	ago.TimeOut = 30
	if err := ago.Connect(); err != nil {
		t.Errorf("connect mongdb [%s] failed [%s]", ago.Addrs, err.Error())
	} else {
		t.Logf("connect mongo %s success\n", addrs)
	}
	ago.Close()
}

func TestMongo_Ping(t *testing.T) {
	addrs := []string{"127.0.0.1:27017"}
	ago := &Mongo{}
	ago.Addrs = addrs
	ago.Username = "tim"
	ago.Password = "Spurs21"
	ago.DbName = "emar"
	ago.CollectionName = "campaign"
	ago.TimeOut = 30
	if err := ago.Connect(); err == nil {
		if err = ago.Ping(); err != nil {
			t.Errorf("ping mongo %s failed[%s]", addrs, err.Error())
		} else {
			t.Logf("ping mongo %s success\n", addrs)
		}
	} else {
		t.Errorf("connect mongo %s failed[%s]", addrs, err.Error())
	}
	ago.Close()
}
