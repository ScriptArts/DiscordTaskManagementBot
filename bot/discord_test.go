package bot

import (
	"os"
	"testing"
)

func TestGetDiscordClientInvalidToken(t *testing.T) {
	os.Setenv("DISCORD_BOT_TOKEN", "asjkldhflkjasdh")

	discord, err := GetDiscordClient()
	if err != nil {
		t.Fatal("GetDiscordClient()", err)
	}

	err = discord.Open()
	if err == nil {
		t.Fatal("不正なトークンでAPI認証成功")
		discord.Close() // 認証に成功してしまっているので明示的にClose
	}
}

// 存在するトークンを用いたテストは開発者の環境で行う
//func TestGetDiscordClientTrustToken(t *testing.T) {
//	sampleToken := os.Getenv("DISCORD_BOT_TEST_TOKEN")
//	os.Setenv("DISCORD_BOT_TOKEN", sampleToken)
//
//	discord, err := GetDiscordClient()
//	if err != nil {
//		t.Fatal("GetDiscordClient()", err)
//	}
//
//	err = discord.Open()
//	if err != nil {
//		t.Fatal("存在するトークンでAPI認証失敗", err)
//	} else {
//		discord.Close() // 認証に成功しているので明示的にClose
//	}
//}
