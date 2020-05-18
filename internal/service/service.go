package service

import (
	"github.com/hpifu/go-kit/hhttp"
	"github.com/realwrtoff/go_mod_demo/internal/cache"
	"github.com/sirupsen/logrus"
)


type Service struct {
	secure    bool
	domain    string
	mgo     *cache.Mongo
	pubCidCfg   *cache.MemKv
	httpClient  *hhttp.HttpClient
	infoLog   *logrus.Logger
	warnLog   *logrus.Logger
	accessLog *logrus.Logger
}

func (s *Service) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	s.infoLog = infoLog
	s.warnLog = warnLog
	s.accessLog = accessLog
}

func NewService(
	secure bool,
	domain string,
	mgo *cache.Mongo,
	pubCidCfg *cache.MemKv,
	httpClient *hhttp.HttpClient,
) *Service {
	return &Service{
		secure:    secure,
		domain:    domain,
		mgo: mgo,
		pubCidCfg: pubCidCfg,
		httpClient: httpClient,
		infoLog:   logrus.New(),
		warnLog:   logrus.New(),
		accessLog: logrus.New(),
	}
}