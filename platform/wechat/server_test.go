package wechat

import (
	"encoding/base64"
	"fmt"
	"github.com/weibaohui/go-kit/httpkit"
	"github.com/weibaohui/go-kit/strkit"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func server() *WechatServer {
	w := WechatServer{}
	w.SetAppID("wx43de43f93421f607").
		SetToken("cYl8R3UY8a3i3YPC").
		SetBase64AESKey("94xfc6XMlm28fIKTcljBoBCWcMy3aUZkQjqGgh2xyEn")
	return &w
}
func TestServeHTTP(t *testing.T) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := strkit.GetRandomString(10)
	msg := "helloworld"

	s := server()

	random := []byte("123456")
	encryptedMsg := AESEncryptMsg(random, []byte(msg), s.appID, s.aesKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)
	msgSignature := MsgSign(s.token, timestamp, nonce, base64EncryptedMsg)
	///?signature=%s&timestamp=%s&nonce=%s&openid=%s&encrypt_type=aes&msg_signature=%s
	query := fmt.Sprintf(""+
		"/"+
		"?timestamp=%s"+
		"&nonce=%s"+
		"&encrypt_type=aes"+
		"&msg_signature=%s",
		timestamp, nonce, msgSignature)
	t.Log(query)
	t.Logf("%d", time.Now().Unix())

	server := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))
	t.Log(server.URL)
	finalUrl := server.URL + query
	t.Logf("final url: %s", finalUrl)

	str, err := httpkit.NewRequest(finalUrl, "POST").Body(base64EncryptedMsg).
		String()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(str)
	if str != msg {
		t.Fatalf("except %s,got %s", msg, str)
	}

}

func BenchmarkWechatServer_ServeHTTP(b *testing.B) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := strkit.GetRandomString(10)
	msg := "helloworld"

	s := server()

	random := []byte("123456")
	encryptedMsg := AESEncryptMsg(random, []byte(msg), s.appID, s.aesKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)
	msgSignature := MsgSign(s.token, timestamp, nonce, base64EncryptedMsg)
	///?signature=%s&timestamp=%s&nonce=%s&openid=%s&encrypt_type=aes&msg_signature=%s
	query := fmt.Sprintf(""+
		"/"+
		"?timestamp=%s"+
		"&nonce=%s"+
		"&encrypt_type=aes"+
		"&msg_signature=%s",
		timestamp, nonce, msgSignature)

	server := httptest.NewServer(http.HandlerFunc(s.ServeHTTP))

	finalUrl := server.URL + query

	for i := 0; i < b.N; i++ {

		str, err := httpkit.NewRequest(finalUrl, "POST").Body(base64EncryptedMsg).
			String()
		if err != nil {
			b.Fatal(err.Error())
		}

		if str != msg {
			b.Fatalf("except %s,got %s", msg, str)
		}
	}
}

func TestServeHTTPReal(t *testing.T) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := strkit.GetRandomString(10)
	msg := "helloworld"

	s := server()

	random := []byte("123456")
	encryptedMsg := AESEncryptMsg(random, []byte(msg), s.appID, s.aesKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)
	msgSignature := MsgSign(s.token, timestamp, nonce, base64EncryptedMsg)
	///?signature=%s&timestamp=%s&nonce=%s&openid=%s&encrypt_type=aes&msg_signature=%s
	query := fmt.Sprintf(""+

		"?appID=%s"+
		"&timestamp=%s"+
		"&nonce=%s"+
		"&encrypt_type=aes"+
		"&msg_signature=%s",
		s.appID, timestamp, nonce, msgSignature)
	t.Log(query)

	finalUrl := "http://127.0.0.1:8089/v1/long" + query
	t.Logf("final url: %s", finalUrl)

	str, err := httpkit.NewRequest(finalUrl, "POST").Body(base64EncryptedMsg).
		String()
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(str)
	except := `{"appID":"wx43de43f93421f607","msg":"helloworld","x":"x"}`
	if str != except {
		t.Fatalf("except %s,got %s", msg, str)
	}

}
