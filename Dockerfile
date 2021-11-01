FROM golang:1.15.12-buster
WORKDIR /app
COPY . .
RUN apt update && apt install -y net-tools nmap && go build -o tg_bot cmd/main.go
CMD ["/app/tg_bot"]