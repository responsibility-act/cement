package security

import (
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	mpoauth2 "github.com/chanxuehong/wechat.v2/mp/oauth2"
	"github.com/dchest/uniuri"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/empirefox/cement/cerr"
	"github.com/empirefox/cement/config"
	"github.com/empirefox/cement/front"
	"github.com/empirefox/cement/models"
	"github.com/empirefox/esecend/db-service"
	"github.com/empirefox/reform"
	"github.com/golang/glog"
	"github.com/patrickmn/go-cache"
	"github.com/sony/sonyflake"
)

type Handler struct {
	conf  *config.Security
	sf    *sonyflake.Sonyflake
	cache *cache.Cache
	db    *dbsrv.DbService
}

func NewHandler(config *config.Config, db *dbsrv.DbService) *Handler {
	return &Handler{
		conf:  &config.Security,
		sf:    sonyflake.NewSonyflake(sonyflake.Settings{StartTime: time.Date(2017, time.May, 1, 0, 0, 0, 0, time.UTC)}),
		cache: cache.New(config.Security.ExpiresMinute*time.Minute, config.Security.ClearsMinute*time.Minute),
		db:    db,
	}
}

func (h *Handler) Login(userinfo *mpoauth2.UserInfo, user1 uint) (interface{}, error) {
	var usr models.User
	var err error
	var refreshToken = uniuri.NewLen(32)
	var encRefreshToken []byte

	if encRefreshToken, err = bcrypt.GenerateFromPassword([]byte(refreshToken), 4); err != nil {
		return nil, err
	}

	tx, err := h.db.Tx()
	if err != nil {
		return nil, err
	}
	defer tx.RollbackIfNeeded()

	db := tx.GetDB()
	if err = db.FindOneTo(&usr, "$UnionId", userinfo.UnionId); err == reform.ErrNoRows {
		usr = models.User{
			OpenId:       userinfo.OpenId,
			UnionId:      userinfo.UnionId,
			RefreshToken: &encRefreshToken,
			CreatedAt:    time.Now().Unix(),

			Nickname:     userinfo.Nickname,
			Sex:          userinfo.Sex,
			City:         userinfo.City,
			Province:     userinfo.Province,
			HeadImageURL: userinfo.HeadImageURL, // TODO Save to our cdn
			User1:        user1,
		}
		glog.Errorln(usr, userinfo.Nickname)
		err = db.Insert(&usr)
	} else if err == nil {
		usr.RefreshToken = &encRefreshToken
		usr.SigninAt = time.Now().Unix()
		err = db.UpdateColumns(&usr, "RefreshToken", "SigninAt")
	}
	if err == nil {
		err = tx.Commit()
	}
	if err != nil {
		return nil, err
	}

	tok, err := h.NewToken(&usr)
	if err != nil {
		return nil, err
	}

	return &front.UserTokenResponse{
		AccessToken:  tok,
		RefreshToken: &refreshToken,
		User:         usr.Info(),
	}, nil
}

// NewToken generate token string json
func (h *Handler) NewToken(usr *models.User) (*string, error) {
	return h.NewTokenWithIat(usr, time.Now().Unix())
}

func (h *Handler) NewTokenWithIat(usr *models.User, now int64) (*string, error) {
	sonyid, err := h.sf.NextID()
	if err != nil {
		return nil, cerr.SonyFlakeTimeout
	}

	claims := &front.TokenClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.FormatUint(sonyid, 36),
			ExpiresAt: now + h.conf.TokenLife*60,
			IssuedAt:  now,
			Subject:   "Weixin",
		},
		OpenId: usr.OpenId,
		UserId: usr.ID,
		User1:  usr.User1,
		Phone:  usr.Phone,
		Nonce:  uniuri.NewLen(32),
	}

	key := []byte(uniuri.NewLen(128))
	h.cache.Set(claims.Id, key, cache.DefaultExpiration)

	token := jwt.NewWithClaims(jwt.GetSigningMethod(h.conf.SignAlg), claims)
	tok, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return &tok, nil
}

func (h *Handler) RefreshToken(tok *jwt.Token, refreshToken []byte) (token *string, err error) {
	claims := tok.Claims.(*front.TokenClaims)
	if claims.ExpiresAt-time.Now().Unix() > h.conf.RefreshIn*60 {
		return nil, cerr.NoNeedRefreshToken
	}
	if len(refreshToken) == 0 {
		return nil, cerr.NoRefreshToken
	}

	var usr models.User
	if err = h.db.GetDB().FindByPrimaryKeyTo(&usr, claims.UserId); err != nil {
		return nil, err
	}
	if usr.RefreshToken == nil || len(*usr.RefreshToken) == 0 {
		return nil, cerr.NoRefreshToken
	}
	if err = h.CompareRefreshToken(*usr.RefreshToken, refreshToken); err != nil {
		return nil, cerr.InvalidRefreshToken
	}

	token, err = h.NewToken(&usr)
	if err != nil {
		return nil, err
	}
	h.cache.Delete(claims.Id)
	return token, nil
}

func (h *Handler) RevokeToken(tok *jwt.Token) error {
	claims := tok.Claims.(*front.TokenClaims)
	h.cache.Delete(claims.Id)
	if claims.UserId == 0 {
		return nil
	}
	return h.db.GetDB().UpdateColumns(&models.User{ID: claims.UserId}, "RefreshToken")
}

func (h *Handler) FindKeyfunc(tok *jwt.Token) (interface{}, error) {
	if tok.Method.Alg() != h.conf.SignAlg {
		return nil, cerr.InvalidSignAlg
	}

	claims := tok.Claims.(*front.TokenClaims)
	key, ok := h.cache.Get(claims.Id)
	if !ok {
		return nil, cerr.InvalidClaimId
	}
	return key, nil

}

func (h *Handler) ParseToken(req *http.Request) (tok *jwt.Token, tokUsr interface{}, err error) {
	tok, err = request.ParseFromRequestWithClaims(req, request.OAuth2Extractor, &front.TokenClaims{}, h.FindKeyfunc)
	if err != nil {
		return tok, nil, err
	}

	claims := tok.Claims.(*front.TokenClaims)
	usr := &models.User{
		OpenId: claims.OpenId,
		ID:     claims.UserId,
		User1:  claims.User1,
		Phone:  claims.Phone,
	}

	return tok, usr, err
}

func (h *Handler) CompareRefreshToken(hashed, input []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, input)
}
