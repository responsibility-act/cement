package sms

import (
	"fmt"
	"time"

	"github.com/dchest/uniuri"
	"github.com/empirefox/cement/perr"
	"github.com/patrickmn/go-cache"
)

type Config struct {
	Dev            bool
	CodeChars      string
	CodeLen        int
	RetryMinSecond time.Duration `default:"50ns"`
	ExpiresMinute  time.Duration `default:"2ns"`
	ClearsMinute   time.Duration `default:"4ns"`
}

type LimitedCode struct {
	ID     string
	Code   string
	UserID uint
	GenAt  int64
}

type sender struct {
	config    Config
	vendor    Vendor
	cache     *cache.Cache
	retryDur  time.Duration
	codeChars []byte
}

func NewSender(config Config, vendor Vendor) Sender {
	return &sender{
		config:    config,
		vendor:    vendor,
		cache:     cache.New(config.ExpiresMinute*time.Minute, config.ClearsMinute*time.Minute),
		retryDur:  config.RetryMinSecond * time.Second,
		codeChars: []byte(config.CodeChars),
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
		if time.Now().Add(-s.retryDur).Unix() < lcode.GenAt || lcode.UserID != userId {
			return perr.SmsNeedCold
		}
	}

	lcode := LimitedCode{
		Code:   uniuri.NewLenChars(s.config.CodeLen, s.codeChars),
		UserID: userId,
		GenAt:  time.Now().Unix(),
	}
	s.cache.Set(key, &lcode, cache.DefaultExpiration)

	if s.config.Dev {
		fmt.Println(">>> phone:", phone, "code:", lcode.Code)
		return nil
	}
	return s.vendor.Send(phone, lcode.Code)
}
