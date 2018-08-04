package auth

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/tetafro/nott-backend-go/internal/httpx/request"
	"github.com/tetafro/nott-backend-go/internal/httpx/response"
)

// Controller handles HTTP API requests.
type Controller struct {
	users  UsersRepo
	tokens TokensRepo
	log    *logrus.Logger
}

// NewController creates new controller.
func NewController(users UsersRepo, tokens TokensRepo, log *logrus.Logger) *Controller {
	return &Controller{users: users, tokens: tokens, log: log}
}

// Login handles request for logging in using email+password.
func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
	var body loginRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		response.New().
			WithStatus(http.StatusBadRequest).
			WithError("Invalid JSON").
			Write(w)
		return
	}
	defer req.Body.Close()

	user, err := c.users.GetByEmail(body.Email)
	if err == ErrNotFound {
		response.New().
			WithStatus(http.StatusUnauthorized).
			WithError(response.Error("Invalid email or password")).
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
			WithError(response.Error("Invalid email or password")).
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

// Profile handles request for getting current logged in user.
func (c *Controller) Profile(w http.ResponseWriter, req *http.Request) {
	user := GetUser(req)
	response.New().
		WithStatus(http.StatusOK).
		WithData(user).
		Write(w)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
