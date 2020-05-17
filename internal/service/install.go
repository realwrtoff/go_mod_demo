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
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	objId := bson.ObjectIdHex(req.ClickId)
	rec := s.mgo.Collection.FindId(objId)
	clkIns := &ClickInstall{}
	if err := rec.One(clkIns); err != nil {
		return nil, nil, http.StatusNotFound, fmt.Errorf("clickid [%v] not found. err: [%v]", req.ClickId, err)
	}
	// 重复发送
	if clkIns.RespTime != 0 {
		return nil, nil, http.StatusAlreadyReported, fmt.Errorf("clickid [%v] duplicated", req.ClickId)
	}
	// 查找扣量回调等配置信息
	if s.channel[clkIns.Token] == nil {
		return nil, nil, http.StatusNotFound, fmt.Errorf("click token [%v] not found", clkIns.Token)
	}
	if s.channel[clkIns.Token][clkIns.Cid] == nil {
		return nil, nil, http.StatusNotFound, fmt.Errorf("click token[%v] cid[%v] not found", clkIns.Token, clkIns.Cid)
	}

	if s.channel[clkIns.Token][clkIns.Cid].Billing == "install" {
		s.channel[clkIns.Token][clkIns.Cid].Counter += 1
		var reduce bool
		// 需要扣量
		if s.channel[clkIns.Token][clkIns.Cid].Step > 0 {
			if s.channel[clkIns.Token][clkIns.Cid].Counter % s.channel[clkIns.Token][clkIns.Cid].Step == 0 {
				reduce = true
			}
		}
		// 更新mongo
		if err := s.mgo.Collection.UpdateId(objId, bson.M{"$set": bson.M{"resp_time": time.Now().Unix(), "reduce": reduce}}); err != nil{
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
		// 更新mongo
		if err := s.mgo.Collection.UpdateId(objId, bson.M{"$set": bson.M{"resp_time": time.Now().Unix()}}); err != nil{
			return nil, nil, http.StatusInternalServerError, err
		}
	}
	return req, &InstallRes{
		Message: "ok",
	}, http.StatusOK, nil
}