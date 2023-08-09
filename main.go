package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🏁 Starting Bot")

	godotenv.Load()

	req, err := http.NewRequest(http.MethodGet, os.Getenv("CLOUDFLARE_KV_URL")+"/values/go", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLOUDFLARE_KV_KEY")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	repo := string(body)

	req, err = http.NewRequest(http.MethodGet, os.Getenv("CLOUDFLARE_KV_URL")+"/metadata/go", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("CLOUDFLARE_KV_KEY")))

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	type Metadata struct {
		Result struct {
			AuthorName string `json:"author_name"`
			RepoName   string `json:"repo_name"`
			TotalStars string `json:"total_stars"`
		}
	}
	var metadata Metadata
	json.Unmarshal([]byte(body), &metadata)

	payload := strings.NewReader("status=" + metadata.Result.AuthorName + "/" + metadata.Result.RepoName + " - ⭐ " + metadata.Result.TotalStars + "\n" + repo)
	req, err = http.NewRequest(http.MethodPost, "https://botsin.space/api/v1/statuses", payload)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("MASTODON_ACCESS_TOKEN")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err = http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode == 200 {
		fmt.Println("✅ Successfully posted to Mastodon")
	}
}
