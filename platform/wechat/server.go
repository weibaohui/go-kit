package wechat

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

//  oriId:        可选; 公众号的原始ID(微信公众号管理后台查看), 如果设置了值则该Server只能处理 ToUserName 为该值的公众号的消息(事件);
//  appId:        可选; 公众号的AppId, 如果设置了值则安全模式时该Server只能处理 AppId 为该值的公众号的消息(事件);
//  token:        必须; 公众号用于验证签名的token;
//  base64AESKey: 可选; aes加密解密key, 43字节长(base64编码, 去掉了尾部的'='), 安全模式必须设置;
//
//var (
//	appID        string
//	token        string
//	base64AESKey string
//	aesKey       []byte
//)

type Server struct {
	AppID        string
	Token        string
	base64AESKey string
	AESKey       []byte
}

//
//func init() {
//	appID = "wx43de43f93421f607"
//	token = "cYl8R3UY8a3i3YPC"
//	base64AESKey = "94xfc6XMlm28fIKTcljBoBCWcMy3aUZkQjqGgh2xyEn"
//
//	if len(base64AESKey) != 43 {
//		errors.New("the length of base64AESKey must equal to 43")
//	}
//	aesKey, _ = base64.StdEncoding.DecodeString(base64AESKey + "=")
//}

func NewServer() *Server {
	return &Server{}
}
func (s *Server) SetAppID(appID string) *Server {
	s.AppID = appID
	return s
}

func (s *Server) SetToken(token string) *Server {
	s.Token = token
	return s
}

func (s *Server) SetBase64AESKey(key string) *Server {
	s.base64AESKey = key

	if len(s.base64AESKey) != 43 {
		errors.New("the length of base64AESKey must equal to 43")
	}
	aesKey, _ := base64.StdEncoding.DecodeString(s.base64AESKey + "=")
	s.AESKey = aesKey
	return s
}

func (s *Server) Start() {
	http.HandleFunc("/", s.ServeHTTP)

	http.ListenAndServe(":9011", nil)
}

// ServeHTTP 例子
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//fmt.Printf("%v\n", r.URL.Query())

	query := r.URL.Query()
	errorHandler := DefaultErrorHandler

	switch r.Method {
	case "POST":
		switch encryptType := query.Get("encrypt_type"); encryptType {
		case "aes":

			_, msg, err := s.MsgAESEncrypt(r)
			if err != nil {
				errorHandler.ServeError(w, r, err)
				return
			}
			//fmt.Println("xxx---start")
			//fmt.Println(appID)
			//fmt.Println(msg)
			//fmt.Println("xxx---end")
			fmt.Fprint(w, msg)

		case "", "raw":
			s, err := s.MsgNoEncrypt(r)
			if err != nil {
				errorHandler.ServeError(w, r, err)
				return
			}
			fmt.Fprint(w, s)
		default:

			errorHandler.ServeError(w, r, errors.New("unknown encrypt_type: "+encryptType))
		}
	case "GET": // 验证回调URL是否有效
		s, err := s.EchoStr(r)
		if err != nil {
			errorHandler.ServeError(w, r, err)
			return
		}
		io.WriteString(w, s)
	}
}

// 验证回调URL是否有效
func (s *Server) EchoStr(r *http.Request) (string, error) {
	query := r.URL.Query()

	haveSignature := query.Get("signature")
	if haveSignature == "" {
		return "", errors.New("not found signature query parameter")
	}
	timestamp := query.Get("timestamp")
	if timestamp == "" {
		return "", errors.New("not found timestamp query parameter")
	}
	nonce := query.Get("nonce")
	if nonce == "" {
		return "", errors.New("not found nonce query parameter")
	}
	echoStr := query.Get("echoStr")
	if echoStr == "" {
		return "", errors.New("not found echoStr query parameter")
	}

	wantSignature := Sign(s.Token, timestamp, nonce)
	if haveSignature != wantSignature {
		err := fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature)
		return "", err

	}
	return echoStr, nil
}

// 明文方式获取原文
func (s *Server) MsgNoEncrypt(r *http.Request) (string, error) {
	query := r.URL.Query()
	haveSignature := query.Get("signature")
	if haveSignature == "" {
		return "", errors.New("not found signature query parameter")
	}
	timestampString := query.Get("timestamp")
	if timestampString == "" {
		return "", errors.New("not found timestamp query parameter")
	}
	_, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		err = fmt.Errorf("can not parse timestamp query parameter %q to int64", timestampString)
		return "", err
	}

	nonce := query.Get("nonce")
	if nonce == "" {
		return "", errors.New("not found nonce query parameter")
	}

	wantSignature := Sign(s.Token, timestampString, nonce)
	if haveSignature != wantSignature {
		err = fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature)
		return "", err

	}

	msgPlaintext, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(msgPlaintext), nil
}

// AES方式解密出原文
// 返回 appID,msg,error
func (s *Server) MsgAESEncrypt(r *http.Request) (string, string, error) {
	query := r.URL.Query()
	haveMsgSignature := query.Get("msg_signature")
	if haveMsgSignature == "" {
		return "", "", errors.New("not found msg_signature query parameter")
	}
	timestampString := query.Get("timestamp")
	if timestampString == "" {
		return "", "", errors.New("not found timestamp query parameter")
	}
	_, err := strconv.ParseInt(timestampString, 10, 64)
	if err != nil {
		err = fmt.Errorf("can not parse timestamp query parameter %q to int64", timestampString)
		return "", "", err
	}
	nonce := query.Get("nonce")
	if nonce == "" {
		return "", "", errors.New("not found nonce query parameter")
	}

	buffer := textBufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer textBufferPool.Put(buffer)

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		return "", "", err
	}
	requestBodyBytes := buffer.Bytes()
	fmt.Errorf("%s\n", string(requestBodyBytes))
	//收到的body就是base64加密的
	//errorHandler.ServeError(w, r, errors.New("Base64EncryptedMsg        "+string(requestBodyBytes)))

	wantMsgSignature := MsgSign(s.Token, timestampString, nonce, string(requestBodyBytes))
	if haveMsgSignature != wantMsgSignature {
		err := fmt.Errorf("check msg_signature failed, have: %s, want: %s", haveMsgSignature, wantMsgSignature)
		return "", "", err
	}

	encryptedMsg := make([]byte, base64.StdEncoding.DecodedLen(len(requestBodyBytes)))
	encryptedMsgLen, err := base64.StdEncoding.Decode(encryptedMsg, []byte(requestBodyBytes))
	if err != nil {
		return "", "", err
	}
	encryptedMsg = encryptedMsg[:encryptedMsgLen]

	_, msgPlaintext, haveAppIdBytes, err := AESDecryptMsg(encryptedMsg, s.AESKey)
	return string(haveAppIdBytes), string(msgPlaintext), nil
}
