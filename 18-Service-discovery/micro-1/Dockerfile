FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Copy the source code into the container
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64


RUN go build -o micro-1 .

# Stage 2: Create the runtime image
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/micro-1 /app/micro-1

EXPOSE 80

CMD ["/app/micro-1"]