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


FROM alpine:latest

WORKDIR /app

# 로컬에서 빌드한 실행 파일을 현재 작업 디렉토리로 복사
COPY main .

# 실행 파일이 존재하는지 확인
RUN if [ ! -f "./main" ]; then echo "Error: main file not found"; exit 1; fi

# 실행 파일에 실행 가능한 권한 부여
RUN chmod +x main

# 실행 파일의 권한 확인
RUN ls -l main

# ENTRYPOINT로 실행 파일을 설정
ENTRYPOINT ["./main"]
