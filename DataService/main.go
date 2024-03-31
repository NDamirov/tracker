package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "main/proto/data_service"
)

var db *sql.DB

func main() {
	var err error
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_NAME"))
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaskDataServer(grpcServer, &server{})

	err = grpcServer.Serve(listener)
	panic(err)
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
	_, err := db.Exec("INSERT INTO tasks(author_id, task_status, task_description, created_at) VALUES ($1, $2, $3, $4)",
		task.AuthorId, task.Status, task.Description, task.CreatedAt)
	if err != nil {
		grpclog.Errorf("Failed to insert task: %s", err)
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
	task, err := getTask(request.Task.TaskId)
	if err != nil {
		grpclog.Errorf("Failed to get task: %s", err)
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
	_, err = db.Exec("UPDATE tasks SET description = $1, status = $2 WHERE author_id = $3 AND task_id = $4", request.Task.Description, request.Task.Status, task.AuthorId, task.TaskId)
	if err != nil {
		grpclog.Errorf("Failed to update task: %s", err)
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
	task, err := getTask(request.TaskId)
	if err != nil {
		grpclog.Errorf("Failed to get task: %s", err)
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
		grpclog.Errorf("Failed to delete task: %s", err)
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
	task, err := getTask(request.TaskId)
	if err != nil {
		grpclog.Errorf("Failed to get task: %s", err)
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
	rows, err := db.Query("SELECT task_id FROM tasks WHERE author_id = $1 ORDER BY task_id DESC OFFSET $2 LIMIT $3", request.UserId, request.ResultsPerPage*request.PageNumber, request.ResultsPerPage)
	tasks := []*pb.Task{}
	if err != nil {
		grpclog.Errorf("Failed to insert task: %s", err)
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
	return nil, nil
}
