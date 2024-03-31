package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pb "main/proto/data_service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type Task struct {
	Id          int64  `json:"id"`
	AuthorId    int64  `json:"author_id"`
	Description string `json:"description"`
	Status      string `json:"string"`
	CreatedAt   int64  `json:"created_at"`
}

var db *sql.DB
var grpcDataClient pb.TaskDataClient
var ctx context.Context

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

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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

func checkLogin(w http.ResponseWriter, decoder *json.Decoder) *int64 {
	var creds UserCreds

	err := decoder.Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	var user_id int64
	err = db.QueryRow("SELECT id FROM users WHERE login = $1 AND phash = $2", creds.Login, ComputeHash(creds.Password)).Scan(&user_id)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusForbidden)
		return nil
	}

	return &user_id
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user_id := checkLogin(w, decoder)
	if user_id == nil {
		return
	}

	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := grpcDataClient.CreateTask(ctx, &pb.Task{
		AuthorId:    *user_id,
		Description: task.Description,
		Status:      task.Status,
		CreatedAt:   timestamppb.Now(),
	})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user_id := checkLogin(w, decoder)
	if user_id == nil {
		return
	}
	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := grpcDataClient.UpdateTask(ctx, &pb.UpdateTaskRequest{
		Task: &pb.Task{
			AuthorId:    task.AuthorId,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   timestamppb.New(time.Unix(task.CreatedAt, 0)),
		},
		UserId: *user_id,
	})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user_id := checkLogin(w, decoder)
	if user_id == nil {
		return
	}
	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := grpcDataClient.DeleteTask(ctx, &pb.DeleteTaskRequest{
		TaskId: task.Id,
		UserId: *user_id,
	})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != 0 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user_id := checkLogin(w, decoder)
	if user_id == nil {
		return
	}
	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := grpcDataClient.GetTask(ctx, &pb.GetTaskRequest{
		TaskId: task.Id,
		UserId: *user_id,
	})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(Task{
		Id:          resp.Task.TaskId,
		AuthorId:    resp.Task.AuthorId,
		Description: resp.Task.Description,
		Status:      resp.Task.Status,
		CreatedAt:   resp.Task.CreatedAt.Seconds,
	})
	w.WriteHeader(http.StatusOK)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	user_id := checkLogin(w, decoder)
	if user_id == nil {
		return
	}
	resp, err := grpcDataClient.GetTasks(ctx, &pb.GetTasksRequest{})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tasks := []Task{}
	for _, task := range resp.Tasks {
		tasks = append(tasks, Task{
			Id:          task.TaskId,
			AuthorId:    task.AuthorId,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt.Seconds,
		})
	}
	json.NewEncoder(w).Encode(tasks)
	w.WriteHeader(http.StatusOK)
}

func main() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))
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

	conn, err := grpc.Dial(os.Getenv("DATA_SERVICE_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcDataClient = pb.NewTaskDataClient(conn)
	ctx, _ = context.WithTimeout(context.Background(), time.Second)

	router := mux.NewRouter()
	router.HandleFunc("/user/create", CreateUser).Methods("POST")
	router.HandleFunc("/user/update", UpdateUser).Methods("PUT")
	router.HandleFunc("/user/login", UserLogin).Methods("POST")

	router.HandleFunc("/post/create", CreateTask).Methods("POST")
	router.HandleFunc("/post/update", UpdateTask).Methods("PUT")
	router.HandleFunc("/post/delete", DeleteTask).Methods("POST")
	router.HandleFunc("/post/get_task", GetTask).Methods("GET")
	router.HandleFunc("/post/get_tasks", GetTasks).Methods("GET")

	http.ListenAndServe(":8080", router)
	select {}
}
