package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

const (
	INSTALL = 1
	ACTIVE
)

type ChannelCid struct {
	Status    int    `form:"status" json:"status,omitempty"`         // 单子状态
	Counter   int    `form:"cnt" json:"cnt,omitempty"`               // 单子计数器
	Step      int    `form:"step" json:"step,omitempty"`             // 单子步长
	Url       string `form:"url" json:"url,omitempty"`               // 单子广告主的请求地址
	OriginCid string `form:"origin_cid" json:"origin_cid,omitempty"` // 请求广告主的身份
	Name      string `form:"name" json:"name,omitempty"`             // 请求广告主的身份
	Billing   int    `form:"billing" json:"billing,omitempty"`       // 计费类型
}

type ChannelReq struct {
	Token      string `form:"pub"  json:"pub,omitempty"` // 渠道token
	CampaignId string `form:"cid" json:"cid,omitempty"`  // 单子id
	ChannelCid        // 渠道单子信息
}

type ChannelRes struct {
	Message string     `form:"message" json:"message"`
	Channel ChannelReq `form:"channel" json:"channel"`
}

// 设置渠道和单子的映射， 存储单子相关信息
func (s *Service) Channel(rid string, c *gin.Context) (interface{}, interface{}, int, error) {
	req := &ChannelReq{}

	if err := c.Bind(req); err != nil {
		return nil, nil, http.StatusBadRequest, fmt.Errorf("bind failed. err: [%v]", err)
	}

	if s.channel[req.Token] == nil {
		s.channel[req.Token] = make(map[string]*ChannelCid)
	}
	if s.channel[req.Token][req.CampaignId] == nil {
		// urldecode 广告主url
		unescapeUrl, err := url.QueryUnescape(req.Url)
		if err != nil {
			return req, &ChannelRes{
				Message: err.Error(),
				Channel: *req,
			}, http.StatusOK, nil
		}
		// 初始化
		s.channel[req.Token][req.CampaignId] = &ChannelCid{
			Status:  req.Status,
			Counter: req.Counter,
			Step:    req.Step,
			Url:     unescapeUrl,
			Name:    req.Name,
		}
	} else {
		// 更新
		if req.Status != 0 {
			s.channel[req.Token][req.CampaignId].Status = req.Status
		}
		if req.Step != 0 {
			s.channel[req.Token][req.CampaignId].Step = req.Step
		}
		if len(req.Url) != 0 {
			unescapeUrl, err := url.QueryUnescape(req.Url)
			if err != nil {
				return req, &ChannelRes{
					Message: err.Error(),
					Channel: *req,
				}, http.StatusOK, nil
			}
			s.channel[req.Token][req.CampaignId].Url = unescapeUrl
		}
		if len(req.Name) != 0 {
			s.channel[req.Token][req.CampaignId].Name = req.Name
		}
	}
	return req, &ChannelRes{
		Message: "update",
		Channel: *req,
	}, http.StatusOK, nil
}
