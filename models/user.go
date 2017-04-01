package models

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/empirefox/bongine/front"
)

var (
	ComparePaykey = bcrypt.CompareHashAndPassword
	EncPaykey     = bcrypt.GenerateFromPassword
)

type User struct {
	front.WritableUserInfo `pg:",override"`
	front.ReadonlyUserInfo

	// claims below 3 fields with UserInfo.ID
	OpenId string
	Phone  string `sql:"unique"`
	User1  uint

	UnionId string `sql:"unique"`

	// for jwt, auto generated when
	// login   sign with new key
	// logout  remove exist keys
	// refresh set old key life with 1min, add the old jti to new head if still valid
	//         sign with new key
	// jwt is saved in mem K-V(jti:key) cache, not in user table
	// Key string

	// RefreshToken is not lookup every time
	// Only query when need refresh
	// Remove when logout
	RefreshToken *[]byte // bcrypt, no expires

	Paykey *[]byte // for pay, user set, bcrypt
}

func (u *User) Info() *front.UserInfo {
	return &front.UserInfo{
		Writable:         &u.WritableUserInfo,
		ReadonlyUserInfo: u.ReadonlyUserInfo,
		HasPayKey:        u.Paykey != nil && len(*u.Paykey) > 0,
	}
}
