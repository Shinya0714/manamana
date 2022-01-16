FROM golang:latest
COPY ./go/app /app/go
WORKDIR /app/go
RUN go mod init github.com/Shinya0714/testgo && go build -o main . 
CMD ["/app/go/main"]