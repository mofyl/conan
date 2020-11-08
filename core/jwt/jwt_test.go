package jwt

import (
	"testing"
	"conan/log"
)

func TestEncode(t *testing.T) {

	j := NewJwt()

	tokenStr, err := j.Encode(TokenInfo{
		ChannelName: "aaaaa",
		ExternInfo:  "",
	})
	if err != nil {
		log.Info("encode err is  %s", err.Error())
		return
	}

	log.Info("encode Token is %s", tokenStr)
	err = j.decode(tokenStr)

	if err != nil {
		log.Info("decode Err is %s ", err.Error())
		return
	}
}
