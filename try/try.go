package main

import (
	"encoding/json"
	"os"

	"github.com/empirefox/cement/xogen"
)

func main() {
	user := &xogen.User{}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	e.Encode(user)
}
