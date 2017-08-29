package captchar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/color"
	"image/png"
	"time"

	"go.uber.org/zap"

	"github.com/afocus/captcha"
	"github.com/dchest/uniuri"
	"github.com/empirefox/cement/clog"
	"github.com/empirefox/cement/perr"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/patrickmn/go-cache"
)

type Config struct {
	StrType       int           `validate:"gte=0,lte=3"` // NUM=0 LOWER UPPER ALL
	Distur        int           `validate:"gt=0,lte=16"                     default:"8"`
	FontPaths     []string      `validate:"required,gt=0,dive,required,uri`
	BgColorLen    int           `validate:"gte=5"                           default:"5"`
	FrontColorLen int           `validate:"gte=5"                           default:"5"`
	Width         int           `validate:"gte=48"                          default:"92"`
	Height        int           `validate:"gte=20"                          default:"32"`
	CodeLen       int           `validate:"gte=4,lte=10"                    default:"4"`
	ExpiresSecond time.Duration `validate:"gte=30"                          default:"60ns"`
	ClearsSecond  time.Duration `validate:"gte=30,gtefield=ExpiresSecond"   default:"120ns"`
	NbfInSecond   int64         `validate:"gte=3,ltfield=ExpiresSecond"     default:"3"`
}

type Captchar interface {
	New(userId uint) (*Captcha, error)
	Verify(userId uint, key, value string) bool
}

type Captcha struct {
	ID     string
	Base64 *json.RawMessage
}

type Cached struct {
	UserID    uint
	Value     string
	CreatedAt int64
}

type captchar struct {
	config   Config
	cc       *captcha.Captcha
	capCache *cache.Cache
	log      *zap.Logger
}

func NewCaptchar(config Config, logger clog.Logger) (Captchar, error) {
	lm := logger.Module("captcha")
	cc := captcha.New()

	if err := cc.SetFont(config.FontPaths...); err != nil {
		lm.Error("FontPaths init", zap.Error(err))
		return nil, perr.BadFontPaths
	}

	var bgColors []color.Color
	for _, c := range colorful.FastWarmPalette(config.BgColorLen) {
		bgColors = append(bgColors, c)
	}
	var frontColors []color.Color
	for _, c := range colorful.FastHappyPalette(config.FrontColorLen) {
		frontColors = append(frontColors, c)
	}

	cc.SetSize(config.Width, config.Height)
	cc.SetDisturbance(captcha.DisturLevel(config.Distur))
	cc.SetBkgColor(bgColors...)
	cc.SetFrontColor(frontColors...)

	return &captchar{
		config:   config,
		cc:       cc,
		capCache: cache.New(config.ExpiresSecond*time.Second, config.ClearsSecond*time.Second),
	}, nil
}

func (c *captchar) New(userId uint) (*Captcha, error) {
	img, value := c.cc.Create(c.config.CodeLen, captcha.StrType(c.config.StrType))

	var b bytes.Buffer
	b.WriteByte('"')
	b64 := base64.NewEncoder(base64.StdEncoding, &b)
	defer b64.Close()
	if err := png.Encode(b64, img); err != nil {
		c.log.Error("png failed", zap.Error(err))
		return nil, perr.TryAgain
	}
	b.WriteByte('"')

	key := uniuri.New()
	cached := Cached{
		UserID:    userId,
		Value:     value,
		CreatedAt: time.Now().Unix(),
	}
	for c.capCache.Add(key, &cached, cache.DefaultExpiration) != nil {
		key = uniuri.NewLen(20)
	}

	data := json.RawMessage(b.Bytes())
	return &Captcha{
		ID:     key,
		Base64: &data,
	}, nil
}

func (c *captchar) Verify(userId uint, key, value string) bool {
	fact, ok := c.capCache.Get(key)
	if !ok {
		return false
	}

	c.capCache.Delete(key)
	cached := fact.(*Cached)
	return userId == cached.UserID &&
		time.Now().Unix() > cached.CreatedAt+c.config.NbfInSecond &&
		bytes.Compare(bytes.ToLower([]byte(value)), bytes.ToLower([]byte(cached.Value))) == 0
}
