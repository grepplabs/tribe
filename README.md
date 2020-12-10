# tribe - user management and identity server

## Getting started

```
make build

export UPPER_DB_LOG=TRACE
./tribe serve admin --server-port=8080
google-chrome http://127.0.0.1:8080/v1/docs
```

## Usage

Realms

```
curl -X POST -H 'Content-Type: application/json' -d '{"realm_id": "main"}' http://localhost:8080/v1/realms

curl http://localhost:8080/v1/realms/main
```

Users
```
curl -X POST -H 'Content-Type: application/json' -d '{"username": "michal", "password": "hello", "enabled": true}' http://localhost:8080/v1/realms/main/users
curl -X POST -H 'Content-Type: application/json' -d '{"username": "michal2", "password": "hello2", "enabled": false, "email": "michal2@example.com"}' http://localhost:8080/v1/realms/main/users

curl http://localhost:8080/v1/realms/main/users/michal

```

```
for i in {1..100}
do
  echo "{\"username\": \"user-${i}-$(uuidgen)\", \"password\": \"hello\"}" | curl -X POST -H 'Content-Type: application/json' -d @- http://localhost:8080/v1/realms/main/users
done
```

