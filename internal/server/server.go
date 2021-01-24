package server

import (
	"encoding/json"
	"fmt"
	"github.com/mattcarabine/perugo/internal/perugo"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	rooms map[string]perugo.Room
}

var (
	secret      []byte
	server      = Server{rooms: map[string]perugo.Room{}}
	tokenExpiry = time.Hour * 10
)

func generateJwt() (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// Create the Claims
	now := time.Now()
	claims := jwt.StandardClaims{
		ExpiresAt: now.Add(tokenExpiry).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
		Issuer:    "perugo-api",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(secret)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tokenString, err := generateJwt()
	if err != nil {
		zap.S().Error("failed to generate secret token", zap.Error(err), zap.String("remote", r.RemoteAddr))
		return
	}
	zap.S().Info("generated token", zap.String("remote", r.RemoteAddr))

	_, err = w.Write([]byte(tokenString))
	if err != nil {
		zap.S().Error("failed to write secret token", zap.Error(err), zap.String("remote", r.RemoteAddr))
	}
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")

		valid, err := validateJWT(auth)
		if !valid || err != nil {
			zap.S().Warn("invalid JWT token", zap.Error(err), zap.String("remote", r.RemoteAddr))
			w.WriteHeader(401)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateJWT(tokenString string) (bool, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})

	if err != nil {
		zap.S().Warn("unable to parse JWT", zap.String("token", tokenString), zap.Error(err))
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	} else {
		return false, nil
	}
}

func roomHandler(w http.ResponseWriter, r *http.Request) {
	rawBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		zap.S().Warn("failed to read room body", zap.Error(err), zap.Any("request", r))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var room perugo.Room
	err = json.Unmarshal(rawBody, &room)
	if err != nil {
		zap.S().Warn("failed to unmarshal body to a Room", zap.Error(err), zap.ByteString("body", rawBody))
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if room.Id == "" {
		zap.S().Warn("no id found in create room payload")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server.rooms[room.Id] = room

	response, err := json.Marshal(room)
	if err != nil {
		zap.S().Error("failed to marshal room", zap.Error(err))
		return
	}

	_, err = w.Write(response)
	if err != nil {
		zap.S().Error("failed to send response", zap.Error(err), zap.String("remote", r.RemoteAddr), zap.ByteString("response", response))
	}
}

func SetupServer(addr, signingSecret string) error {
	secret = []byte(signingSecret)

	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods(http.MethodPost)
	api := r.PathPrefix("/api/").Subrouter()
	api.Use(jwtMiddleware)
	api.HandleFunc("/room", roomHandler).Methods(http.MethodPost)
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	zap.S().Info("started server, now listening", zap.String("address", addr), zap.String("test", "value"))

	return srv.ListenAndServe()
}
