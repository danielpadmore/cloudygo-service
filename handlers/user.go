package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/danielpadmore/cloudygo-service/data"
	"github.com/danielpadmore/cloudygo-service/logs"
	"github.com/dgrijalva/jwt-go"
)

// JWTSecret is an arbitrary secret
const JWTSecret = "dennis_the_menace"

// User contains database connection data
type User struct {
	logger     logs.Logger
	connection data.Connection
}

// AuthData describes the shape of inbound authentication details
type AuthData struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// AuthResponse describes outbound authenticated data
type AuthResponse struct {
	UserID   string `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Token    string `json:"token,omitempty"`
}

// NewUser creates a new user
func NewUser(logger logs.Logger, connection data.Connection) *User {
	return &User{logger, connection}
}

// No handling for root of this route
func (user *User) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	user.logger.Debug(newLog("Attempted request at invalid user endpoint: %s", r.URL.String()))
	http.NotFound(rw, r)
}

// Register is a handler to register a new user
func (user *User) Register(rw http.ResponseWriter, r *http.Request) {
	user.logger.Info(newLog("Register request made at %s", r.URL.String()))
	body := AuthData{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		user.logger.Info(newLog("Unable to parse register body: %s", err.Error()))
		http.Error(rw, "Unable to parse request", http.StatusBadRequest)
		return
	}

	u, err := user.connection.CreateUser(body.Username, body.Password)
	if err != nil {
		user.logger.Info(newLog("Unable to register user %s: %s", body.Username, err.Error()))
		http.Error(rw, fmt.Sprintf("Unable to register user %s", body.Username), http.StatusBadRequest)
		return
	}

	tokenString, err := generateJWTToken(u.ID, u.Username)
	if err != nil {
		user.logger.Error(newLog("Unable to generate JWT token"))
		http.Error(rw, "Unable to generate JWT token", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(AuthResponse{
		UserID:   u.ID,
		Username: u.Username,
		Token:    tokenString,
	})

}

// SignIn authenticates and returns a new JWT token with user details
func (user *User) SignIn(rw http.ResponseWriter, r *http.Request) {
	user.logger.Info(newLog("Sign in request made at %s", r.URL.String()))

	body := AuthData{}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		user.logger.Info(newLog("Unable to parse sign in body: %s", err.Error()))
		http.Error(rw, "Unable to parse request", http.StatusBadRequest)
		return
	}

	u, err := user.connection.AuthenticateUser(body.Username, body.Password)
	if err != nil {
		user.logger.Info(newLog("Unable to sign in user %s: %s", body.Username, err.Error()))
		http.Error(rw, "Invalid creds, try again!", http.StatusUnauthorized)
		return
	}

	tokenString, err := generateJWTToken(u.ID, u.Username)
	if err != nil {
		user.logger.Info(newLog("Unable to generate JWT token: %s", err.Error()))
		http.Error(rw, "Unable to generate JWT token", http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(AuthResponse{
		UserID:   u.ID,
		Username: u.Username,
		Token:    tokenString,
	})
}

func generateJWTToken(id string, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  id,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(JWTSecret))
}
