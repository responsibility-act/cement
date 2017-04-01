package config

import (
	"time"

	"github.com/empirefox/bogger"
	"github.com/iris-contrib/middleware/secure"
	"github.com/kataras/iris"
	"github.com/uber-go/zap"
)

type Server struct {
	Host      string
	Port      int    `env:"PORT" default:"443"`
	Cert      []byte `json:"-" xps:"server.crt"`
	Key       []byte `json:"-" xps:"server.key"`
	DevMode   bool   `env:"DEV_MODE"` // remove if not used
	TLS       bool   `env:"TLS" validate:"rq=Cert,rq=Key"`
	ZapLevel  string `env:"ZAP_LEVEL" default:"error" validate:"zap_level"`
	ZapIris   bool   `env:"ZAP_IRIS"`
	FakeToken bool   `env:"FAKE_TOKEN"`
}

type Security struct {
	SignAlg       string        `default:"HS256" validate:"sign_alg"`
	TokenLife     int64         `default:"60"` // 60 minute
	RefreshIn     int64         `default:"5"`  // last 5 minute
	ExpiresMinute time.Duration `default:"61ns"`
	ClearsMinute  time.Duration `default:"10ns"`
	SecendOrigin  string        `validate:"required,url"` // without '/'
	CorsOrigins   []string      `validate:"required,gt=0,dive,required,url"`
	CorsAgeSecond int
}

type Captcha struct {
	StrType       int           `validate:"gte=0,lte=3"` // NUM LOWER UPPER ALL
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

type Order struct {
	EvalTimeoutDay        uint          `default:"15"`
	CompleteTimeoutDay    uint          `default:"10"`
	HistoryTimeoutDay     uint          `default:"5"`
	CheckoutExpiresMinute time.Duration `default:"30ns"`
	WxPayExpiresMinute    time.Duration `default:"120ns"`
	MaintainTimeMinute    uint          `default:"60"`
}

type Money struct {
	StoreSaleFeePercent uint
	User1RebatePercent  uint
	Store1RebatePercent uint
	WithdrawDesc        string
}

type Weixin struct {
	WebScope       string `default:"snsapi_base" validate:"eq=snsapi_base|eq=snsapi_userinfo"`
	AppId          string `validate:"required"`
	ApiKey         string `validate:"required"`
	MchKey         string `validate:"required"`
	MchId          string `validate:"required"`
	Cert           []byte `json:"-" validate:"gt=0" xps:"wxapi.crt"`
	Key            []byte `json:"-" validate:"gt=0" xps:"wxapi.key"`
	PayBody        string `validate:"required"`
	PayNotifyURL   string `validate:"required,uri"`
	TransCheckName string
}

type Alidayu struct {
	Appkey         string `validate:"required"`
	AppSecret      string `validate:"required"`
	CodeChars      string
	CodeLen        int
	SignName       string        `validate:"required"`
	Template       string        `validate:"required"`
	RetryMinSecond time.Duration `default:"50ns"`
	ExpiresMinute  time.Duration `default:"2ns"`
	ClearsMinute   time.Duration `default:"4ns"`
}

// pg.v5 Options
type Postgres struct {
	// Network type, either tcp or unix.
	// Default is tcp.
	Network string
	// TCP host:port or Unix socket depending on Network.
	Addr string

	User     string `validate:"required"`
	Password string
	Database string `validate:"required"`

	// Maximum number of retries before giving up.
	// Default is to not retry failed queries.
	MaxRetries int
	// Whether to retry queries cancelled because of statement_timeout.
	RetryStatementTimeout bool

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	WriteTimeout time.Duration

	// Maximum number of socket connections.
	// Default is 20 connections.
	PoolSize int
	// Time for which client waits for free connection if all
	// connections are busy before returning an error.
	// Default is 5 seconds.
	PoolTimeout time.Duration
	// Time after which client closes idle connections.
	// Default is to not close idle connections.
	IdleTimeout time.Duration
	// Connection age at which client retires (closes) the connection.
	// Primarily useful with proxies like HAProxy.
	// Default is to not close aged connections.
	MaxAge time.Duration
	// Frequency of idle checks.
	// Default is 1 minute.
	IdleCheckFrequency time.Duration

	// When true Tx does not issue BEGIN, COMMIT, or ROLLBACK.
	// Also underlying database connection is immediately returned to the pool.
	// This is primarily useful for running your database tests in one big
	// transaction, because PostgreSQL does not support nested transactions.
	DisableTransaction bool
}

type Paging struct {
	PageSize uint64 `default:"25"`
	MaxSize  uint64 `default:"100"`
}

type QiniuPrefix struct {
	Common string `default:"c/"`
	User   string `default:"u/"`
	Site   string `default:"s/"`

	Assets  string `default:"a/"`
	Head    string `default:"h/"`
	Product string `default:"p/"`
	Taks    string `default:"t/"`
	Other   string `default:"o/"`
}

type Config struct {
	Env         Env `json:"-"`
	Server      Server
	Iris        iris.Configuration `validate:"-"`
	IrisSecure  secure.Options     `validate:"-"`
	Security    Security
	Captcha     Captcha
	Order       Order
	Money       Money
	Weixin      Weixin
	Alidayu     Alidayu
	Postgres    Postgres
	Paging      Paging
	QiniuPrefix QiniuPrefix
	Qiniu       bogger.Config

	Logger zap.Logger `json:"-"`
}
