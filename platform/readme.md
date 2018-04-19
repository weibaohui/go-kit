timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := strkit.GetRandomString(10)

	msg := `{
    "pendingCode": "000000008",
    "lastUpdateDate": "20180418150828",
    "pendingStatus": "1",
    "pendingNote":"na164"
}`
	s := wechat.Server{}
	s.SetAppID("bpm").
		SetToken("cYl8R3UY8a3i3YPC").
		SetBase64AESKey("94xfc6XMlm28fIKTcljBoBCWcMy3aUZkQjqGgh2xyEn")

	random := []byte("123456")

	encryptedMsg := wechat.AESEncryptMsg(random, []byte(msg), s.AppID, s.AESKey)
	base64EncryptedMsg := base64.StdEncoding.EncodeToString(encryptedMsg)
	msgSignature := wechat.MsgSign(s.Token, timestamp, nonce, base64EncryptedMsg)
	// ?signature=%s&timestamp=%s&nonce=%s&openid=%s&encrypt_type=aes&msg_signature=%s
	query := fmt.Sprintf(""+
		"?appID=%s"+
		"&timestamp=%s"+
		"&nonce=%s"+
		"&encrypt_type=aes"+
		"&msg_signature=%s",
		s.AppID, timestamp, nonce, msgSignature)

	finalUrl := "http://127.0.0.1:8077/endToDo" + query

	str, err := httpkit.NewRequest(finalUrl, "POST").EnableDebug().Body(base64EncryptedMsg).String()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(str)