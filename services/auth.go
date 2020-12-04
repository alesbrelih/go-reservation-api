package services

import (
	"time"

	"github.com/alesbrelih/go-reservation-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var InvalidTokenError = errors.New("Invalid jwt token")

func NewAuthHandler(secret string, accessExp, refreshExp time.Duration) AuthService {
	return &authService{
		secret:     []byte(secret),
		accessExp:  accessExp,
		refreshExp: refreshExp,
	}
}

type AuthService interface {
	GenerateJwtPair(id string) (*models.TokenPair, error)
	RefreshToken(refreshToken string) (*models.TokenPair, error)
	GetClaims(token string) (*jwt.StandardClaims, error)
}

type authService struct {
	secret     []byte
	accessExp  time.Duration
	refreshExp time.Duration
}

func (a *authService) GenerateJwtPair(id string) (*models.TokenPair, error) {
	accessClaims := &jwt.StandardClaims{
		Subject:   id,
		ExpiresAt: time.Now().Add(a.accessExp).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := access.SignedString(a.secret)
	if err != nil {
		return nil, errors.Wrap(err, "Error signing access token")
	}

	refreshClaims := &jwt.StandardClaims{
		Subject:   id,
		ExpiresAt: time.Now().Add(a.refreshExp).Unix(),
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refresh.SignedString(a.secret)
	if err != nil {
		return nil, errors.Wrap(err, "Error signing refresh token")
	}

	tokenPair := &models.TokenPair{
		Refresh: refreshToken,
		Access:  accessToken,
	}
	return tokenPair, nil
}

func (a *authService) RefreshToken(refreshToken string) (*models.TokenPair, error) {

	claims, err := a.GetClaims(refreshToken)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting claims in refresh token")
	}
	return a.GenerateJwtPair(claims.Subject)
}

func (a *authService) GetClaims(jwtString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(jwtString, &jwt.StandardClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Wrap(InvalidTokenError, "Missmatch signing method")
		}
		return a.secret, nil
	})

	if err != nil {
		return nil, errors.Wrap(InvalidTokenError, "Invalid token")
	}

	if !token.Valid {
		return nil, errors.Wrap(InvalidTokenError, "Token is not valid")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, errors.Wrap(InvalidTokenError, "Claims cant be parsed")
	}
	return claims, nil
}
