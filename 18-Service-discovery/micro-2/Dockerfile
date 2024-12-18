FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Copy the source code into the container
COPY . .

ENV GOOS=linux
ENV GOARCH=amd64


RUN go build -o micro-2 .

# Stage 2: Create the runtime image
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/micro-2 /app/micro-2

EXPOSE 80

CMD ["/app/micro-2"]