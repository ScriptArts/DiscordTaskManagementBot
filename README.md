# DiscordTaskManagementBot

個人にメンションされたタスクの管理ボット

- env

`.env.template` をもとに、`.env`ファイルを作成する。

```
# envファイル存在するディレクトリへのPath
DISCORD_TASK_MANAGEMENT_DIR=
# DiscordBotのトークン
DISCORD_BOT_TOKEN=
# データベースタイプ (sqlite3, pg, mysql, mssql)
DISCORD_TASK_MANAGEMENT_DATABASE_TYPE=sqlite3
# データベースの接続文字列
DISCORD_TASK_MANAGEMENT_DATABASE_CONNECTION_STR=app.db
# クリエイター追加コマンドなどの管理者コマンド実行権限（DiscordのユーザID）
# ex) 11111111,22222222
DISCORD_BOT_OP_LIST=
```

`go run server.go` or (`go build server.go` and run binary file)