### build stage
FROM golang:1.18 AS builder

ENV GO111MODULE auto
ENV CGO_ENABLED=0

# modを先に入れてキャッシュをきかせる
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download

# その後ディレクトリを全てコピーする
COPY . /src
WORKDIR /src/cmd/famili-api
RUN go build -o /src/bin/famili-api

### final stage
FROM scratch

WORKDIR /app
COPY --from=builder /src/bin/famili-api /app/famili-api
ENTRYPOINT ["/app/famili-api"]
