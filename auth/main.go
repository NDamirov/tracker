package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const JwtSecret = "SECRET"

type UserCreds struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserInfo struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Birth   string `json:"birth"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
}

type AuthToken struct {
	Token string `json:"token"`
}

var db *sql.DB

func ComputeHash(message string) string {
	hash := sha256.New()
	hash.Write([]byte(message))
	return hex.EncodeToString(hash.Sum(nil))
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var creds UserCreds

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO users (login, phash) VALUES ($1, $2)", creds.Login, ComputeHash(creds.Password))
	if err != nil {
		http.Error(w, fmt.Sprint("User with login already exists", err), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var creds UserCreds

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var login string
	err = db.QueryRow("SELECT login FROM users WHERE login = $1 AND phash = $2", creds.Login, ComputeHash(creds.Password)).Scan(&login)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusForbidden)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login":  login,
		"expiry": time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	})

	tokenString, err := token.SignedString([]byte(JwtSecret))
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO user_tokens(login, token) VALUES ($1, $2) ON CONFLICT (login) DO UPDATE SET token=EXCLUDED.token", login, tokenString)
	if err != nil {
		http.Error(w, "Failed to store token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AuthToken{Token: tokenString})
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userInfo UserInfo

	err := json.NewDecoder(r.Body).Decode(&userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenString := r.Header.Get("token")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}

		return []byte(JwtSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		login := claims["login"]

		is_ok := 0
		err = db.QueryRow("SELECT COUNT(*) FROM user_tokens WHERE login = $1 AND token = $2", login, tokenString).Scan(&is_ok)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusInternalServerError)
			return
		}
		if is_ok != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = db.Exec("UPDATE users SET name = $1, surname = $2, birth = $3, email = $4, phone = $5 WHERE login = $6",
			userInfo.Name, userInfo.Surname, userInfo.Birth, userInfo.Email, userInfo.Phone, login)
		if err != nil {
			http.Error(w, "Invalid request", http.StatusForbidden)
			return
		}
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))
	fmt.Println(connStr)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	for i := 0; i < 5; i++ {
		err = db.Ping()
		if err != nil {
			log.Println(err)
			if i == 4 {
				panic(err)
			}
		}
		time.Sleep(1000 * time.Millisecond)
	}
	fmt.Sprintln("OK")

	router := mux.NewRouter()
	router.HandleFunc("/user/create", CreateUser).Methods("POST")
	router.HandleFunc("/user/update", UpdateUser).Methods("PUT")
	router.HandleFunc("/user/login", UserLogin).Methods("POST")

	http.ListenAndServe(":8080", router)
	select {}
}
