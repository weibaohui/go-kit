package wechat

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/weibaohui/go-kit/strkit"
	"strconv"
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

// AESResponse 回复aes加密的消息给微信服务器.
//  msg:       经过 encoding/xml.Marshal 得到的结果符合微信消息格式的任何数据结构
//  timestamp: 时间戳, 如果为 0 则默认使用 Context.Timestamp
//  nonce:     随机数, 如果为 "" 则默认使用 Context.Nonce
//  random:    16字节的随机字符串, 如果为 nil 则默认使用 Context.Random
func AESResponse(msg interface{}, timestamp int64, nonce string, random []byte) (msgSignature string, err error) {
	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}
	if nonce == "" {
		nonce = strkit.GetRandomString(5)
	}
	if len(random) == 0 {
		random = []byte(strkit.GetRandomString(10))
	}

	msgPlaintext, err := json.Marshal(msg)
	if err != nil {
		return
	}

	encryptedMsg := AESEncryptMsg(random, msgPlaintext, appId, aesKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)
	timestampString := strconv.FormatInt(timestamp, 10)
	msgSignature = MsgSign(token, timestampString, nonce, base64EncryptedMsg)

	return

}
