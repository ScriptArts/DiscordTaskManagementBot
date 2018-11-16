package utils

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func LoadEnv() error {
	dir := filepath.Join(os.Getenv("DISCORD_TASK_MANAGEMENT_DIR"), ".env")
	return godotenv.Load(dir)
}
