FROM golang:1.20 AS go-build

ENV GO111MODULE=on

WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

# ENV GOPROXY=direct

RUN go mod download

COPY . .

# RUN go test ./internal/tests -v

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash
WORKDIR /root/
COPY --from=go-build /app .

CMD ["./main"]

EXPOSE 3000