package sms

import (
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/empirefox/esecend/cerr"
	"github.com/empirefox/esecend/config"
	"github.com/golang/glog"
	"github.com/opensource-conet/alidayu"
	"github.com/patrickmn/go-cache"
)

const (
	BindPhone = "B"
	SetPaykey = "P"
)

type Sender interface {
	Send(prefix string, userId uint, phone string) error
	Verify(prefix string, userId uint, phone, code string) bool
}

type LimitedCode struct {
	ID     string
	Code   string
	UserID uint
	GenAt  int64
}

type sender struct {
	config    *config.Alidayu
	cache     *cache.Cache
	retryMin  time.Duration
	codeChars []byte
}

func NewSender(config *config.Config) Sender {
	dayu := &config.Alidayu
	alidayu.Appkey = dayu.Appkey
	alidayu.AppSecret = dayu.AppSecret
	return &sender{
		config:    dayu,
		cache:     cache.New(dayu.ExpiresMinute*time.Minute, dayu.ClearsMinute*time.Minute),
		retryMin:  dayu.RetryMinSecond * time.Second,
		codeChars: []byte(dayu.CodeChars),
	}
}

func (s *sender) Verify(prefix string, userId uint, phone, code string) bool {
	key := prefix + phone

	limitedCode, ok := s.cache.Get(key)
	if !ok {
		return false
	}
	s.cache.Delete(key)

	lcode := limitedCode.(*LimitedCode)
	return lcode.UserID == userId && lcode.Code == code
}

func (s *sender) Send(prefix string, userId uint, phone string) error {
	key := prefix + phone

	if limitedCode, ok := s.cache.Get(key); ok {
		lcode := limitedCode.(*LimitedCode)
		if time.Now().Add(-s.retryMin).Unix() < lcode.GenAt || lcode.UserID != userId {
			return cerr.RetrySmsFailed
		}
	}

	lcode := LimitedCode{
		Code:   uniuri.NewLenChars(s.config.CodeLen, s.codeChars),
		UserID: userId,
		GenAt:  time.Now().Unix(),
	}
	s.cache.Set(key, &lcode, cache.DefaultExpiration)

	//	fmt.Println("phone:", phone, "sent code:", lcode.Code)
	res, err := alidayu.SendOnce(phone, s.config.SignName, s.config.Template, fmt.Sprintf(`{"code":"%s"}`, lcode.Code))
	if err != nil {
		glog.Errorln(err)
		return cerr.SendSmsError
	}

	if !res.Success {
		glog.Errorln(*res.ResultError)
		return cerr.SendSmsFailed
	}

	return nil
}
