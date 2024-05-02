FROM golang:1.22.2-alpine3.19 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

# 필요한 경우 git 설치
RUN apk add --no-cache git

# GitHub 액세스 토큰을 환경 변수로 설정
ENV GITHUB_TOKEN=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

WORKDIR /build

COPY go.mod main.go ./

RUN git config --global url."https://${GITHUB_TOKEN}@github.com/".insteadOf "https://github.com/"
# GitHub 액세스 토큰을 사용하여 모듈 다운로드
RUN go env -w GOPRIVATE=github.com/TeamWAF
RUN go mod tidy

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM scratch

COPY --from=builder /dist/main .

ENTRYPOINT ["/main"]
