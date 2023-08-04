package middleware

import (
	"bootcamp_course_microservice/infras"
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

type Authentication struct {
	DB     *infras.Conn
	Secret []byte
}

func ProvideAuthentication(db *infras.Conn, secret []byte) *Authentication {
	return &Authentication{
		DB:     db,
		Secret: secret,
	}
}

func IsTeacherRole(tokenString string) (bool, error) {
	// Make an HTTP GET request to the bootcamp-auth microservice to validate the JWT
	authServiceURL := "http://localhost:8080/v1/validate-auth"
	req, err := http.NewRequest("GET", authServiceURL, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+tokenString)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var responseData map[string]string
		err := json.NewDecoder(resp.Body).Decode(&responseData)
		if err != nil {
			return false, err
		}

		role, ok := responseData["role"]
		if !ok {
			return false, nil
		}

		return role == "teacher", nil
	}

	return false, nil
}

func (a *Authentication) VerifyJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		isTeacher, err := IsTeacherRole(tokenString)
		if err != nil {
			http.Error(w, "Error validating role", http.StatusInternalServerError)
			return
		}

		if !isTeacher {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// User is authorized, add the token to the request context for use in protected endpoints
		ctx := context.WithValue(r.Context(), "user", tokenString)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
