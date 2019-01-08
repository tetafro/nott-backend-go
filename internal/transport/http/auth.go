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
	users   storage.UsersRepo
	tokener auth.Tokener
	log     logrus.FieldLogger
}

// NewAuthController creates new controller.
func NewAuthController(u storage.UsersRepo, t auth.Tokener, log logrus.FieldLogger) *AuthController {
	return &AuthController{users: u, tokener: t, log: log}
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
	user := auth.User{Email: body.Email}
	user.Password, err = auth.HashPassword(body.Password)
	if err != nil {
		c.log.Errorf("Failed to hash password: %v", err)
		internalServerError(w)
		return
	}
	user, err = c.users.Create(user)
	if err != nil {
		c.log.Errorf("Failed to create user: %v", err)
		internalServerError(w)
		return
	}

	// Generate token
	t, err := c.tokener.Issue(user)
	if err != nil {
		c.log.Errorf("Failed to issue token: %v", err)
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
	t, err := c.tokener.Issue(user)
	if err != nil {
		c.log.Errorf("Failed to issue token: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, t)
}

// GetProfile handles request for getting current logged in user.
func (c *AuthController) GetProfile(w http.ResponseWriter, req *http.Request) {
	userID := getUserID(req)

	user, err := c.users.GetByID(userID)
	if err == domain.ErrNotFound {
		notFound(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to get user: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, user)
}

// UpdateProfile handles request for getting current logged in user.
func (c *AuthController) UpdateProfile(w http.ResponseWriter, req *http.Request) {
	userID := getUserID(req)

	var err error

	user := auth.User{}
	if err = json.NewDecoder(req.Body).Decode(&user); err != nil {
		badRequest(w, "invalid json")
		return
	}
	user.ID = userID

	if err = user.Validate(); err != nil {
		badRequest(w, "invalid user: "+err.Error())
		return
	}

	user, err = c.users.Update(user)
	if err == domain.ErrNotFound {
		notFound(w)
		return
	}
	if err != nil {
		c.log.Errorf("Failed to update user: %v", err)
		internalServerError(w)
		return
	}

	respond(w, http.StatusOK, user)
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
