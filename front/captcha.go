package front

import (
	"encoding/json"
)

type Captcha struct {
	ID     string
	Base64 *json.RawMessage
}
