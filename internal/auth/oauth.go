package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// OAuthProvider is a list of currently available OAuth providers.
type OAuthProvider struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	userInfoURL string
	config      oauth2.Config
}

// NewGithubProvider initializes GitHub OAuth provider.
func NewGithubProvider(host, id, secret string) *OAuthProvider {
	host = strings.TrimRight(host, "/")
	return &OAuthProvider{
		Name:        "GitHub",
		URL:         github.Endpoint.AuthURL + "?scope=user:email&client_id=" + id,
		userInfoURL: "https://api.github.com/user",
		config: oauth2.Config{
			ClientID:     id,
			ClientSecret: secret,
			RedirectURL:  host + "/login-github",
			Scopes:       []string{},
			Endpoint:     github.Endpoint,
		},
	}
}

// GetEmail gets user from Github using provided code, and returns his email.
func (p *OAuthProvider) GetEmail(code string) (string, error) {
	// Get access token
	t, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return "", fmt.Errorf("get access token: %v", err)
	}

	client := p.config.Client(context.Background(), t)
	resp, err := client.Get(p.userInfoURL)
	if err != nil {
		return "", fmt.Errorf("get user info: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %v", err)
	}

	var gu githubUser
	if err := json.Unmarshal(data, &gu); err != nil {
		return "", fmt.Errorf("unmarshal user: %v", err)
	}

	return gu.Email, nil
}

type githubUser struct {
	Email string `json:"email"`
}
