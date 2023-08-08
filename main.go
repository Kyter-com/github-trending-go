package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🏁 Starting Bot")

	// Load our .env file
	godotenv.Load()

	req, err := http.NewRequest(http.MethodGet, os.Getenv("CLOUDFLARE_KV_URL")+"/values/go", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLOUDFLARE_KV_KEY")))
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println("res status code", res.StatusCode)
}
