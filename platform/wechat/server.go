package wechat

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"unicode"
)

//  oriId:        可选; 公众号的原始ID(微信公众号管理后台查看), 如果设置了值则该Server只能处理 ToUserName 为该值的公众号的消息(事件);
//  appId:        可选; 公众号的AppId, 如果设置了值则安全模式时该Server只能处理 AppId 为该值的公众号的消息(事件);
//  token:        必须; 公众号用于验证签名的token;
//  base64AESKey: 可选; aes加密解密key, 43字节长(base64编码, 去掉了尾部的'='), 安全模式必须设置;

var (
	oriId        string
	appId        string
	token        string
	base64AESKey string
	aesKey       []byte
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

func Start() {
	http.HandleFunc("/", ServeHTTP)

	http.ListenAndServe(":9011", nil)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {

	fmt.Printf("%v\n", r.URL.Query())

	query := r.URL.Query()
	errorHandler := DefaultErrorHandler

	switch r.Method {
	case "POST":
		switch encryptType := query.Get("encrypt_type"); encryptType {
		case "aes":
			haveSignature := query.Get("signature")
			if haveSignature == "" {
				errorHandler.ServeError(w, r, errors.New("not found signature query parameter"))
				return
			}
			haveMsgSignature := query.Get("msg_signature")
			if haveMsgSignature == "" {
				errorHandler.ServeError(w, r, errors.New("not found msg_signature query parameter"))
				return
			}
			timestampString := query.Get("timestamp")
			if timestampString == "" {
				errorHandler.ServeError(w, r, errors.New("not found timestamp query parameter"))
				return
			}
			timestamp, err := strconv.ParseInt(timestampString, 10, 64)
			if err != nil {
				err = fmt.Errorf("can not parse timestamp query parameter %q to int64", timestampString)
				errorHandler.ServeError(w, r, err)
				return
			}
			fmt.Println("timestamp", timestamp)
			nonce := query.Get("nonce")
			if nonce == "" {
				errorHandler.ServeError(w, r, errors.New("not found nonce query parameter"))
				return
			}

			wantSignature := Sign(token, timestampString, nonce)
			if haveSignature != wantSignature {
				err = fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature)
				errorHandler.ServeError(w, r, err)
				return
			}

			buffer := textBufferPool.Get().(*bytes.Buffer)
			buffer.Reset()
			defer textBufferPool.Put(buffer)

			if _, err = buffer.ReadFrom(r.Body); err != nil {
				errorHandler.ServeError(w, r, err)
				return
			}
			requestBodyBytes := buffer.Bytes()

			var requestHttpBody cipherRequestHttpBody
			if err = xmlUnmarshal(requestBodyBytes, &requestHttpBody); err != nil {
				errorHandler.ServeError(w, r, err)
				return
			}

			wantMsgSignature := MsgSign(token, timestampString, nonce, string(requestHttpBody.Base64EncryptedMsg))
			if haveMsgSignature != wantMsgSignature {
				err = fmt.Errorf("check msg_signature failed, have: %s, want: %s", haveMsgSignature, wantMsgSignature)
				errorHandler.ServeError(w, r, err)
				return
			}

			encryptedMsg := make([]byte, base64.StdEncoding.DecodedLen(len(requestHttpBody.Base64EncryptedMsg)))
			encryptedMsgLen, err := base64.StdEncoding.Decode(encryptedMsg, requestHttpBody.Base64EncryptedMsg)
			if err != nil {
				errorHandler.ServeError(w, r, err)
				return
			}
			encryptedMsg = encryptedMsg[:encryptedMsgLen]

			random, msgPlaintext, haveAppIdBytes, err := AESDecryptMsg(encryptedMsg, aesKey)
			fmt.Println("xxx---start")
			fmt.Println(string(random))
			fmt.Println(string(msgPlaintext))
			fmt.Println(string(haveAppIdBytes))
			fmt.Println(err)
			fmt.Println("xxx---end")

		case "", "raw":
			haveSignature := query.Get("signature")
			if haveSignature == "" {
				errorHandler.ServeError(w, r, errors.New("not found signature query parameter"))
				return
			}
			timestampString := query.Get("timestamp")
			if timestampString == "" {
				errorHandler.ServeError(w, r, errors.New("not found timestamp query parameter"))
				return
			}
			timestamp, err := strconv.ParseInt(timestampString, 10, 64)
			if err != nil {
				err = fmt.Errorf("can not parse timestamp query parameter %q to int64", timestampString)
				errorHandler.ServeError(w, r, err)
				return
			}
			fmt.Println("timestamp", timestamp)
			nonce := query.Get("nonce")
			if nonce == "" {
				errorHandler.ServeError(w, r, errors.New("not found nonce query parameter"))
				return
			}

			wantSignature := Sign(token, timestampString, nonce)
			if haveSignature != wantSignature {
				err = fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature)
				errorHandler.ServeError(w, r, err)
				return

			}

			msgPlaintext, err := ioutil.ReadAll(r.Body)
			if err != nil {
				errorHandler.ServeError(w, r, err)
				return
			}
			fmt.Printf("原始内容：%v", string(msgPlaintext))
		default:

			errorHandler.ServeError(w, r, errors.New("unknown encrypt_type: "+encryptType))
		}
	case "GET": // 验证回调URL是否有效
		haveSignature := query.Get("signature")
		if haveSignature == "" {
			errorHandler.ServeError(w, r, errors.New("not found signature query parameter"))
			return
		}
		timestamp := query.Get("timestamp")
		if timestamp == "" {
			errorHandler.ServeError(w, r, errors.New("not found timestamp query parameter"))
			return
		}
		nonce := query.Get("nonce")
		if nonce == "" {
			errorHandler.ServeError(w, r, errors.New("not found nonce query parameter"))
			return
		}
		echostr := query.Get("echostr")
		if echostr == "" {
			errorHandler.ServeError(w, r, errors.New("not found echostr query parameter"))
			return
		}

		wantSignature := Sign(token, timestamp, nonce)
		if haveSignature != wantSignature {
			err := fmt.Errorf("check signature failed, have: %s, want: %s", haveSignature, wantSignature)
			errorHandler.ServeError(w, r, err)
			return

		}
		io.WriteString(w, echostr)
	}
}

