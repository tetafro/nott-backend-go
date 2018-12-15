package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/auth"
	"github.com/tetafro/nott-backend-go/internal/domain"
	"github.com/tetafro/nott-backend-go/internal/storage"
)

// AuthController handles HTTP API requests.
type AuthController struct {
	users  storage.UsersRepo
	tokens storage.TokensRepo
	log    logrus.FieldLogger
}

// NewAuthController creates new controller.
func NewAuthController(
	users storage.UsersRepo,
	tokens storage.TokensRepo,
	log logrus.FieldLogger,
) *AuthController {
	return &AuthController{users: users, tokens: tokens, log: log}
}

// Register handles request for registering new users.
func (c *AuthController) Register(w http.ResponseWriter, req *http.Request) {
	var body authRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		badRequest(w, "invalid json")
		return
	}
	defer req.Body.Close()

	_, err := c.users.GetByEmail(body.Email)
	if err == nil {
		badRequest(w, "email is already taken")
		return
	}
	if err != domain.ErrNotFound {
		c.log.Errorf("Failed to check user: %v", err)
		internalServerError(w)
		return
	}

	// Create user in the repository
	u := auth.User{Email: body.Email}
	u.Password, err = auth.HashPassword(body.Password)
	if err != nil {
		c.log.Errorf("Failed to hash password: %v", err)
		internalServerError(w)
		return
	}
	u, err = c.users.Create(u)
	if err != nil {
		c.log.Errorf("Failed to create user: %v", err)
		internalServerError(w)
		return
	}

	// Generate token
	t, err := c.tokens.Create(auth.Token{UserID: u.ID})
	if err != nil {
		c.log.Errorf("Failed to create token: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, t)
}

// Login handles request for logging in using email+password.
func (c *AuthController) Login(w http.ResponseWriter, req *http.Request) {
	var body authRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		badRequest(w, "invalid json")
		return
	}
	defer req.Body.Close()

	user, err := c.users.GetByEmail(body.Email)
	if err == domain.ErrNotFound {
		respond(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		c.log.Errorf("Failed to get user: %v", err)
		internalServerError(w)
		return
	}

	if !auth.CheckPassword(body.Password, user.Password) {
		respond(w, http.StatusBadRequest, "invalid email or password")
		return
	}

	// Generate token
	t, err := c.tokens.Create(auth.Token{UserID: user.ID})
	if err != nil {
		c.log.Errorf("Failed to create token: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, t)
}

// Logout handles request for logging out.
func (c *AuthController) Logout(w http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	if token == "" {
		unauthorized(w)
		return
	}
	if err := c.tokens.Delete(auth.Token{String: token}); err != nil {
		internalServerError(w)
		return
	}
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

// GetProfile handles request for getting current logged in user.
func (c *AuthController) GetProfile(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	respond(w, http.StatusOK, user)
}

// UpdateProfile handles request for getting current logged in user.
func (c *AuthController) UpdateProfile(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)

	var err error

	u := auth.User{}
	if err = json.NewDecoder(req.Body).Decode(&u); err != nil {
		badRequest(w, "invalid json")
		return
	}

	if err = u.Validate(); err != nil {
		badRequest(w, "invalid user: "+err.Error())
		return
	}

	// Get all available fields
	var modified bool
	if user.Email != u.Email {
		user.Email = u.Email
		modified = true
	}

	// Save model if there was any changes
	if modified {
		u, err = c.users.Update(*user)
		if err == domain.ErrNotFound {
			notFound(w)
			return
		}
		if err != nil {
			c.log.Errorf("Failed to update user: %v", err)
			internalServerError(w)
			return
		}
	}

	respond(w, http.StatusOK, u)
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
