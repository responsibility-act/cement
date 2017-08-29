package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/empirefox/cement/clog"
	"github.com/empirefox/cement/qiniu"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	logger, err := clog.NewLogger(clog.Config{Dev: true})
	if err != nil {
		panic(err)
	}
	qn := qiniu.NewQiniu(qiniu.Config{
		Ak:              "ZGhoovpdfp8qNzeQmPFMjW0TWfDLkuJ47szA3pdD",
		Sk:              "oFpChuMomUSdOdxPQnGlyxH0nyWSuxC0GXoKcwD1",
		Bucket:          "dogger",
		UpLifeMinute:    1,
		MaxUpLifeMinute: 10,
		UpHost:          "http://upload.qiniu.com",
		UpHostSecure:    "https://up.qbox.me",
	}, logger)
	s := NewServer(qn, logger)
	err = s.Start(":9999")
	if err != nil {
		panic(err)
	}
}

type Server struct {
	*echo.Echo
	qn  *qiniu.Qiniu
	log *zap.Logger
}

func NewServer(qn *qiniu.Qiniu, logger clog.Logger) *Server {
	e := echo.New()
	e.Debug = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	s := &Server{
		Echo: e,
		qn:   qn,
		log:  logger.Module("wepay"),
	}

	s.GET("/qiniu/headtoken/:life", s.GetQiniuHeadToken)
	s.GET("/qiniu/uptoken/:key/:life", s.GetQiniuUptoken)
	s.POST("/qiniu/:prefix", s.PostQiniuList)
	s.DELETE("/qiniu/:key", s.DeleteQiniu)
	s.POST("/echo", s.PostEcho)

	return s
}

func (s *Server) GetQiniuHeadToken(ctx echo.Context) error {
	userId := 100
	life, _ := strconv.ParseUint(ctx.Param("life"), 10, 32)
	secure := s.secure(ctx)
	return ctx.JSON(http.StatusOK, echo.Map{
		"Uptoken": s.qn.Uptoken(fmt.Sprintf("h/%d", userId), uint32(life), secure),
	})
}

func (s *Server) GetQiniuUptoken(ctx echo.Context) error {
	key, err := base64.URLEncoding.DecodeString(ctx.Param("key"))
	if err != nil {
		return ErrBadParamKey
	}
	life, _ := strconv.ParseUint(ctx.Param("life"), 10, 32)
	secure := s.secure(ctx)
	return ctx.JSON(http.StatusOK, echo.Map{
		"Uptoken": s.qn.Uptoken(string(key), uint32(life), secure),
	})
}

func (s *Server) PostQiniuList(ctx echo.Context) error {
	prefix, err := base64.URLEncoding.DecodeString(ctx.Param("prefix"))
	if err != nil {
		return ErrBadParamPrefix
	}

	items, err := s.qn.List(string(prefix))
	if err != nil {
		return err
	}

	s.log.Debug("list ok", zap.ByteString("prefix", prefix))
	return ctx.JSON(http.StatusOK, items)

}

func (s *Server) DeleteQiniu(ctx echo.Context) error {
	key, err := base64.URLEncoding.DecodeString(ctx.Param("key"))
	if err != nil {
		return ErrBadParamKey
	}

	err = s.qn.Delete(string(key))
	if err != nil {
		return err
	}

	s.log.Debug("delete ok", zap.ByteString("key", key))
	return ctx.NoContent(http.StatusOK)
}

func (s *Server) PostEcho(ctx echo.Context) error {
	body := ctx.Request().Body
	return ctx.Stream(http.StatusOK, "application/octet-stream", body)
}

func (s *Server) secure(ctx echo.Context) bool {
	return strings.HasPrefix(ctx.Request().Header.Get("Origin"), "https://")
}
