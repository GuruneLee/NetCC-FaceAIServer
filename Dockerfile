FROM dlckdgk4858/go-face:1.0

RUN mkdir /app
WORKDIR /app

RUN go mod init Face-AI-server

COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]
