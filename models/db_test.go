package models

import (
	"os"
	"testing"
)

func TestGetDatabaseEmptyEnv(t *testing.T) {
	// 環境変数初期化
	os.Setenv("DISCORD_TASK_MANAGEMENT_DATABASE_TYPE", "")
	os.Setenv("DISCORD_TASK_MANAGEMENT_DATABASE_CONNECTION_STR", "")

	_, err := GetDatabase()
	if err == nil {
		t.Fatal("コネクション文字列を設定してないが、コネクションを開けている")
	}
}

func TestGetDatabase(t *testing.T) {
	os.Setenv("DISCORD_TASK_MANAGEMENT_DATABASE_TYPE", "sqlite3")
	os.Setenv("DISCORD_TASK_MANAGEMENT_DATABASE_CONNECTION_STR", "test.db")

	// 1度目
	db, err := GetDatabase()
	if err != nil {
		t.Fatal(err)
	}

	// 2度目はキャッシュされているコネクションから取得
	db, _ = GetDatabase()
	if err := db.DB().Ping(); err != nil {
		t.Fatal("コネクションがきれていないのに、エラーが返ってきている", err.Error())
	}

	// 明示的にクローズ
	db.Close()

	if err := db.DB().Ping(); err == nil {
		t.Fatal("コネクションを閉じたが、エラーが返ってきていない")
	}

	// コネクションが閉じていた場合、再生成をするのでその処理のテスト
	db, _ = GetDatabase()
	if err := db.DB().Ping(); err != nil {
		t.Fatal("コネクションがきれていないのに、エラーが返ってきている", err.Error())
	}
}
