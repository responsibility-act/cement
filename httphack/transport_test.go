package httphack

import (
	"net/http/httputil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHackClient(t *testing.T) {
	require := require.New(t)

	client := NewHackClient("127.0.0.1:9999")
	res, err := client.Get("https://www.wechat.com/dev-transport.go")
	require.Nil(err)
	dump, err := httputil.DumpResponse(res, true)
	require.Nil(err)
	os.Stdout.Write(dump)
}
