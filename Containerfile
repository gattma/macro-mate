FROM golang:1.23-alpine as build-env
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

RUN GOOS=linux GOARCH=amd64 go build -o /macro-mate github.com/gattma/macro-mate/src/cmd/macro-mate-server

FROM alpine:3.14

RUN apk update \
    && apk upgrade\
    && apk add --no-cache tzdata curl

WORKDIR /app
COPY --from=build-env /macro-mate .
COPY --from=build-env /app/src/cmd/macro-mate-server /app/
COPY web /app/web

EXPOSE 80
CMD [ "./macro-mate" ]