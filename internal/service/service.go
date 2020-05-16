package service

import (
	"github.com/hpifu/go-kit/hhttp"
	"github.com/realwrtoff/go_mod_demo/internal/mongo"
	"github.com/sirupsen/logrus"
	"time"
)

var HttpClients *hhttp.HttpClient

func init()  {
	HttpClients = hhttp.NewHttpClient(20, 200*time.Millisecond, 200*time.Millisecond)
}

type Service struct {
	secure    bool
	domain    string
	mgo     *mongo.Mongo
	channel map[string]map[string]*ChannelCid
	infoLog   *logrus.Logger
	warnLog   *logrus.Logger
	accessLog *logrus.Logger
}

func (s *Service) SetLogger(infoLog, warnLog, accessLog *logrus.Logger) {
	s.infoLog = infoLog
	s.warnLog = warnLog
	s.accessLog = accessLog
}

func (s *Service) SetMongo(mgo *mongo.Mongo) {
	s.mgo = mgo
}

func NewService(
	secure bool,
	domain string,
) *Service {
	return &Service{
		secure:    secure,
		domain:    domain,
		channel: make(map[string]map[string]*ChannelCid),
		infoLog:   logrus.New(),
		warnLog:   logrus.New(),
		accessLog: logrus.New(),
	}
}
