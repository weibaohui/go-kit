package wechat

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/weibaohui/go-kit/httpkit"
	"github.com/weibaohui/go-kit/strkit"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	oriId, appId = "gh_970c702c7a97", "wx43de43f93421f607"
	token = "cYl8R3UY8a3i3YPC"
	base64AESKey = "94xfc6XMlm28fIKTcljBoBCWcMy3aUZkQjqGgh2xyEn"

	if len(base64AESKey) != 43 {
		errors.New("the length of base64AESKey must equal to 43")
	}
	aesKey, _ = base64.StdEncoding.DecodeString(base64AESKey + "=")

}
func TestServeHTTP(t *testing.T) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := strkit.GetRandomString(10)
	msg := "helloworld"
	signature := Sign(token, timestamp, nonce)

	random := []byte("cc9632a98304f81c")
	encryptedMsg := AESEncryptMsg(random, []byte(msg), appId, aesKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)
	msgSignature := MsgSign(token, timestamp, nonce, base64EncryptedMsg)
	///?signature=%s&timestamp=%s&nonce=%s&openid=%s&encrypt_type=aes&msg_signature=%s
	query := fmt.Sprintf(""+
		"/"+
		"?signature=%s"+
		"&timestamp=%s"+
		"&nonce=%s"+
		"&encrypt_type=aes"+
		"&msg_signature=%s",
		signature, timestamp, nonce, msgSignature)
	t.Log(query)
	t.Logf("%d", time.Now().Unix())

	server := httptest.NewServer(http.HandlerFunc(ServeHTTP))
	t.Log(server.URL)
	finalUrl := server.URL + query
	t.Logf("final url: %s", finalUrl)

	str, err := httpkit.NewRequest(finalUrl, "POST").EnableDebug().EnableDump().Body(base64EncryptedMsg).
		String()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(str)
	//if str != msg {
	//	t.Fatalf("except %s,got %s", msg, str)
	//}

}
