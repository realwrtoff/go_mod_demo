package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/realwrtoff/go_mod_demo/internal/cache"
	"net/http"
	"net/url"
)

type ChannelReq struct {
	cache.PubChannelKey
	cache.CidInfo
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

	key := &cache.PubChannelKey {
		Pub: req.Pub,
		Cid: req.Cid,
	}
	value, ok := s.pubCidCfg.Get(*key)
	if !ok {
		// urldecode 广告主url
		unescapeUrl, err := url.QueryUnescape(req.AdvertiserAddr)
		if err != nil {
			res.Message = err.Error()
			return req, res, http.StatusInternalServerError, nil
		}
		// 初始化
		cfg := &cache.CidInfo{
			Status:         req.Status,
			Counter:        req.Counter,
			Step:           req.Step,
			AdvertiserAddr: unescapeUrl,
			AdvertiserCid:  req.AdvertiserCid,
			AppId: req.AppId,
			MyName:         req.MyName,
			BillingType:    req.BillingType,
		}
		s.pubCidCfg.Set(*key, cfg)
		err = s.pubCidMgoKv.Set(key, cfg)
		if err != nil {
			res.Message = fmt.Sprintf("insert failed err[%s]", err.Error())
		} else {
			res.Message = fmt.Sprintf("insert pub %s cid %s advertiser addr %s, step %d ok", req.Pub, req.Cid, cfg.AdvertiserAddr, cfg.Step)
		}
	} else {
		cidInfo := value.(*cache.CidInfo)
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
		err := s.pubCidMgoKv.Set(key, cidInfo)
		if err != nil {
			res.Message = fmt.Sprintf("update failed err[%s]", err.Error())
		} else {
			res.Message = fmt.Sprintf("update pub %s cid %s advertiser addr %s, step %d ok", req.Pub, req.Cid, cidInfo.AdvertiserAddr, cidInfo.Step)
		}
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

	key := &cache.PubChannelKey {
		Pub: req.Pub,
		Cid: req.Cid,
	}
	value, ok := s.pubCidCfg.Get(*key)
	if ok {
		cidInfo := value.(*cache.CidInfo)
		res.Message = fmt.Sprintf("%v", cidInfo)
	} else {
		res.Message = fmt.Sprintf("%v not found", key)
	}
	return req, res, http.StatusOK, nil
}
