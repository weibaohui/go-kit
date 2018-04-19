package main

import "github.com/weibaohui/go-kit/platform/wechat"

func main() {
	w := wechat.Server{}
	w.SetAppID("wx43de43f93421f607").
		SetToken("cYl8R3UY8a3i3YPC").
		SetBase64AESKey("94xfc6XMlm28fIKTcljBoBCWcMy3aUZkQjqGgh2xyEn")
	w.Start()

}
