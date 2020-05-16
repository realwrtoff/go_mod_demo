package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"time"
)

type ClickReq struct {
	Token string `form:"pub" bson:"pub" json:"pub"`
	CallBack string `form:"callback" bson:"callback" json:"callback"`
	Cid string `form:"cid" bson:"cid" json:"cid"`
	Sid string `form:"sid" bson:"sid,omitempty" json:"sid,omitempty"`
	Ip string `form:"ip" bson:"ip,omitempty" json:"ip,omitempty"`
	Os string `form:"os" bson:"os,omitempty" json:"os,omitempty"`
	OsVersion string `form:"osversion" bson:"osversion,omitempty" json:"os_version,omitempty"`
	DeviceType string `form:"devicetype" bson:"devicetype,omitempty" json:"device_type,omitempty"`
	Idfa string `form:"idfa" bson:"idfa,omitempty" json:"idfa,omitempty"`
	IdfaMd5 string `form:"idfa_md5" bson:"idfa_md5,omitempty" json:"idfa_md5,omitempty"`
	Imei string `form:"imei" bson:"imei,omitempty" json:"imei,omitempty"`
	ImeiMd5 string `form:"idfa_md5" bson:"imei_md5,omitempty" json:"imei_md5,omitempty"`
	AndroidId string `form:"androidid" bson:"androidid,omitempty" json:"androidid,omitempty"`
	AndroidIdMd5 string `form:"androidid_md5" bson:"androidid_md5,omitempty" json:"androidid_md5,omitempty"`
	AdvertiserId string `form:"advertiserid" bson:"advertiserid,omitempty" json:"advertiser_id,omitempty"`
	AdvertiserIdMd5 string `form:"advertiserid_md5" bson:"advertiserid_md5,omitempty" json:"advertiserid_md5,omitempty"`
}

type ClickRes struct {
	Code int `form:"code" json:"code"`
	Message string `form:"message" json:"message"`
	Data *ClickReq `form:"data" json:"data"`
}

// 需要确认Click存储时是否嵌套
type ClickInstall struct {
	ClickId bson.ObjectId `bson:"_id" json:"click_id"`
	ReqTime int64 `bson:"req_time,omitempty" json:"req_time,omitempty"`
	RespTime int64 `bson:"resp_time,omitempty" json:"resp_time,omitempty"`
	Reduce bool `bson:"reduce,omitempty" json:"reduce,omitempty"`
	Token string `form:"pub" bson:"pub" json:"pub"`
	CallBack string `form:"callback" bson:"callback" json:"callback"`
	Cid string `form:"cid" bson:"cid" json:"cid"`
	Sid string `form:"sid" bson:"sid,omitempty" json:"sid,omitempty"`
	Ip string `form:"ip" bson:"ip,omitempty" json:"ip,omitempty"`
	Os string `form:"os" bson:"os,omitempty" json:"os,omitempty"`
	OsVersion string `form:"osversion" bson:"osversion,omitempty" json:"os_version,omitempty"`
	DeviceType string `form:"devicetype" bson:"devicetype,omitempty" json:"device_type,omitempty"`
	Idfa string `form:"idfa" bson:"idfa,omitempty" json:"idfa,omitempty"`
	IdfaMd5 string `form:"idfa_md5" bson:"idfa_md5,omitempty" json:"idfa_md5,omitempty"`
	Imei string `form:"imei" bson:"imei,omitempty" json:"imei,omitempty"`
	ImeiMd5 string `form:"idfa_md5" bson:"imei_md5,omitempty" json:"imei_md5,omitempty"`
	AndroidId string `form:"androidid" bson:"androidid,omitempty" json:"androidid,omitempty"`
	AndroidIdMd5 string `form:"androidid_md5" bson:"androidid_md5,omitempty" json:"androidid_md5,omitempty"`
	AdvertiserId string `form:"advertiserid" bson:"advertiserid,omitempty" json:"advertiser_id,omitempty"`
	AdvertiserIdMd5 string `form:"advertiserid_md5" bson:"advertiserid_md5,omitempty" json:"advertiserid_md5,omitempty"`
}

