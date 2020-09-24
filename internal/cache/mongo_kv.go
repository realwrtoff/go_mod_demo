package cache

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type PubChannelKey struct {
	Pub            string `form:"pub" bson:"pub" json:"pub,omitempty"`                                     // 渠道
	Cid            string `form:"cid" bson:"cid" json:"cid,omitempty"`                                     // 单子id
}

type CidInfo struct {
	Status         int    `form:"status" json:"status,omitempty"`                   // 单子状态
	Counter        int    `form:"counter" json:"counter,omitempty"`                 // 单子计数器
	Step           int    `form:"step" json:"step,omitempty"`                       // 单子步长
	AdvertiserAddr string `form:"advertiser_addr" json:"advertiser_addr,omitempty"` // 单子广告主的请求地址
	AdvertiserCid  string `form:"advertiser_cid" json:"advertiser_cid,omitempty"`   // 广告主侧单子id
	AppId          string `form:"app_id" json:"app_id,omitempty"`                   // 广告ID
	MyName         string `form:"my_name" json:"my_name,omitempty"`                 // 请求广告主的身份
	BillingType    string `form:"billing_type" json:"billing_type,omitempty"`       // 计费类型 install, active ?
}

type PubChannelS struct {
	Pub            string `form:"pub" bson:"pub" json:"pub,omitempty"`                                     // 渠道
	Cid            string `form:"cid" bson:"cid" json:"cid,omitempty"`
	Status         int    `form:"status" json:"status,omitempty"`                   // 单子状态
	Counter        int    `form:"counter" json:"counter,omitempty"`                 // 单子计数器
	Step           int    `form:"step" json:"step,omitempty"`                       // 单子步长
	AdvertiserAddr string `form:"advertiser_addr" json:"advertiser_addr,omitempty"` // 单子广告主的请求地址
	AdvertiserCid  string `form:"advertiser_cid" json:"advertiser_cid,omitempty"`   // 广告主侧单子id
	AppId          string `form:"app_id" json:"app_id,omitempty"`                   // 广告ID
	MyName         string `form:"my_name" json:"my_name,omitempty"`                 // 请求广告主的身份
	BillingType    string `form:"billing_type" json:"billing_type,omitempty"`       // 计费类型 install, active ?
}

func NewMgoKv(mgo *Mongo) *MgoKv {
	return &MgoKv{
		kvs: mgo,
	}
}

func (m *MgoKv)Load(memKv *MemKv) error {
	res := m.kvs.Collection.Find(bson.M {"status": 1})
	num, err := res.Count()
	if err != nil {
		return err
	} else if num == 0 {
		return nil
	}
	var elems = make([]*PubChannelS, num)
	if err = res.All(&elems); err != nil {
		return err
	}
	// 计数器启动时候清零开始
	for _, elem := range elems {
		key := &PubChannelKey{
			Pub: elem.Pub,
			Cid: elem.Cid,
		}
		cidInfo := &CidInfo{
			Status: elem.Status,
			Counter: 0,
			Step: elem.Step,
			AdvertiserAddr: elem.AdvertiserAddr,
			AdvertiserCid: elem.AdvertiserCid,
			AppId: elem.AppId,
			MyName: elem.MyName,
			BillingType: elem.BillingType,
		}
		memKv.Set(*key, cidInfo)
	}
	return nil
}

func (m *MgoKv) Get(key interface{}) (interface{}, error) {
	pc := key.(*PubChannelKey)
	query := m.kvs.Collection.Find(bson.M {"pub": pc.Pub, "cid": pc.Cid})
	record := &CidInfo{}
	if err := query.One(record); err != nil {
		return nil, err
	}
	return record, nil
}

func (m *MgoKv) Set(key interface{}, value interface{}) error {
	pk := key.(*PubChannelKey)
	cidInfo := value.(*CidInfo)
	pc := &PubChannelS{}
	pc.Pub = pk.Pub
	pc.Cid = pk.Cid
	rec, err := m.Get(pk)
	if err != nil && err != mgo.ErrNotFound {
		return err
	} else if err == mgo.ErrNotFound {
		pc.Status = cidInfo.Status
		pc.Step = cidInfo.Step
		pc.AdvertiserAddr = cidInfo.AdvertiserAddr
		pc.AdvertiserCid = cidInfo.AdvertiserCid
		pc.AppId = cidInfo.AppId
		pc.MyName = cidInfo.MyName
		pc.BillingType = cidInfo.BillingType
	} else {
		recInfo := rec.(*CidInfo)
		// 更新状态， 一段丑陋且无奈的代码
		if cidInfo.Status != 0 {
			pc.Status = cidInfo.Status
		} else {
			pc.Status = recInfo.Status
		}
		// 更新步长
		pc.Step = cidInfo.Step

		if cidInfo.AdvertiserAddr != "" {
			pc.AdvertiserAddr = cidInfo.AdvertiserAddr
		} else {
			pc.AdvertiserAddr = recInfo.AdvertiserAddr
		}
		if cidInfo.AdvertiserCid != "" {
			pc.AdvertiserCid = cidInfo.AdvertiserCid
		} else {
			pc.AdvertiserCid = recInfo.AdvertiserCid
		}
		if cidInfo.AppId != "" {
			pc.AppId = cidInfo.AppId
		} else {
			pc.AppId = recInfo.AppId
		}
		if cidInfo.MyName != "" {
			pc.MyName = cidInfo.MyName
		} else {
			pc.MyName = recInfo.MyName
		}
		if cidInfo.BillingType != "" {
			pc.BillingType = cidInfo.BillingType
		} else {
			pc.BillingType = recInfo.BillingType
		}
	}
	selector := bson.M{"pub": pc.Pub, "cid": pc.Cid}
	if _, err := m.kvs.Collection.Upsert(selector, pc); err != nil {
		return err
	}
	return nil
}
