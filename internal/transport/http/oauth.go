package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// OAuthController handles HTTP API requests.
type OAuthController struct {
	providers map[string]*auth.OAuthProvider
	users     storage.UsersRepo
	tokens    storage.TokensRepo
	log       logrus.FieldLogger
}

// NewOAuthController creates new OAuth controller.
func NewOAuthController(
	p map[string]*auth.OAuthProvider,
	u storage.UsersRepo,
	t storage.TokensRepo,
	log logrus.FieldLogger,
) *OAuthController {
	return &OAuthController{providers: p, users: u, tokens: t, log: log}
}

// Providers handles request for getting list of currently
// available OAuth2 providers.
func (c *OAuthController) Providers(w http.ResponseWriter, req *http.Request) {
	// Convert map to list
	pp := make([]*auth.OAuthProvider, len(c.providers))
	i := 0
	for _, p := range c.providers {
		pp[i] = p
		i++
	}
	respond(w, http.StatusOK, pp)
}

// Github handles callback requests from Github.
func (c *OAuthController) Github(w http.ResponseWriter, req *http.Request) {
	p, ok := c.providers["github"]
	if !ok {
		badRequest(w, "Github provider is currently disabled")
		return
	}

	var body oauthRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		badRequest(w, "invalid json")
		return
	}
	defer req.Body.Close()

	email, err := p.GetEmail(body.Code)
	if err != nil {
		c.log.Errorf("Failed to get user from OAuth provider: %v", err)
		internalServerError(w)
		return
	}

	token, err := c.handleUser(email)
	if err != nil {
		c.log.Errorf("Failed to handle user: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, token)
}

func (c *OAuthController) handleUser(email string) (auth.Token, error) {
	// Get or create user
	u, err := c.users.GetByEmail(email)
	switch err {
	case nil:
		// Got existing user, proceed
	case domain.ErrNotFound:
		// Create new user
		u, err = c.users.Create(auth.User{Email: email})
		if err != nil {
			return auth.Token{}, fmt.Errorf("create user: %v", err)
		}
	default:
		// Unexpected error
		return auth.Token{}, fmt.Errorf("get user: %v", err)
	}

	// Generate token
	t, err := c.tokens.Create(auth.Token{UserID: u.ID})
	if err != nil {
		return auth.Token{}, fmt.Errorf("create token: %v", err)
	}

	return t, nil
}

type oauthRequest struct {
	Code string `json:"code"`
}
