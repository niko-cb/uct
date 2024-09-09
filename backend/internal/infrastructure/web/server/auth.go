package server

import (
	echojwt "github.com/labstack/echo-jwt/v4"
)

// Auth is a middleware to authenticate the user using JWT
func (s *server) Auth(secretKey string) {
	// The secret key would normally be stored in a secure location.
	// This is just for demonstration purpose
	// The secret key could be stored in a cloud service like AWS Secrets Manager or GCP Secret Manager
	// Or in another secure location like a Kubernetes secret or a .env file, etc.
	s.Use(echojwt.JWT([]byte(secretKey)))
}
