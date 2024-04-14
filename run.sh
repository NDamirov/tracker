# curl -X 'POST' -v \
#   '127.0.0.1:8095/user/create' \
#   -H 'accept: */*' \
#   -H 'Content-Type: application/json' \
#   -d '{
#   "login": "user",
#   "password": "password"
# }'

# curl -X 'POST' -v \
#   '127.0.0.1:8095/user/login' \
#   -H 'accept: */*' \
#   -H 'Content-Type: application/json' \
#   -d '{
#   "login": "user",
#   "password": "password"
# }'

# curl -X 'PUT' \
#   '127.0.0.1:8095/user/update' \
#   -H 'accept: */*' \
#   -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE3MTI2NTkxMDksImxvZ2luIjoidXNlciJ9.YBpfl4SvKX9QGr7fd-jEyNiaGGUY_ABkpqt4gr5NR9Y' \
#   -H 'Content-Type: application/json' \
#   -d '{
#   "name": "Ivan",
#   "surname": "Ivanov",
#   "birth": "2024-04-08",
#   "email": "user@example.com",
#   "phone": "string"
# }' -v

curl -X 'POST' \
  '127.0.0.1:8095/task/create' \
  -H 'accept: */*' \
  -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE3MTI5MzM3NTQsImxvZ2luIjoidXNlciJ9.CxqIW3BPVoPw4BZpbW-Fp-vtDYIKZV_gzLEjKMH63B4' \
  -H 'Content-Type: application/json' \
  -d '{
  "description": "string",
  "status": "string",
  "created_at": 0
}' -v

curl -X 'PUT' \
  '127.0.0.1:8095/task/update' \
  -H 'accept: */*' \
  -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE3MTI5MzM3NTQsImxvZ2luIjoidXNlciJ9.CxqIW3BPVoPw4BZpbW-Fp-vtDYIKZV_gzLEjKMH63B4' \
  -H 'Content-Type: application/json' \
  -d '{
  "id": 1,
  "description": "string",
  "status": "string2",
  "created_at": 0
}' -v

curl -X 'POST' \
  '127.0.0.1:8095/task/delete' \
  -H 'accept: application/json' \
  -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE3MTI5MzM3NTQsImxvZ2luIjoidXNlciJ9.CxqIW3BPVoPw4BZpbW-Fp-vtDYIKZV_gzLEjKMH63B4' \
  -H 'Content-Type: application/json' \
  -d '{
  "task_id": 1
}' -v

curl -X 'GET' \
  '127.0.0.1:8095/task/get_task' \
  -H 'accept: application/json' \
  -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE3MTI5MzM3NTQsImxvZ2luIjoidXNlciJ9.CxqIW3BPVoPw4BZpbW-Fp-vtDYIKZV_gzLEjKMH63B4' \
  -H 'Content-Type: application/json' \
  -d '{
  "task_id": 1
}' -v

# curl -X 'GET' \
#   '127.0.0.1:8095/task/get_tasks' \
#   -H 'accept: application/json' \
#   -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE3MTI5MzIwODEsImxvZ2luIjoidXNlciJ9.WzMYJHK_FfTfRCBHLfeEZuv6kmJnojqMzLFFv3j-SNQ' \
#   -H 'Content-Type: application/json' \
#   -d '{
#   "page_number": 5,
#   "results_per_page": 2
# }' -v