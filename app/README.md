# Titanic

Web server that can be used to manage passengers that were aboard the Titanic. For API documentation please see [swagger.yml](https://github.com/labrenbe/titanic/-/blob/main/app/swagger.yml).

## Requirements

- Go >= 1.19 installed
- Connection to a postgres database
- This directory should be on your GOPATH

## Starting the titanic server locally

```
go run titanic.go
```

## Expected environment variables
- DB_HOST
- DB_PORT
- DB_USER
- DB_PASSWORD
- DB_NAME
