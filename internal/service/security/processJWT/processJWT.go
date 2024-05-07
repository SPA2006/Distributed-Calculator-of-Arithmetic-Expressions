package processjwt

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

var (
	mySigningKey = []byte("grpc-calc...tsssss")
)

// done
func GenerateNewJWTToken(userID int64) (string, error) {
	now := time.Now()
	convertedUserID := strconv.Itoa(int(userID))

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["user_id"] = convertedUserID
	claims["exp"] = now.Add(time.Minute * 5).Unix()
	claims["nbf"] = now.Unix()
	claims["iat"] = now.Unix()

	tokenString, err := token.SignedString([]byte(mySigningKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// done
func VerifyJWTToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		userID, ok := claims["user_id"].(string)

		if !ok {
			return "", fmt.Errorf("unavailable get user with that id")
		}

		return userID, nil
	}

	return "", fmt.Errorf("invalid token")
}

func GetJWTTokenFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", errors.New("unable to find authorization token")
	}
	tokenString := cookie.Value

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("ivalid token format")
		}

		user, ok := claims["name"].(string)
		if !ok {
			return nil, errors.New("invalid 'name' format in token")
		}

		return []byte(user), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return "", errors.New("invalid token format or invalid token")
	}

	user := token.Claims.(jwt.MapClaims)["name"].(string)
	return user, nil
}
