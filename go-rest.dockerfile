# builder img
FROM golang:1.22-alpine as builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CG_ENABLED=0 go build -o gorest ./cmd/api
RUN chmod +x /app/gorest

# tiny exec img
FROM alpine:latest
RUN mkdir /app
COPY --from=builder /app/gorest /app
EXPOSE 8001
CMD ["/app/gorest"]