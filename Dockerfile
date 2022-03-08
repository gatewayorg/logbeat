FROM golang:1.15-alpine as builder
LABEL builder=gateway-logbeat
WORKDIR /src
RUN go env -w  GOPROXY=https://goproxy.cn,direct

ADD . .

RUN CGO_ENABLE=0 GOOS=linux go build -race -o ./build/server ./cmd/server/

FROM alpine
COPY --from=builder /src/build/ /
COPY entrypoint.sh entrypoint.sh
RUN chmod +x entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]

