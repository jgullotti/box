package main

import (
	"flag"
	"fmt"

	"github.com/ttacon/box"
	"golang.org/x/oauth2"
)

var (
	clientId     = flag.String("cid", "", "OAuth Client ID")
	clientSecret = flag.String("csec", "", "OAuth Client Secret")

	accessToken  = flag.String("atok", "", "Access Token")
	refreshToken = flag.String("rtok", "", "Refresh Token")

	fileId    = flag.String("fid", "", "File (ID) to grab")
	expiresAt = flag.String("expires-at", "", "what time the lock should expire at")
)

func main() {
	flag.Parse()

	if len(*clientId) == 0 || len(*clientSecret) == 0 ||
		len(*accessToken) == 0 || len(*refreshToken) == 0 ||
		len(*fileId) == 0 {
		fmt.Println("unfortunately all flags must be provided")
		return
	}

	// Set our OAuth2 configuration up
	var (
		configSource = box.NewConfigSource(
			&oauth2.Config{
				ClientID:     *clientId,
				ClientSecret: *clientSecret,
				Scopes:       nil,
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://app.box.com/api/oauth2/authorize",
					TokenURL: "https://app.box.com/api/oauth2/token",
				},
				RedirectURL: "http://localhost:8080/handle",
			},
		)
		tok = &oauth2.Token{
			TokenType:    "Bearer",
			AccessToken:  *accessToken,
			RefreshToken: *refreshToken,
		}
		c = configSource.NewClient(tok)
	)

	var lock *box.Lock
	if len(*expiresAt) != 0 {
		lock = &box.Lock{
			Type:                "lock",
			ExpiresAt:           *expiresAt,
			IsDownloadPrevented: false,
		}
	}

	resp, err := c.FileService().Lock(*fileId, lock)
	fmt.Println("resp: ", resp)
	fmt.Println("err: ", err)

	// Print out the new tokens for next time
	fmt.Printf("\n%#v\n", tok)
}
