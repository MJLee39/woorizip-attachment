# FROM golang:1.22.2-alpine3.19 AS builder

# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=arm64

# WORKDIR /build

# COPY go.mod main.go ./

# RUN go env -w GOPRIVATE=github.com/TeamWAF

# RUN go mod tidy

# RUN go build -o main .

# WORKDIR /dist

# RUN cp /build/main .

# FROM scratch

# COPY --from=builder /dist/main .

# ENTRYPOINT ["/main"]


FROM scratch

WORKDIR /app

COPY ./main .

ENTRYPOINT ["./main"]