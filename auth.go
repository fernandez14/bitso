package bitso

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/dghubble/sling"
	"io/ioutil"
	"strconv"
	"time"
)

func bitsoAuth(s *sling.Sling, secret []byte) *sling.Sling {
	t := time.Now().UnixNano()
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(strconv.Itoa(int(t))))
	r, err := s.Request()
	if err != nil {
		return s
	}
	h.Write([]byte(r.Method))
	h.Write([]byte(r.URL.Path))
	if r.Method == "POST" {
		b, err := ioutil.ReadAll(r.Body)
		if err == nil {
			h.Write(b)
		}
	}
	sha := hex.EncodeToString(h.Sum(nil))
	return s.Set("authorization", "Bitso "+sha)
}