// =====================================================================================================================

type cipherRequestHttpBody struct {
	XMLName            struct{} `xml:"xml"`
	ToUserName         string   `xml:"ToUserName"`
	Base64EncryptedMsg []byte   `xml:"Encrypt"`
}

var (
	msgStartElementLiteral = []byte("<xml>")
	msgEndElementLiteral   = []byte("</xml>")

	msgToUserNameStartElementLiteral = []byte("<ToUserName>")
	msgToUserNameEndElementLiteral   = []byte("</ToUserName>")

	msgEncryptStartElementLiteral = []byte("<Encrypt>")
	msgEncryptEndElementLiteral   = []byte("</Encrypt>")

	cdataStartLiteral = []byte("<![CDATA[")
	cdataEndLiteral   = []byte("]]>")
)

//<xml>
//    <ToUserName><![CDATA[gh_b1eb3f8bd6c6]]></ToUserName>
//    <Encrypt><![CDATA[DlCGq+lWQuyjNNK+vDaO0zUltpdUW3u4V00WCzsdNzmZGEhrU7TPxG52viOKCWYPwTMbCzgbCtakZHyNxr5hjoZJ7ORAUYoIAGQy/LDWtAnYgDO+ppKLp0rDq+67Dv3yt+vatMQTh99NII6x9SEGpY3O2h8RpG99+NYevQiOLVKqiQYzan21sX/jE4Y3wZaeudsb4QVjqzRAPaCJ5nS3T31uIR9fjSRgHTDRDOzjQ1cHchge+t6faUhniN5VQVTE+wIYtmnejc55BmHYPfBnTkYah9+cTYnI3diUPJRRiyVocJyHlb+XOZN22dsx9yzKHBAyagaoDIV8Yyb/PahcUbsqGv5wziOgLJQIa6z93/VY7d2Kq2C2oBS+Qb+FI9jLhgc3RvCi+Yno2X3cWoqbsRwoovYdyg6jme/H7nMZn77PSxOGRt/dYiWx2NuBAF7fNFigmbRiive3DyOumNCMvA==]]></Encrypt>
//</xml>
func xmlUnmarshal(data []byte, p *cipherRequestHttpBody) error {
	data = bytes.TrimSpace(data)
	if !bytes.HasPrefix(data, msgStartElementLiteral) || !bytes.HasSuffix(data, msgEndElementLiteral) {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}
	data2 := data[len(msgStartElementLiteral) : len(data)-len(msgEndElementLiteral)]

	// ToUserName
	ToUserNameElementBytes := data2
	i := bytes.Index(ToUserNameElementBytes, msgToUserNameStartElementLiteral)
	if i == -1 {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}
	ToUserNameElementBytes = ToUserNameElementBytes[i+len(msgToUserNameStartElementLiteral):]
	ToUserNameElementBytes = bytes.TrimLeftFunc(ToUserNameElementBytes, unicode.IsSpace)
	if !bytes.HasPrefix(ToUserNameElementBytes, cdataStartLiteral) {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}
	ToUserNameElementBytes = ToUserNameElementBytes[len(cdataStartLiteral):]
	i = bytes.Index(ToUserNameElementBytes, cdataEndLiteral)
	if i == -1 {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}
	ToUserName := ToUserNameElementBytes[:i]
	ToUserNameElementBytes = ToUserNameElementBytes[i+len(cdataEndLiteral):]
	ToUserNameElementBytes = bytes.TrimLeftFunc(ToUserNameElementBytes, unicode.IsSpace)
	if !bytes.HasPrefix(ToUserNameElementBytes, msgToUserNameEndElementLiteral) {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}
	ToUserNameElementBytes = ToUserNameElementBytes[len(msgToUserNameEndElementLiteral):]

	// Encrypt
	EncryptElementBytes := ToUserNameElementBytes
	i = bytes.Index(EncryptElementBytes, msgEncryptStartElementLiteral)
	if i == -1 {
		EncryptElementBytes = data2
		i = bytes.Index(EncryptElementBytes, msgEncryptStartElementLiteral)
		if i == -1 {
			log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
			return xml.Unmarshal(data, p)
		}
	}
	EncryptElementBytes = EncryptElementBytes[i+len(msgEncryptStartElementLiteral):]
	EncryptElementBytes = bytes.TrimLeftFunc(EncryptElementBytes, unicode.IsSpace)
	if !bytes.HasPrefix(EncryptElementBytes, cdataStartLiteral) {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}
	EncryptElementBytes = EncryptElementBytes[len(cdataStartLiteral):]
	i = bytes.Index(EncryptElementBytes, cdataEndLiteral)
	if i == -1 {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}
	Encrypt := EncryptElementBytes[:i]
	EncryptElementBytes = EncryptElementBytes[i+len(cdataEndLiteral):]
	EncryptElementBytes = bytes.TrimLeftFunc(EncryptElementBytes, unicode.IsSpace)
	if !bytes.HasPrefix(EncryptElementBytes, msgEncryptEndElementLiteral) {
		log.Printf("[WARNING] xmlUnmarshal failed, data:\n%s\n", data)
		return xml.Unmarshal(data, p)
	}

	p.ToUserName = string(ToUserName)
	p.Base64EncryptedMsg = Encrypt
	return nil
}
