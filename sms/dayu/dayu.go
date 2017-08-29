package dayu

import (
	"go.uber.org/zap"

	"github.com/empirefox/cement/clog"
	"github.com/empirefox/cement/perr"
	"github.com/empirefox/cement/sms"
	"github.com/opensource-conet/alidayu"
)

type Config struct {
	Appkey    string `validate:"required"`
	AppSecret string `validate:"required"`
	SignName  string `validate:"required"`
	Template  string `validate:"required"`
}

type dayu struct {
	config Config
	log    *zap.Logger
}

func NewDayu(config Config, logger clog.Logger) sms.Vendor {
	alidayu.Appkey = config.Appkey
	alidayu.AppSecret = config.AppSecret
	return &dayu{
		config: config,
		log:    logger.Module("alidayu"),
	}
}

func (s *dayu) Send(phone, code string) error {
	res, err := alidayu.SendOnce(phone, s.config.SignName, s.config.Template, `{"code":"`+code+`"}`)
	if err != nil {
		s.log.Error("send fail", zap.Error(err))
		return perr.SmsSendFail
	}
	if !res.Success {
		s.log.Error("send fail", zap.Int("code", res.ResultError.Code), zap.String("subcode", res.ResultError.SubCode))
		return perr.SmsSendFail
	}
	return nil
}
