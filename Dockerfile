FROM golang:1.23-bullseye
WORKDIR /build

COPY ./go.mod .

ENV GO111MODULE=on
ENV APP_HOME=/outagealert
ENV PGHOST=/var/run/postgresql
RUN mkdir -p "$APP_HOME"

RUN go mod tidy
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

EXPOSE 1323
# RUN make sqlcgen
CMD ["go", "run", "main.go"]
