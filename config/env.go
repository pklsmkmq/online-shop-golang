package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	SUPABASE_URL string
	SUPABASE_KEY string
	JWT_SECRET   []byte
)

func LoadEnv() {
	godotenv.Load()

	SUPABASE_URL = os.Getenv("SUPABASE_URL")
	SUPABASE_KEY = os.Getenv("SUPABASE_ANON_KEY")
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

	if SUPABASE_URL == "" || SUPABASE_KEY == "" {
		panic("ENV SUPABASE_URL dan SUPABASE_ANON_KEY wajib di set")
	}
}
