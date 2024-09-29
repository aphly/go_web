package test

import (
	"fmt"
	"go_web/app/core/crypt"
	"testing"
)

func TestS(t *testing.T) {
	key := []byte("x1cs456789abcdef0123456789ascdef")
	en, err := crypt.EncryptAES256ECB("xxxx hoowsl isiw world! womeshi haopengyou =_+ss%", key)
	if err != nil {
		return
	}
	fmt.Println(en)
	de, err := crypt.DecryptAES256ECB(en, key)
	if err != nil {
		return
	}
	fmt.Println(de)
}
