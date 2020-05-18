package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type InstallReq struct {
	ClickId string `form:"click_id"`
}

type InstallRes struct {
	Message string `form:"message" json:"message"`
}

func (s *Service) Install(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &InstallReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("rid[%s] bind failed. err: [%v]", rid, err)
	}
	s.infoLog.Infof("receive click_id %s install", req.ClickId)

	res := &InstallRes{}
	objId := bson.ObjectIdHex(req.ClickId)
	rec := s.mgo.Collection.FindId(objId)
	record := &Record{}
	if err := rec.One(record); err != nil {
		res.Message = fmt.Sprintf("clickid [%v] not found. err: [%v]", req.ClickId, err)
		s.warnLog.Warn(res.Message)
		return nil, res, http.StatusNotFound, err
	}
	// 重复发送
	if record.RespTime != 0 {
		res.Message = fmt.Sprintf("clickid [%v] duplicated", req.ClickId)
		s.warnLog.Warn(res.Message)
		return nil, res, http.StatusAlreadyReported, fmt.Errorf("clickid [%v] duplicated", req.ClickId)
	}
	// 查找扣量回调等配置信息
	key := fmt.Sprintf("%s_%s", record.Pub, record.Cid)
	value, ok := s.pubCidCfg.Get(key)
	if !ok {
		res.Message = fmt.Sprintf("clickid [%v] key[%s] cid info not found", req.ClickId, key)
		s.warnLog.Warn(res.Message)
		return nil, res, http.StatusNotFound, fmt.Errorf("clickid [%v] key[%s] cid info not found", req.ClickId, key)
	}
	cidInfo := value.(*CidInfo)

	if cidInfo.BillingType == "install" {
		cidInfo.Counter += 1
		var reduce bool
		// 需要扣量
		if cidInfo.Step > 0 {
			if cidInfo.Counter % cidInfo.Step == 0 {
				reduce = true
			}
		}
		// 更新mongo
		if err := s.mgo.Collection.UpdateId(objId, bson.M{"$set": bson.M{"resp_time": time.Now().Unix(), "reduce": reduce}}); err != nil{
			res.Message = err.Error()
			s.warnLog.Errorf("update click_id %s resp_time & reduce failed.err[%s]", req.ClickId, res.Message)
			return nil, res, http.StatusInternalServerError, err
		}
		// 回调
		if !reduce {
			callBack := s.httpClient.GET(record.CallBack, nil, nil, nil)
			if callBack.Err != nil {
				s.warnLog.Infof("call back %s failed", record.CallBack)
				res.Message = callBack.Err.Error()
				return nil, res, http.StatusNotFound, fmt.Errorf("call back %s faild", record.Cid)
			}
		} else {
			s.infoLog.Infof("call back %s reduced", record.CallBack)
		}
	} else {
		// 更新mongo
		if err := s.mgo.Collection.UpdateId(objId, bson.M{"$set": bson.M{"resp_time": time.Now().Unix()}}); err != nil{
			res.Message = err.Error()
			s.warnLog.Errorf("update click_id %s resp_time failed.err[%s]", req.ClickId, res.Message)
			return nil, res, http.StatusInternalServerError, err
		}
	}
	res.Message = fmt.Sprintf("install %s ok", req.ClickId)
	return req, res, http.StatusOK, nil
}