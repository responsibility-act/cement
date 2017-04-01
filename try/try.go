package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/empirefox/bongine/captchar"
	"github.com/empirefox/bongine/config"
)

func main() {
	conf, err := config.LoadFromXpsWithEnv()
	if err != nil {
		panic(err)
	}

	captcha, err := captchar.NewCaptchar(&conf.Captcha)
	if err != nil {
		panic(err)
	}

	for _ = range make([]int, 10) {
		data, err := captcha.New(0)
		if err != nil {
			panic(err)
		}

		b := []byte(*data.Base64)
		b = b[1 : len(b)-1]
		d := make([]byte, base64.StdEncoding.DecodedLen(len(b)))
		n, err := base64.StdEncoding.Decode(d, b)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("captcha-%s.png", data.ID), d[:n], os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

}
