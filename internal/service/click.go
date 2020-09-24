package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/realwrtoff/go_mod_demo/internal/cache"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	"time"
)

type ClickReq struct {
	Pub             string `form:"pub" bson:"pub" json:"pub"`											// 渠道
	CallBack        string `form:"callback" bson:"callback" json:"callback"`							// 渠道回调地址
	Cid             string `form:"cid" bson:"cid" json:"cid"`											// 产品id
	Ip              string `form:"ip" bson:"ip,omitempty" json:"ip,omitempty"`
	Os              string `form:"os" bson:"os,omitempty" json:"os,omitempty"`
	OsVersion       string `form:"osversion" bson:"osversion,omitempty" json:"os_version,omitempty"`
	DeviceType      string `form:"devicetype" bson:"devicetype,omitempty" json:"device_type,omitempty"`
	Idfa            string `form:"idfa" bson:"idfa,omitempty" json:"idfa,omitempty"`
	IdfaMd5         string `form:"idfa_md5" bson:"idfa_md5,omitempty" json:"idfa_md5,omitempty"`
	Imei            string `form:"imei" bson:"imei,omitempty" json:"imei,omitempty"`
	ImeiMd5         string `form:"idfa_md5" bson:"imei_md5,omitempty" json:"imei_md5,omitempty"`
	AndroidId       string `form:"androidid" bson:"androidid,omitempty" json:"androidid,omitempty"`
	AndroidIdMd5    string `form:"androidid_md5" bson:"androidid_md5,omitempty" json:"androidid_md5,omitempty"`
	AdvertiserId    string `form:"advertiserid" bson:"advertiserid,omitempty" json:"advertiser_id,omitempty"`
	AdvertiserIdMd5 string `form:"advertiserid_md5" bson:"advertiserid_md5,omitempty" json:"advertiserid_md5,omitempty"`
}

type ClickRes struct {
	Code    int       `form:"code" json:"code"`
	Message string    `form:"message" json:"message"`
	Data    interface{} `form:"data" json:"data"`
}

// 需要确认Click存储时是否嵌套
type Record struct {
	ClickId         bson.ObjectId `bson:"_id" json:"click_id"`
	DevId           string        `form:"dev_id" bson:"dev_id,omitempty" json:"dev_id, omitempty"`
	AppId           string        `form:"app_id" bson:"app_id,omitempty" json:"app_id, omitempty"`
	ReqTime         int64         `bson:"req_time,omitempty" json:"req_time,omitempty"`
	RespTime        int64         `bson:"resp_time,omitempty" json:"resp_time,omitempty"`
	ActiveTime      int64         `bson:"active_time,omitempty" json:"active_time,omitempty"`
	Reduce          bool          `bson:"reduce,omitempty" json:"reduce,omitempty"`
	Pub             string        `form:"pub" bson:"pub" json:"pub"`
	CallBack        string        `form:"callback" bson:"callback" json:"callback"`
	Cid             string        `form:"cid" bson:"cid" json:"cid"`
	Ip              string        `form:"ip" bson:"ip,omitempty" json:"ip,omitempty"`
	Os              string        `form:"os" bson:"os,omitempty" json:"os,omitempty"`
	OsVersion       string        `form:"osversion" bson:"osversion,omitempty" json:"os_version,omitempty"`
	DeviceType      string        `form:"devicetype" bson:"devicetype,omitempty" json:"device_type,omitempty"`
	Idfa            string        `form:"idfa" bson:"idfa,omitempty" json:"idfa,omitempty"`
	IdfaMd5         string        `form:"idfa_md5" bson:"idfa_md5,omitempty" json:"idfa_md5,omitempty"`
	Imei            string        `form:"imei" bson:"imei,omitempty" json:"imei,omitempty"`
	ImeiMd5         string        `form:"idfa_md5" bson:"imei_md5,omitempty" json:"imei_md5,omitempty"`
	AndroidId       string        `form:"androidid" bson:"androidid,omitempty" json:"androidid,omitempty"`
	AndroidIdMd5    string        `form:"androidid_md5" bson:"androidid_md5,omitempty" json:"androidid_md5,omitempty"`
	AdvertiserId    string        `form:"advertiserid" bson:"advertiserid,omitempty" json:"advertiser_id,omitempty"`
	AdvertiserIdMd5 string        `form:"advertiserid_md5" bson:"advertiserid_md5,omitempty" json:"advertiserid_md5,omitempty"`
}

type WeiYiRes struct {
	Code string `json:"code"`
	Flag string `json:"flag"`
	Message string `json:"message"`
	Data string `json:"data"`
}

