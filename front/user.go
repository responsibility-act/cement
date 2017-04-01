package front

import "github.com/dgrijalva/jwt-go"

type WritableUserInfo struct {
	tableName    struct{} `sql:"cc_user,alias:u"`
	ID           uint
	Nickname     string
	HeadImageURL string
	Sex          int
	City         string
	Province     string
	Birthday     int64
	Hobby        string
	Intro        string
	UpdatedAt    int64
}

type ReadonlyUserInfo struct {
	CreatedAt int64
	SigninAt  int64
}

type UserInfo struct {
	Writable *WritableUserInfo
	ReadonlyUserInfo
	HasPayKey bool `sql:"-"`
}

type SetUserInfoResponse struct {
	UpdatedAt int64
}

// out by auth middleware
type UserTokenResponse struct {
	AccessToken  *string
	RefreshToken *string
	User         *UserInfo
}

type TokenClaims struct {
	jwt.StandardClaims
	OpenId string `json:"oid,omitempty"`
	UserId uint   `json:"uid,omitempty"`
	User1  uint   `json:"us1,omitempty"`
	Phone  string `json:"mob,omitempty"`
	Nonce  string `json:"non,omitempty"`
}

type PreBindPhonePayload struct {
	Phone string
}

type BindPhonePayload struct {
	Phone        string `binding:"required"`
	Code         string `binding:"required"`
	CaptchaID    string `binding:"required"`
	Captcha      string `binding:"required"`
	RefreshToken string
}

type RefreshTokenResponse struct {
	OK          bool
	AccessToken *string
}

type SetPaykeyPayload struct {
	Key       string
	Code      string
	CaptchaID string `binding:"required"`
	Captcha   string `binding:"required"`
}
