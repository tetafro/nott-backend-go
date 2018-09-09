package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/errors"
	"github.com/tetafro/nott-backend-go/internal/http/request"
	"github.com/tetafro/nott-backend-go/internal/http/response"
)

// Controller handles HTTP API requests.
type Controller struct {
	users  UsersRepo
	tokens TokensRepo
	log    logrus.FieldLogger
}

// NewController creates new controller.
func NewController(users UsersRepo, tokens TokensRepo, log logrus.FieldLogger) *Controller {
	return &Controller{users: users, tokens: tokens, log: log}
}

// Login handles request for logging in using email+password.
func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError(fmt.Errorf("Invalid JSON")).
			Write(w)
		return
	}
	defer req.Body.Close()

	user, err := c.users.GetByEmail(body.Email)
	if err == ErrNotFound {
		response.New().
			WithStatus(http.StatusUnauthorized).
			WithError(fmt.Errorf("Invalid email or password")).
			Write(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to get user: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	if !checkPassword(body.Password, user.Password) {
		response.New().
			WithStatus(http.StatusUnauthorized).
			WithError(fmt.Errorf("Invalid email or password")).
			Write(w)
		return
	}

	t, err := c.tokens.Create(Token{UserID: user.ID})
	if err != nil {
		c.log.Errorf("Failed to create token: %v", err)
		response.InternalServerError().Write(w)
		return
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(t).
		Write(w)
}

// Logout handles request for logging out.
func (c *Controller) Logout(w http.ResponseWriter, req *http.Request) {
	token := request.GetToken(req)
	if token == "" {
		response.Unauthorized().Write(w)
		return
	}
	if err := c.tokens.Delete(Token{String: token}); err != nil {
		response.InternalServerError().Write(w)
		return
	}
	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

// GetProfile handles request for getting current logged in user.
func (c *Controller) GetProfile(w http.ResponseWriter, req *http.Request) {
	user := GetUser(req)
	response.New().
		WithStatus(http.StatusOK).
		WithData(user).
		Write(w)
}

// UpdateProfile handles request for getting current logged in user.
func (c *Controller) UpdateProfile(w http.ResponseWriter, req *http.Request) {
	user := GetUser(req)

	var err error

	u := User{}
	if err = json.NewDecoder(req.Body).Decode(&u); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError(fmt.Errorf("Invalid JSON")).
			Write(w)
		return
	}

	if err = u.Validate(); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError(fmt.Errorf("Invalid user: %v", err)).
			Write(w)
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
		if err == errors.ErrNotFound {
			response.NotFound().Write(w)
			return
		}
		if err != nil {
			c.log.Errorf("Failed to update user: %v", err)
			response.InternalServerError().Write(w)
			return
		}
	}

	response.New().
		WithStatus(http.StatusOK).
		WithData(u).
		Write(w)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
