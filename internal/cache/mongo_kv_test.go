package cache

import (
	"gopkg.in/mgo.v2"
	"testing"
)

func TestMgoKv_Get(t *testing.T) {
	addrs := []string{"127.0.0.1:21617"}
	ago := &Mongo{}

	ago.Addrs = addrs
	ago.DbName = "haina"
	ago.CollectionName = "pub_cid_info"
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

	mgoKv := NewMgoKv(ago)
	key := &PubChannelKey{
		Pub: "jim",
		Cid: "money",
	}
	res, err := mgoKv.Get(key)
	if err != nil && err != mgo.ErrNotFound{
		t.Errorf("cache get %v failed err[%s]", key, err.Error())
	} else if err == mgo.ErrNotFound {
		t.Logf("get %v not found\n", key)
	} else {
		t.Logf("get %v success [%v]\n", key, res)
	}
	ago.Close()
}

func TestMemKv_Set(t *testing.T) {
	addrs := []string{"127.0.0.1:21617"}
	ago := &Mongo{}

	ago.Addrs = addrs
	ago.DbName = "haina"
	ago.CollectionName = "pub_cid_info"
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

	mgoKv := NewMgoKv(ago)
	key := &PubChannelKey{
		Pub: "jim",
		Cid: "money",
	}
	val := &CidInfo{
		Status: 1,
		AdvertiserAddr: "http://advertise.main.com",
		MyName: "whoami",
		BillingType: "install",
	}
	if err := mgoKv.Set(key, val); err != nil {
		t.Errorf("set %v failed[%s]", key, err.Error())
	} else {
		t.Logf("set %v success\n", key)
	}

	ago.Close()
}
