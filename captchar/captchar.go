package captchar

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/color"
	"image/png"
	"time"

	"github.com/afocus/captcha"
	"github.com/dchest/uniuri"
	"github.com/empirefox/bongine/config"
	"github.com/empirefox/bongine/front"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/patrickmn/go-cache"
)

type Captchar interface {
	New(userId uint) (*front.Captcha, error)
	Verify(userId uint, key, value string) bool
}

type Cached struct {
	UserID    uint
	Value     string
	CreatedAt int64
}

type captchar struct {
	config   *config.Captcha
	cc       *captcha.Captcha
	capCache *cache.Cache
}

func NewCaptchar(config *config.Captcha) (Captchar, error) {
	cc := captcha.New()

	if err := cc.SetFont(config.FontPaths...); err != nil {
		return nil, err
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

func (c *captchar) New(userId uint) (*front.Captcha, error) {
	img, value := c.cc.Create(c.config.CodeLen, captcha.StrType(c.config.StrType))

	var b bytes.Buffer
	b.WriteByte('"')
	b64 := base64.NewEncoder(base64.StdEncoding, &b)
	defer b64.Close()
	if err := png.Encode(b64, img); err != nil {
		return nil, err
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
	return &front.Captcha{
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
