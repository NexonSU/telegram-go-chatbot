FROM golang:1.19 AS builder

RUN apt update
RUN apt install ca-certificates
RUN update-ca-certificates

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=builder /app .
ADD files /files

ENTRYPOINT ["/app"]