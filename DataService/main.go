package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "main/proto/data_service"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error
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

	grpcAddr := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	listener, err := net.Listen("tcp", grpcAddr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskDataServer(grpcServer, &server{})
	reflection.Register(grpcServer)

	log.Printf("Starting grpc server at '%s'\n", grpcAddr)
	err = grpcServer.Serve(listener)
	log.Printf("Grpc server failed: %s\n", err)
}

type server struct {
	pb.UnimplementedTaskDataServer
}

type Task struct {
	TaskId      int64
	AuthorId    int64
	Status      string
	Description string
	CreatedAt   *timestamppb.Timestamp
}

func getTask(taskId int64) (*Task, error) {
	row := db.QueryRow("SELECT task_id, author_id, task_status, task_description, created_at FROM tasks WHERE task_id = $1", taskId)
	var task Task
	var createdAt int64
	switch err := row.Scan(&task.TaskId, &task.AuthorId, &task.Status, &task.Description, &createdAt); err {
	case nil:
	case sql.ErrNoRows:
		return nil, nil
	default:
		return nil, err
	}
	task.CreatedAt = timestamppb.New(time.Unix(createdAt, 0))

	return &task, nil
}

func (s *server) CreateTask(c context.Context, task *pb.Task) (*pb.Error, error) {
	log.Println("Creating task")
	_, err := db.Exec("INSERT INTO tasks(author_id, task_status, task_description, created_at) VALUES ($1, $2, $3, $4)",
		task.AuthorId, task.Status, task.Description, task.CreatedAt.AsTime().Unix())
	if err != nil {
		log.Printf("Failed to insert task: %s", err)
		return &pb.Error{
			StatusCode: 1,
			Message:    "Internal error",
		}, nil
	}
	return &pb.Error{
		StatusCode: 0,
		Message:    "Ok",
	}, nil
}

func (s *server) UpdateTask(c context.Context, request *pb.UpdateTaskRequest) (*pb.Error, error) {
	log.Printf("Updating task")
	task, err := getTask(request.Task.TaskId)
	if err != nil {
		log.Printf("Failed to get task: %s", err)
		return &pb.Error{
			StatusCode: 1,
			Message:    "Internal error",
		}, nil
	}
	if task == nil {
		return &pb.Error{
			StatusCode: 2,
			Message:    "Invalid id",
		}, nil
	}
	if task.AuthorId != request.UserId {
		return &pb.Error{
			StatusCode: 3,
			Message:    "Not enough rights",
		}, nil
	}
	_, err = db.Exec("UPDATE tasks SET task_description = $1, task_status = $2 WHERE author_id = $3 AND task_id = $4", request.Task.Description, request.Task.Status, task.AuthorId, task.TaskId)
	if err != nil {
		log.Printf("Failed to update task: %s", err)
		return &pb.Error{
			StatusCode: 1,
			Message:    "Internal error",
		}, nil
	}
	return &pb.Error{
		StatusCode: 0,
		Message:    "Ok",
	}, nil
}

func (s *server) DeleteTask(c context.Context, request *pb.DeleteTaskRequest) (*pb.Error, error) {
	log.Printf("Deleting task")
	task, err := getTask(request.TaskId)
	if err != nil {
		log.Printf("Failed to get task: %s", err)
		return &pb.Error{
			StatusCode: 1,
			Message:    "Internal error",
		}, nil
	}
	if task == nil {
		return &pb.Error{
			StatusCode: 2,
			Message:    "Invalid id",
		}, nil
	}
	if task.AuthorId != request.UserId {
		return &pb.Error{
			StatusCode: 3,
			Message:    "Not enough rights",
		}, nil
	}
	_, err = db.Exec("DELETE FROM tasks WHERE task_id = $1", request.TaskId)
	if err != nil {
		log.Printf("Failed to delete task: %s", err)
		return &pb.Error{
			StatusCode: 1,
			Message:    "Internal error",
		}, err
	}
	return &pb.Error{
		StatusCode: 0,
		Message:    "Ok",
	}, nil
}

func (s *server) GetTask(c context.Context, request *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	log.Println("Getting task")
	task, err := getTask(request.TaskId)
	if err != nil {
		log.Printf("Failed to get task: %s", err)
		return &pb.GetTaskResponse{
			Error: &pb.Error{
				StatusCode: 1,
				Message:    "Internal error",
			},
		}, nil
	}
	if task == nil {
		return &pb.GetTaskResponse{
			Error: &pb.Error{
				StatusCode: 2,
				Message:    "Invalid id",
			},
		}, nil
	}
	if task.AuthorId != request.UserId {
		return &pb.GetTaskResponse{
			Error: &pb.Error{
				StatusCode: 3,
				Message:    "Not enough rights",
			},
		}, nil
	}
	return &pb.GetTaskResponse{
		Error: &pb.Error{
			StatusCode: 0,
			Message:    "Ok",
		},
		Task: &pb.Task{
			TaskId:      task.TaskId,
			AuthorId:    task.AuthorId,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
		},
	}, nil
}

func (s *server) GetTasks(c context.Context, request *pb.GetTasksRequest) (*pb.GetTasksResponse, error) {
	log.Println("Getting tasks")
	rows, err := db.Query("SELECT task_id FROM tasks WHERE author_id = $1 ORDER BY task_id DESC OFFSET $2 LIMIT $3", request.UserId, request.ResultsPerPage*request.PageNumber, request.ResultsPerPage)
	tasks := []*pb.Task{}
	if err != nil {
		log.Printf("Failed to insert task: %s", err)
		return &pb.GetTasksResponse{
			Error: &pb.Error{
				StatusCode: 1,
				Message:    "Internal error",
			},
			Tasks: tasks,
		}, nil
	}
	for rows.Next() {
		var taskId int64
		rows.Scan(&taskId)
		task, _ := getTask(taskId)
		tasks = append(tasks, &pb.Task{
			TaskId:      task.TaskId,
			AuthorId:    task.AuthorId,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
		})
	}
	return &pb.GetTasksResponse{
		Error: &pb.Error{
			StatusCode: 0,
			Message:    "Ok",
		},
		Tasks: tasks,
	}, nil
}
