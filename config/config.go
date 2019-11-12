package config

import (
	"time"

	"github.com/empirefox/cement/captchar"
	"github.com/empirefox/cement/clog"
	"github.com/empirefox/cement/qiniu"
	"github.com/empirefox/cement/sms"
	"github.com/empirefox/cement/sms/dayu"
	"github.com/empirefox/cement/wepay"

	"github.com/iris-contrib/middleware/secure"
	"github.com/kataras/iris/v12"
)

type Server struct {
	Host      string
	Port      int    `env:"PORT" default:"443"`
	Cert      []byte `json:"-" xps:"server.crt"`
	Key       []byte `json:"-" xps:"server.key"`
	Dev       bool   `env:"DEV"` // this is global dev setting
	TLS       bool   `env:"TLS" validate:"dep=Cert,dep=Key"`
	FakeToken bool   `env:"FAKE_TOKEN"`
}

type Security struct {
	SignAlg       string        `default:"HS256" validate:"eq=HS256|eq=HS384|eq=HS512"`
	TokenLife     int64         `default:"60"` // 60 minute
	RefreshIn     int64         `default:"5"`  // last 5 minute
	ExpiresMinute time.Duration `default:"61ns"`
	ClearsMinute  time.Duration `default:"10ns"`
	SecendOrigin  string        `validate:"required,url"` // without '/'
	CorsOrigins   []string      `validate:"required,gt=0,dive,required,url"`
	CorsAgeSecond int
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
	Server      Server
	Clog        clog.Config
	Iris        iris.Configuration `validate:"-"`
	IrisSecure  secure.Options     `validate:"-"`
	Security    Security
	Order       Order
	Money       Money
	Captcha     captchar.Config
	Wepay       wepay.Config
	Sms         sms.Config
	Alidayu     dayu.Config
	Postgres    Postgres
	Paging      Paging
	Qiniu       qiniu.Config
	QiniuPrefix QiniuPrefix
}

func (c *Config) GetEnvPtrs() []interface{} {
	return []interface{}{&c.Server, &c.Clog}
}
