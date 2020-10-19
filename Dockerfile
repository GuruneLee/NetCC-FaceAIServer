FROM dlckdgk4858/go-face:1.0

RUN mkdir /go/src/app
WORKDIR /go/src/app

RUN go mod init Face-AI-server


COPY main.go main.go

ENV MODEL_DIR /go/src/go-face-example/testdata/models

EXPOSE 8080

ENTRYPOINT ["go", "run", "main.go"]
