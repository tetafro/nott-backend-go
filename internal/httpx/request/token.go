package request

import (
	"net/http"
	"strings"
)

// GetToken gets token from HTTP header.
// Header format (RFC2617):
// Authorization: Token token="abcd1234"
func GetToken(req *http.Request) string {
	auth, ok := req.Header["Authorization"]
	if !ok || len(auth) == 0 {
		return ""
	}

	token := auth[0]
	if !strings.HasPrefix(token, `Token token="`) || !strings.HasSuffix(token, `"`) {
		return ""
	}

	return token[13 : len(token)-1]
}
