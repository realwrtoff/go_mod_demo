package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

type CidInfo struct {
	Status         int    `form:"status" json:"status,omitempty"`                   // 单子状态
	Counter        int    `form:"counter" json:"counter,omitempty"`                 // 单子计数器
	Step           int    `form:"step" json:"step,omitempty"`                       // 单子步长
	AdvertiserAddr string `form:"advertiser_addr" json:"advertiser_addr,omitempty"` // 单子广告主的请求地址
	AdvertiserCid  string `form:"advertiser_cid" json:"advertiser_cid,omitempty"`   // 请求广告主的身份
	AppId          string `form:"app_id" json:"app_id,omitempty"`                   // 广告ID
	MyName         string `form:"my_name" json:"my_name,omitempty"`                 // 请求广告主的身份
	BillingType    string `form:"billing_type" json:"billing_type,omitempty"`       // 计费类型 install, active ?
}

type ChannelReq struct {
	Pub     string `form:"pub"  json:"pub,omitempty"` // 渠道
	Cid     string `form:"cid" json:"cid,omitempty"`  // 单子id
	CidInfo        // 渠道单子信息
}

type ChannelRes struct {
	Message string `form:"message" json:"message"`
}

// 设置渠道和单子的映射， 存储单子相关信息
func (s *Service) Channel(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ChannelReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("rid[%s] bind failed. err: [%v]", rid, err)
	}

	res := &ChannelRes{
		Message: "",
	}

	key := fmt.Sprintf("%s_%s", req.Pub, req.Cid)
	value, ok := s.pubCidCfg.Get(key)
	if !ok {
		// urldecode 广告主url
		unescapeUrl, err := url.QueryUnescape(req.AdvertiserAddr)
		if err != nil {
			res.Message = err.Error()
			return req, res, http.StatusInternalServerError, nil
		}
		// 初始化
		cfg := &CidInfo{
			Status:         req.Status,
			Counter:        req.Counter,
			Step:           req.Step,
			AdvertiserAddr: unescapeUrl,
			AdvertiserCid:  req.AdvertiserCid,
			AppId: req.AppId,
			MyName:         req.MyName,
			BillingType:    req.BillingType,
		}
		s.pubCidCfg.Set(key, cfg)
		res.Message = fmt.Sprintf("add pub %s cid %s advertiser addr %s ok", req.Pub, req.Cid, req.AdvertiserAddr)
	} else {
		cidInfo := value.(*CidInfo)
		// 更新
		res.Message = fmt.Sprintf("update pub %s cid %s set", req.Pub, req.Cid)
		if req.Status != 0 {
			cidInfo.Status = req.Status
			res.Message += fmt.Sprintf(" status=[%d]", req.Status)
		}
		if req.Step != 0 {
			cidInfo.Step = req.Step
			res.Message += fmt.Sprintf(" step=[%d]", req.Step)
		}
		if len(req.AdvertiserAddr) != 0 {
			unescapeUrl, err := url.QueryUnescape(req.AdvertiserAddr)
			if err != nil {
				res.Message = err.Error()
				return req, res, http.StatusOK, nil
			}
			cidInfo.AdvertiserAddr = unescapeUrl
			res.Message += fmt.Sprintf(" advertiser_url=[%s]", unescapeUrl)
		}
		if len(req.AdvertiserCid) != 0 {
			cidInfo.AdvertiserCid = req.AdvertiserCid
			res.Message += fmt.Sprintf(" advertiser_cid=[%s]", req.AdvertiserCid)
		}
		if len(req.AppId) != 0 {
			cidInfo.AppId = req.AppId
			res.Message += fmt.Sprintf(" app_id=[%s]", req.AppId)
		}
		if len(req.MyName) != 0 {
			cidInfo.MyName = req.MyName
			res.Message += fmt.Sprintf(" my_name=[%s]", req.MyName)
		}
		if len(req.BillingType) != 0 {
			cidInfo.BillingType = req.BillingType
			res.Message += fmt.Sprintf(" billing_type=[%s]", req.BillingType)
		}
		// 无需写回，已经更新
		// s.pubCidCfg.Set(key, cidInfo)
	}
	return req, res, http.StatusOK, nil
}

// 设置渠道和单子的映射， 存储单子相关信息
func (s *Service) GetChannel(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ChannelReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("rid[%s] bind failed. err: [%v]", rid, err)
	}

	res := &ChannelRes{
		Message: "",
	}

	key := fmt.Sprintf("%s_%s", req.Pub, req.Cid)
	value, ok := s.pubCidCfg.Get(key)
	if ok {
		cidInfo := value.(*CidInfo)
		res.Message = fmt.Sprintf("addr %s", cidInfo.AdvertiserAddr)
	} else {
		res.Message = key + " not found"
	}
	return req, res, http.StatusOK, nil
}
