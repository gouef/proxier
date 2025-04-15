FROM golang:1.24
WORKDIR /src
COPY *.go .
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/proxier ./main.go

FROM alpine:latest

WORKDIR /app
COPY --from=0 /app/proxier /app/proxier
COPY config.yaml .
RUN chmod +x proxier

EXPOSE 80
EXPOSE 443
CMD ["/app/proxier"]


