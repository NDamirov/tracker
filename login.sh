curl -X 'POST' -v \
  '127.0.0.1:8095/user/create' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "login": "user",
  "password": "password"
}'

curl -X 'POST' -v \
  '127.0.0.1:8095/user/login' \
  -H 'accept: */*' \
  -H 'Content-Type: application/json' \
  -d '{
  "login": "user",
  "password": "password"
}'

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

# curl -X 'POST' \
#   '127.0.0.1:8095/task/create' \
#   -H 'accept: */*' \
#   -H 'token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE3MTI5MjAxMTUsImxvZ2luIjoidXNlciJ9.uDEzk_mg0Q50sDJ7yirlyW_CIF1aYjGa0KIwQjYF3zQ' \
#   -H 'Content-Type: application/json' \
#   -d '{
#   "description": "string",
#   "status": "string",
#   "created_at": 0
# }' -v