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

func (a *Authentication) VerifyTeacherJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Make an HTTP GET request to the bootcamp-auth microservice to validate the JWT
		authServiceURL := "http://localhost:8080/v1/validate-auth"
		req, err := http.NewRequest("GET", authServiceURL, nil)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		req.Header.Set("Authorization", "Bearer "+tokenString)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error sending request", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Check the response status code to see if the token is valid
		if resp.StatusCode == http.StatusOK {
			// Token is valid, extract the user's role from the response body
			var responseData map[string]string
			err := json.NewDecoder(resp.Body).Decode(&responseData)
			if err != nil {
				http.Error(w, "Failed to parse response", http.StatusInternalServerError)
				return
			}

			role, ok := responseData["role"]

			if !ok {
				http.Error(w, "Failed to get user role", http.StatusInternalServerError)
				return
			}

			if role != "teacher" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// User is authorized, add the token to the request context for use in protected endpoints
			ctx := context.WithValue(r.Context(), "user", tokenString)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else if resp.StatusCode == http.StatusUnauthorized {
			// Token is invalid or expired

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {

			http.Error(w, "Error", http.StatusInternalServerError)
		}
	})
}
