package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type ActiveReq struct {
	ClickId string `form:"click_id"`
}

type ActiveRes struct {
	Message string `form:"message" json:"message"`
}

func (s *Service) Active(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ActiveReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("rid[%s] bind failed. err: [%v]", rid, err)
	}
	s.infoLog.Infof("active receive callback [%v]", req)

	res := &ActiveRes{}
	var query *mgo.Query
	objId := bson.ObjectIdHex(req.ClickId)
	query = s.mgo.Collection.FindId(objId)

	record := &Record{}
	if err := query.One(record); err != nil {
		res.Message = fmt.Sprintf("req [%v] not found. err: [%v]", req, err)
		s.warnLog.Warn(res.Message)
		return nil, res, http.StatusNotFound, err
	}
	// 重复发送
	if record.ActiveTime != 0 {
		res.Message = fmt.Sprintf("req [%v] duplicated", req)
		s.warnLog.Warn(res.Message)
		return nil, res, http.StatusAlreadyReported, fmt.Errorf("req [%v] duplicated callback active", req)
	}
	// 查找扣量回调等配置信息
	key := fmt.Sprintf("%s_%s", record.Pub, record.Cid)
	value, ok := s.pubCidCfg.Get(key)
	if !ok {
		res.Message = fmt.Sprintf("req [%v] key[%s] cid info not found", req, key)
		s.warnLog.Warn(res.Message)
		return nil, res, http.StatusNotFound, fmt.Errorf("req [%v] key[%s] cid info not found", req, key)
	}
	cidInfo := value.(*CidInfo)

	if cidInfo.BillingType == "active" {
		cidInfo.Counter += 1
		var reduce bool
		// 需要扣量
		if cidInfo.Step > 0 {
			if cidInfo.Counter % cidInfo.Step == 0 {
				reduce = true
			}
		}
		// 更新mongo
		if err := s.mgo.Collection.UpdateId(record.ClickId, bson.M{"$set": bson.M{"active_time": time.Now().Unix(), "reduce": reduce}}); err != nil{
			res.Message = err.Error()
			s.warnLog.Errorf("update req [%v] resp_time&reduce failed.err[%s]", req, res.Message)
			return nil, res, http.StatusInternalServerError, err
		}
		// 回调
		if !reduce {
			callBack := s.httpClient.GET(record.CallBack, nil, nil, nil)
			if callBack.Err != nil {
				s.warnLog.Infof("call back %s failed", record.CallBack)
				res.Message = callBack.Err.Error()
				return nil, res, http.StatusNotFound, fmt.Errorf(" %s call back %s faild", record.Cid, record.CallBack)
			}
		} else {
			s.infoLog.Infof("call back %s reduced", record.CallBack)
		}
	} else {
		// 更新mongo 应该不会走到这一支
		if err := s.mgo.Collection.UpdateId(record.ClickId, bson.M{"$set": bson.M{"active_time": time.Now().Unix()}}); err != nil{
			res.Message = err.Error()
			s.warnLog.Errorf("update req [%v] resp_time failed.err[%s]", req, res.Message)
			return nil, res, http.StatusInternalServerError, err
		}
	}
	res.Message = fmt.Sprintf("active req [%v] ok", req)
	return req, res, http.StatusOK, nil
}