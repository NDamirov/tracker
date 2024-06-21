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

	"github.com/IBM/sarama"

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
	Status      string `json:"status"`
	CreatedAt   int64  `json:"created_at"`
}

type TaskIdRequest struct {
	Id int64 `json:"task_id"`
}

type GetTasksRequest struct {
	PageNum        int `json:"page_number"`
	ResultsPerPage int `json:"results_per_page"`
}

var db *sql.DB
var grpcDataClient pb.TaskDataClient
var producer sarama.SyncProducer

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
		http.Error(w, "User with login already exists", http.StatusForbidden)
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

		_, err = db.Exec("UPDATE users SET name = $1, surname = $2, bdate = to_date($3, 'YYYY-MM-DD'), email = $4, phoneno = $5 WHERE login = $6",
			userInfo.Name, userInfo.Surname, userInfo.Birth, userInfo.Email, userInfo.Phone, login)
		if err != nil {
			http.Error(w, fmt.Sprint("Invalid request ", err), http.StatusForbidden)
			return
		}
	} else {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func checkLogin(w http.ResponseWriter, r *http.Request) *int64 {
	tokenString := r.Header.Get("token")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}

		return []byte(JwtSecret), nil
	})

	var userId int64
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		login := claims["login"]

		is_ok := 1
		err := db.QueryRow("SELECT COUNT(*) FROM user_tokens WHERE login = $1 AND token = $2", login, tokenString).Scan(&is_ok)
		if err != nil || is_ok == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil
		}

		err = db.QueryRow("SELECT id FROM users WHERE login = $1", login).Scan(&userId)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return nil
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}

	return &userId
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId := checkLogin(w, r)
	if userId == nil {
		return
	}

	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := grpcDataClient.CreateTask(context.Background(), &pb.Task{
		AuthorId:    *userId,
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
		log.Printf("Not zero status code: %s\n", resp.Message)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId := checkLogin(w, r)
	if userId == nil {
		return
	}
	var task Task
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := grpcDataClient.UpdateTask(context.Background(), &pb.UpdateTaskRequest{
		Task: &pb.Task{
			TaskId:      task.Id,
			AuthorId:    task.AuthorId,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   timestamppb.New(time.Unix(task.CreatedAt, 0)),
		},
		UserId: *userId,
	})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if resp.StatusCode != 0 {
		http.Error(w, resp.Message, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId := checkLogin(w, r)
	if userId == nil {
		return
	}
	var task TaskIdRequest
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := grpcDataClient.DeleteTask(context.Background(), &pb.DeleteTaskRequest{
		TaskId: task.Id,
		UserId: *userId,
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
	userId := checkLogin(w, r)
	if userId == nil {
		return
	}
	var task TaskIdRequest
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := grpcDataClient.GetTask(context.Background(), &pb.GetTaskRequest{
		TaskId: task.Id,
		UserId: *userId,
	})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch resp.Error.StatusCode {
	case 0:
	case 2, 3:
		http.Error(w, resp.Error.Message, http.StatusForbidden)
		return
	default:
		http.Error(w, resp.Error.Message, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Task{
		Id:          resp.Task.TaskId,
		AuthorId:    resp.Task.AuthorId,
		Description: resp.Task.Description,
		Status:      resp.Task.Status,
		CreatedAt:   resp.Task.CreatedAt.Seconds,
	})
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	userId := checkLogin(w, r)
	if userId == nil {
		return
	}

	var req GetTasksRequest
	err := decoder.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := grpcDataClient.GetTasks(context.Background(), &pb.GetTasksRequest{
		UserId:         *userId,
		PageNumber:     int32(req.PageNum),
		ResultsPerPage: int32(req.ResultsPerPage),
	})
	if err != nil {
		log.Printf("GRPC failed: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch resp.Error.StatusCode {
	case 0:
	case 2, 3:
		http.Error(w, resp.Error.Message, http.StatusForbidden)
		return
	default:
		http.Error(w, resp.Error.Message, http.StatusInternalServerError)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func react(w http.ResponseWriter, r *http.Request, t string) {
	msg := &sarama.ProducerMessage{
		Topic: "my_topic",
		Value: sarama.ByteEncoder(t),
	}

	_, _, err := producer.SendMessage(msg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func LikeTask(w http.ResponseWriter, r *http.Request) {
	react(w, r, "like")
}

func ViewTask(w http.ResponseWriter, r *http.Request) {
	react(w, r, "view")
}

func main() {
	var err error

	brokers := []string{"kafka:29092"}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal(err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))

	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Println(err)
			if i == 4 {
				panic(err)
			}
		} else {
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}

	defer db.Close()

	conn, err := grpc.Dial(os.Getenv("DATA_SERVICE_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcDataClient = pb.NewTaskDataClient(conn)

	router := mux.NewRouter()
	router.HandleFunc("/user/create", CreateUser).Methods("POST")
	router.HandleFunc("/user/update", UpdateUser).Methods("PUT")
	router.HandleFunc("/user/login", UserLogin).Methods("POST")

	router.HandleFunc("/task/create", CreateTask).Methods("POST")
	router.HandleFunc("/task/update", UpdateTask).Methods("PUT")
	router.HandleFunc("/task/delete", DeleteTask).Methods("POST")
	router.HandleFunc("/task/get_task", GetTask).Methods("GET")
	router.HandleFunc("/task/get_tasks", GetTasks).Methods("GET")

	router.HandleFunc("/react/like", LikeTask).Methods("POST")
	router.HandleFunc("/react/view", ViewTask).Methods("POST")

	http.ListenAndServe(":8080", router)
	select {}
}
