FROM golang:1.22.2-alpine3.19 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /build

COPY go.mod main.go ./

# 필요한 경우 git 및 openssh 설치
RUN apk add --no-cache git openssh

# SSH 키를 복사하여 Docker 이미지 내에 추가
COPY ~/.ssh/id_rsa /root/.ssh/id_rsa

# SSH 호스트 키를 인증된 호스트 목록에 추가
RUN ssh-keyscan -t rsa github.com >> /root/.ssh/known_hosts

# SSH 키 권한 설정
RUN chmod 600 /root/.ssh/id_rsa

RUN go env -w GOPRIVATE=github.com/TeamWAF

RUN go mod tidy

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main .

ENTRYPOINT ["/main"]
