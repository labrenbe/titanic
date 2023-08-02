FROM golang:1.17-alpine AS builder

WORKDIR $GOPATH/titanic

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY titanic.go ./

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o /app/titanic


FROM scratch

USER 1000
ENV GIN_MODE=release

COPY --from=builder /app/titanic /app/titanic

EXPOSE 8080

ENTRYPOINT ["/app/titanic"]
