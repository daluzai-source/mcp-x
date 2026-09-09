package posts

import (
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/require"
)

func resultText(t *testing.T, res *mcpsdk.CallToolResult) string {
	t.Helper()

	require.NotNil(t, res)
	require.Len(t, res.Content, 1)

	content, ok := res.Content[0].(*mcpsdk.TextContent)
	require.True(t, ok, "content is not a text block")

	return content.Text
}