func (s *Service) Click(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ClickReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	res := &ClickRes{
		Code: 200,
		Message: "",
		Data: req,
	}
	// 根据token和campaign 查找相关信息
	if s.channel[req.Token] == nil{
		res.Code = http.StatusNotFound
		res.Message = fmt.Sprintf("channel token %v not found", req.Token)
		s.warnLog.Warn(res.Message)
		return req, res, http.StatusNotFound, nil
	}
	if s.channel[req.Token][req.Cid] == nil {
		res.Code = http.StatusNotFound
		res.Message = fmt.Sprintf("channel %v cid %v not found", req.Token, req.Cid)
		s.warnLog.Warn(res.Message)
		return req, res, http.StatusNotFound, nil
	}

	// 检查设备参数
	devId := CheckReqDevId(req)
	if len(devId) == 0 {
		res.Code = http.StatusNotFound
		res.Message = fmt.Sprintf("idfa or imei not found, plz check the request params")
		s.warnLog.Warnf("%s\nparams %v", res.Message, *req)
		return req, res, http.StatusNotFound, nil
	}

	// 写入mongo
	clickId := bson.NewObjectId()
	clickInfo := &ClickInstall{
		ClickId: clickId,
		Token: req.Token,
		CallBack: req.CallBack,
		Cid: req.Cid,
		Sid: req.Sid,
		Ip: req.Ip,
		Os: req.Os,
		OsVersion: req.OsVersion,
		DeviceType: req.DeviceType,
		Idfa: req.Idfa,
		IdfaMd5: req.IdfaMd5,
		Imei: req.Imei,
		ImeiMd5: req.ImeiMd5,
		AndroidId: req.AndroidId,
		AndroidIdMd5: req.AndroidIdMd5,
		AdvertiserId: req.AdvertiserId,
		AdvertiserIdMd5: req.AdvertiserIdMd5,
		ReqTime: time.Now().Unix(),
	}
	s.infoLog.Info(clickId.Hex())
	s.infoLog.Info("req ", *req)
	s.infoLog.Info("click info ", clickInfo)
	if err := s.mgo.Collection.Insert(clickInfo); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		s.warnLog.Error("insert mongo failed, doc ",*clickInfo, " error ", err.Error())
		return req, res, http.StatusInternalServerError, nil
	}

	// 按照甲方请求参数进行封装
	if err := s.RequestAdvertiser(req, clickId.Hex()); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		s.warnLog.Error("request advertiser failed, err [%v]", *clickInfo, err.Error())
		return req, res, http.StatusInternalServerError, nil
	}

	return req, res, http.StatusOK, nil
}

func (s *Service)RequestAdvertiser(req *ClickReq, clickId string) error{
	// 信息替换
	adReq := make(map[string]interface{})
	adReq["pub"] = s.channel[req.Token][req.Cid].Name
	adReq["cid"] = s.channel[req.Token][req.Cid].OriginCid
	callBack := fmt.Sprintf("http://%s/install?click_id=%s", s.domain, clickId)
	escapeUrl := url.QueryEscape(callBack)
	adReq["callback"] = escapeUrl

	// 参数传递
	adReq["ip"] = req.Ip
	adReq["os"] = req.Os
	adReq["osversion"] = req.OsVersion
	adReq["devicetype"] = req.DeviceType
	adReq["idfa"] = req.Idfa
	adReq["idfamd5"] = req.IdfaMd5
	adReq["imei"] = req.Imei
	adReq["imeimd5"] = req.ImeiMd5
	adReq["androidid"] = req.AndroidId
	adReq["androididmd5"] = req.AndroidIdMd5
	adReq["advertiserid"] = req.AdvertiserId
	adReq["advertiseridmd5"] = req.AdvertiserIdMd5

	adUrl := s.channel[req.Token][req.Cid].Url
	httpRes := HttpClients.GET(adUrl, nil, adReq, nil)
	if httpRes.Err != nil {
		s.warnLog.Errorf("request %s params %v, response [%d][%s]", adUrl, adReq, httpRes.Status, httpRes.Err.Error())
		return httpRes.Err
	}
	s.infoLog.Infof("request %s params %v, response [%d]", adUrl, adReq, httpRes.Status)
	return nil
}

func CheckReqDevId(req *ClickReq) string {
	if req.Os == "android" {
		if len(req.Imei) > 0 {
			return req.Imei
		} else if len(req.ImeiMd5) > 0 {
			return req.ImeiMd5
		} else if len(req.AdvertiserId) > 0 {
			return req.AdvertiserId
		} else if len(req.AdvertiserIdMd5) > 0 {
			return req.AdvertiserIdMd5
		} else if len(req.AndroidId) > 0 {
			return req.AndroidId
		} else if len(req.AndroidIdMd5) > 0 {
			return req.AndroidIdMd5
		}
	} else {
		if len(req.Idfa) > 0 {
			return req.Idfa
		} else {
			return req.IdfaMd5
		}
	}
	return ""
}
