package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
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
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	objId := bson.ObjectIdHex(req.ClickId)
	rec := s.mgo.Collection.FindId(objId)
	clkIns := &ClickInstall{}
	if err := rec.One(clkIns); err != nil {
		return nil, nil, http.StatusNotFound, fmt.Errorf("clickid [%v] not found. err: [%v]", req.ClickId, err)
	}
	// 重复发送
	if clkIns.ActiveTime != 0 {
		return nil, nil, http.StatusAlreadyReported, fmt.Errorf("clickid [%v] duplicated callback active", req.ClickId)
	}
	// 查找扣量回调等配置信息
	if s.channel[clkIns.Token] == nil {
		return nil, nil, http.StatusNotFound, fmt.Errorf("active token [%v] not found", clkIns.Token)
	}
	if s.channel[clkIns.Token][clkIns.Cid] == nil {
		return nil, nil, http.StatusNotFound, fmt.Errorf("active token[%v] cid[%v] not found", clkIns.Token, clkIns.Cid)
	}

	if s.channel[clkIns.Token][clkIns.Cid].Billing == "active" {
		s.channel[clkIns.Token][clkIns.Cid].Counter += 1
		var reduce bool
		// 需要扣量
		if s.channel[clkIns.Token][clkIns.Cid].Step > 0 {
			if s.channel[clkIns.Token][clkIns.Cid].Counter % s.channel[clkIns.Token][clkIns.Cid].Step == 0 {
				reduce = true
			}
		}
		// 更新mongo
		if err := s.mgo.Collection.UpdateId(objId, bson.M{"$set": bson.M{"active_time": time.Now().Unix(), "reduce": reduce}}); err != nil{
			return nil, nil, http.StatusInternalServerError, err
		}
		// 回调
		if !reduce {
			callBack := HttpClients.GET(clkIns.CallBack, nil, nil, nil)
			if callBack.Err != nil {
				s.warnLog.Infof("call back %s failed", clkIns.CallBack)
				return nil, nil, http.StatusNotFound, fmt.Errorf("call back %s faild", clkIns.Cid)
			}
		}
	} else {
		// 更新mongo 应该不会走到这一支
		if err := s.mgo.Collection.UpdateId(objId, bson.M{"$set": bson.M{"active_time": time.Now().Unix()}}); err != nil{
			return nil, nil, http.StatusInternalServerError, err
		}
	}
	return req, &ActiveRes{
		Message: "ok",
	}, http.StatusOK, nil
}