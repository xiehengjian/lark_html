package parse

import (
	"context"
	"testing"

	"github.com/chyroc/lark"
	"github.com/stretchr/testify/assert"
)

func TestParseHtmlToLarkPostMessage(t *testing.T) {
	ctx := context.Background()
	cli := lark.New(lark.WithAppCredential("<APP_ID>", "<APP_SECRET>"))
	htmlStr := `<a href="http://url">Link text</a>`
	data, err := ParseHtmlToLarkPostMessage(htmlStr, cli)
	assert.Nil(t, err)
	_, _, err = cli.Message.Send().ToOpenID("<OPEN_ID>").SendPost(ctx, data)
	assert.Nil(t, err)
}
