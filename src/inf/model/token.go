package model

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Token struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

const (
	TOKEN_EXPIRY    = 30 * 24 * 60 * 60 // 30 Days
	ENV_SIGNING_KEY = "WISHLIST_TOK_SIGNING_KEY"
)

type TokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func Authenticate(email, password string) (*Token, error) {

	// Find User details
	usr, err := GetUserWithEmail(email)
	if err != nil {
		return nil, err
	}

	hashedPassword := HashPassword(password)

	// Check auth details are valid
	if usr.Password != hashedPassword {
		return nil, fmt.Errorf("Email / Password doesn't match")
	}

	// Hash Token
	expiry := time.Now().Unix() + TOKEN_EXPIRY
	tok := generateToken(email, expiry)
	signd, err := tok.SignedString([]byte(os.Getenv(ENV_SIGNING_KEY)))
	if err != nil {
		return nil, InternalServerError("Failed to generate authentication token")
	}

	return &Token{
		Token:     signd,
		ExpiresAt: fmt.Sprint(expiry),
	}, nil
}

// Generates a JWT with the user's email and expiry time
func generateToken(email string, expiresAt int64) *jwt.Token {
	cl := TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Email: email,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, cl)
}

// Middleware to check if a token is valid
var TokenValidator = middleware.JWTWithConfig(middleware.JWTConfig{
	Claims:     &TokenClaims{},
	SigningKey: []byte(os.Getenv(ENV_SIGNING_KEY)),
})

// AuthValidator Middleware makes sure the user making the request is the user being altered
func AuthValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Get currently logged in token's ID
		token, ok := c.Get("user").(*jwt.Token)
		if !ok {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to retrieve JWT data from middleware",
			}
		}
		claims, ok := token.Claims.(*TokenClaims)
		if !ok {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to retrieve user claims from middleware JWT",
			}
		}
		loggedInUser, err := GetUserWithEmail(claims.Email)
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Token used is for a user that doesn't exist anymore",
			}
		}

		idSigned, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
		if err != nil {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("Invalid User ID provided: '%s'", c.Param("user_id")),
			}
		}
		id := uint64(idSigned)

		// if user is not updating themselves, abort
		if id != loggedInUser.ID {
			return &echo.HTTPError{
				Code:    http.StatusForbidden,
				Message: "You are forbidden from changing/deleting other users than yourself",
			}
		}

		return next(c)
	}
}
