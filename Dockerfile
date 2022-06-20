FROM golang:alpine AS build
WORKDIR /build
COPY *go* .
RUN go mod tidy
RUN go build -o main .

FROM golang:alpine
RUN apk update && apk add bash
WORKDIR /app
COPY --from=build /build/main .
COPY config.yaml /app/config.yaml
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
