package server

import (
	"net/url"
	"strings"

	"github.com/empirefox/bogger"
	"github.com/empirefox/cement/api"
	"github.com/empirefox/cement/captchar"
	"github.com/empirefox/cement/config"
	"github.com/empirefox/cement/dbs"
	"github.com/iris-contrib/middleware/cors"
	"github.com/iris-contrib/middleware/loggerzap"
	"github.com/iris-contrib/middleware/secure"
	"github.com/kataras/iris"
	"github.com/uber-go/zap"
)

type Server struct {
	*iris.Framework
	config config.Config
	logger zap.Logger
	dbs    dbs.DBS

	captcha *captchar.Captchar
	qiniu   *bogger.Qiniu
}

func NewServer(config *config.Config) (*Server, error) {
	dbs, err := dbs.NewDBService(config)
	if err != nil {
		return nil, err
	}

	app := iris.New(config.Iris)

	if config.Server.ZapIris {
		app.Use(loggerzap.New(loggerzap.Config{
			Status: true,
			IP:     true,
			Method: true,
			Path:   true,
		}))
	}

	app.Use(secure.New(config.Iris))

	app.Use(cors.New(cors.Options{
		AllowedHeaders:   []string{"Origin", "Authorization", "Content-Type"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: false,
		AllowedOrigins:   config.Security.CorsOrigins,
		MaxAge:           config.Security.CorsAgeSecond,
		AllowOriginFunc: func(origin string) bool {
			origin = strings.ToLower(origin)
			for _, o := range config.Security.CorsOrigins {
				if o == origin {
					return true
				}
			}

			u, err := url.ParseRequestURI(origin)
			return err == nil && dbs.SiteExist(u.Host)
		},
	}))
	//	app.OnError(iris.StatusBadRequest, func(ctx *iris.Context) {
	//		ctx.Write("CUSTOM 404 NOT FOUND ERROR PAGE")
	//		ctx.Log("http status: 400 happened!")
	//	})
	captcha, err := captchar.NewCaptchar(&config.Captcha)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Framework: app,
		config:    *config,
		logger:    config.Logger,
		dbs:       dbs,

		captcha: captcha,
		qiniu:   bogger.NewQiniu(config.Qiniu),
	}

	apis := api.NewApis()
	s.Get(apis.GetUptoken, s.GetUptoken)
	s.Post(apis.PostList, s.PostList)
	s.Post(apis.PostDelete, s.PostDelete)

	return s
}