func (s *Service) Click(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ClickReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("rid[%s] bind failed. err: [%v]", rid, err)
	}

	res := &ClickRes{
		Code:    200,
		Message: "",
	}

	// 检查设备参数
	devId := CheckReqDevId(req)
	if len(devId) == 0 {
		res.Code = http.StatusNotFound
		res.Message = fmt.Sprintf("idfa or imei not found, plz check the request params")
		s.warnLog.Warnf("%s\nparams %v", res.Message, *req)
		return req, res, http.StatusNotFound, nil
	}
	// 根据pub和cid 查找相关信息
	key := fmt.Sprintf("%s_%s", req.Pub, req.Cid)
	value, ok := s.pubCidCfg.Get(key)
	if !ok {
		res.Code = http.StatusNotFound
		res.Message = fmt.Sprintf("channel %s cid %s not found", req.Pub, req.Cid)
		s.warnLog.Warn(res.Message)
		return req, res, http.StatusNotFound, nil
	}
	cidInfo := value.(*cache.CidInfo)

	// 写入mongo
	clickId := bson.NewObjectId()
	clickInfo := &Record{
		ClickId:         clickId,
		DevId:           devId,
		AppId:           cidInfo.AppId,
		Pub:             req.Pub,
		CallBack:        req.CallBack,
		Cid:             req.Cid,
		Ip:              req.Ip,
		Os:              req.Os,
		OsVersion:       req.OsVersion,
		DeviceType:      req.DeviceType,
		Idfa:            req.Idfa,
		IdfaMd5:         req.IdfaMd5,
		Imei:            req.Imei,
		ImeiMd5:         req.ImeiMd5,
		AndroidId:       req.AndroidId,
		AndroidIdMd5:    req.AndroidIdMd5,
		AdvertiserId:    req.AdvertiserId,
		AdvertiserIdMd5: req.AdvertiserIdMd5,
		ReqTime:         time.Now().Unix(),
	}
	s.infoLog.Infof("receive pub_cid %s request click id %s.", key, clickId.Hex())

	if err := s.mgo.Collection.Insert(clickInfo); err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		s.warnLog.Error("insert cache failed, doc ", *clickInfo, " error ", err.Error())
		return req, res, http.StatusInternalServerError, nil
	}

	// 按照甲方请求参数进行封装
	httpRes := s.RequestAdvertiser(req, clickId.Hex(), cidInfo)
	if httpRes.Err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = httpRes.Err.Error()
		s.warnLog.Errorf("request advertiser failed [%v], err [%s]", *clickInfo, httpRes.Err.Error())
		return req, res, http.StatusInternalServerError, nil
	}
	res = s.DealResponse(httpRes, clickId.Hex(), cidInfo)
	return req, res, http.StatusOK, nil
}

func (s *Service) RequestAdvertiser(req *ClickReq, clickId string, cidCfg *cache.CidInfo) *hhttp.HttpResult {
	adReq := make(map[string]interface{})
	switch cidCfg.AdvertiserCid {
		case "weiyi":
			adReq["channelId"] = cidCfg.MyName
			adReq["callBackUrl"] = clickId
		default:
			adReq["pub"] = cidCfg.MyName
			adReq["cid"] = cidCfg.AdvertiserCid
			adReq["app_id"] = cidCfg.AppId
			// 潜规则 BillingType要和router一致
			callBack := fmt.Sprintf("http://%s/%s?click_id=%s", s.domain, cidCfg.BillingType, clickId)
			escapeUrl := url.QueryEscape(callBack)
			adReq["callback"] = escapeUrl
	}

	// 参数传递
	if req.Ip != "" {
		adReq["ip"] = req.Ip
	}
	if req.Os != "" {
		adReq["os"] = req.Os
	}
	if req.OsVersion != "" {
		adReq["osversion"] = req.OsVersion
	}
	if req.DeviceType != "" {
		adReq["devicetype"] = req.DeviceType
	}
	if req.Idfa != "" {
		adReq["idfa"] = req.Idfa
	}
	if req.IdfaMd5 != "" {
		adReq["idfamd5"] = req.IdfaMd5
	}
	if req.Imei != "" {
		adReq["imei"] = req.Imei
	}
	if req.ImeiMd5 != "" {
		adReq["imeimd5"] = req.ImeiMd5
	}
	if req.AndroidId != "" {
		adReq["androidid"] = req.AndroidId
	}
	if req.AndroidIdMd5 != "" {
		adReq["androididmd5"] = req.AndroidIdMd5
	}
	if req.AdvertiserId != "" {
		adReq["advertiserid"] = req.AdvertiserId
	}
	if req.AdvertiserIdMd5 != "" {
		adReq["advertiseridmd5"] = req.AdvertiserIdMd5
	}

	adUrl := cidCfg.AdvertiserAddr
	s.infoLog.Info(adUrl)
	s.infoLog.Info(adReq)
	httpRes := s.httpClient.GET(adUrl, nil, adReq, nil)
	return httpRes
}

func (s *Service) DealResponse(httpRes *hhttp.HttpResult, clickId string, cidCfg *cache.CidInfo) *ClickRes {
	res := &ClickRes{
		Code:    200,
		Message: "",
	}
	switch cidCfg.AdvertiserCid {
	case "weiyi":
		wy := WeiYiRes{}
		_ = json.Unmarshal(httpRes.Res, &wy)
		if wy.Code != "0" {
			res.Code = -1
			res.Message = wy.Message
		}
		res.Data = wy.Data
	default:
		res.Data = clickId
	}
	return res
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
