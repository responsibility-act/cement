package wo2

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/dgrijalva/jwt-go"
	"github.com/empirefox/cement/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func init() {
	//	log.Level = logrus.DebugLevel
}

func TestAuther_Middleware(t *testing.T) {
	auther := newAuther(newSecHandler())
	var called bool
	request(auther, "", func(c *gin.Context) { called = true })
	require.True(t, called, "Should be called without MustAuthed")
	called = false

	request(auther, "", auther.MustAuthed, func(c *gin.Context) { called = true })
	require.False(t, called, "Should not be called with MustAuthed")

	request(auther, "TOKEN", auther.MustAuthed, func(c *gin.Context) { called = true })
	require.True(t, called, "Should be called with token to MustAuthed")
	called = false
}

func request(auther *Auther, bearer string, h ...gin.HandlerFunc) *httptest.ResponseRecorder {
	r := gin.Default()
	r.Use(auther.Middleware())
	r.GET("/request", h...)
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/request", nil)
	if bearer != "" {
		req.Header.Set("Authorization", "BEARER "+bearer)
	}
	r.ServeHTTP(res, req)
	return res
}

func TestAuther_OauthSuccess(t *testing.T) {
	auther := newAuther(newSecHandler())
	res := requestOauth(auther, "POST", strings.NewReader(`{"code":"CODE"}`))
	require.Equal(t, 200, res.Code)
	require.Equal(t, `{"token":"TOKEN"}`, strings.TrimSpace(res.Body.String()))
}

func TestAuther_OauthFail(t *testing.T) {
	auther := newAuther(newSecHandler())

	res := requestOauth(auther, "GET", strings.NewReader(`{"code":"CODE"}`))
	require.Equal(t, 404, res.Code)

	res = requestOauth(auther, "POST", strings.NewReader(`{"code":""}`))
	require.Equal(t, 401, res.Code)

	res = requestOauth(auther, "POST", strings.NewReader(`{"code2":"aaa"}`))
	require.Equal(t, 401, res.Code)

	sh1 := &secHandler{
		loginReturn: gin.H{"token": "TOKEN"},
		loginErr:    fmt.Errorf("Failed to login"),
		parsedToken: new(jwt.Token),
		parsedUser:  gin.H{"ID": 111},
		parsedErr:   nil,
	}
	auther = newAuther(sh1)
	res = requestOauth(auther, "POST", strings.NewReader(`{"code":"CODE"}`))
	require.Equal(t, 401, res.Code)
}

func requestOauth(auther *Auther, method string, payload io.Reader) *httptest.ResponseRecorder {
	r := gin.Default()
	r.Use(auther.Middleware())
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(method, "/auth", payload)
	r.ServeHTTP(res, req)
	return res
}

func newSecHandler() *secHandler {
	return &secHandler{
		loginReturn: gin.H{"token": "TOKEN"},
		loginErr:    nil,
		parsedToken: &jwt.Token{Valid: true},
		parsedUser:  gin.H{"ID": 111},
		parsedErr:   fmt.Errorf("cannot parse token"),
	}
}

type secHandler struct {
	loginReturn interface{}
	loginErr    error
	parsedToken *jwt.Token
	parsedUser  interface{}
	parsedErr   error
}

func (h *secHandler) Login(userinfo *mpoauth2.UserInfo) (ret interface{}, err error) {
	return h.loginReturn, h.loginErr
}

func (h *secHandler) ParseToken(c *gin.Context) (tok *jwt.Token, tokUsr interface{}, err error) {
	if c.Request.Header.Get("Authorization") == "BEARER TOKEN" {
		return h.parsedToken, h.parsedUser, nil
	}
	return nil, nil, h.parsedErr
}

type oauth2HttpClientTransport struct{}

func (t oauth2HttpClientTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")
	responseBody := ""
	if req.URL.Path == "/sns/oauth2/access_token" {
		responseBody = `
{
  "access_token":"ACCESS_TOKEN",
  "expires_in":7200,
  "refresh_token":"REFRESH_TOKEN",
  "openid":"OPENID",
  "scope":"SCOPE",
  "unionid":"o6_bmasdasdsad6_2sgVt7hMZOPfL"
}`
	} else {
		responseBody = `
{
  "access_token":"ACCESS_TOKEN",
  "expires_in":7200,
  "refresh_token":"REFRESH_TOKEN",
  "openid":"OPENID",
  "scope":"SCOPE"
}`
	}
	response.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

type getUserInfoHttpClientTransport struct{}

func (t getUserInfoHttpClientTransport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")
	responseBody := `
{
  "openid":"OPENID",
  "nickname":"NICKNAME",
  "sex":1,
  "province":"PROVINCE",
  "city":"CITY",
  "country":"COUNTRY",
  "headimgurl":"http://wx.qlogo.cn/mmopen/g3MonUZtNHkdmzicIlibx6iaFqAc56vxLSUfpb6n5WKSYVY0ChQKkiaJSgQ1dZuTOgvLLrhJbERQQ4eMsv84eavHiaiceqxibJxCfHe/46", 
	"privilege":[
	"PRIVILEGE1",
	"PRIVILEGE2"
   ],
   "unionid":"o6_bmasdasdsad6_2sgVt7hMZOPfL"
}`
	response.Body = ioutil.NopCloser(strings.NewReader(responseBody))
	return response, nil
}

func newAuther(h *secHandler) *Auther {
	oauth2HttpClient := &http.Client{Transport: new(oauth2HttpClientTransport)}
	getUserInfoHttpClient := &http.Client{Transport: new(getUserInfoHttpClientTransport)}
	return &Auther{

		Oauth2HttpClient:      oauth2HttpClient,
		GetUserInfoHttpClient: getUserInfoHttpClient,

		wx:          new(config.Weixin),
		wxOauthPath: "/auth",
		secHandler:  h,
	}
}
