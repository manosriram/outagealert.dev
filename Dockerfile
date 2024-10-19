FROM golang:1.23-bullseye
WORKDIR /build

COPY ./go.mod .

ENV GO111MODULE=on
ENV APP_HOME=/outagealert
ENV PGHOST=/var/run/postgresql
RUN mkdir -p "$APP_HOME"

RUN sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates curl gnupg
RUN curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | sudo gpg --dearmor -o /usr/share/keyrings/doppler-archive-keyring.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" | sudo tee /etc/apt/sources.list.d/doppler-cli.list
RUN sudo apt-get update && sudo apt-get install doppler
RUN doppler configure set token $DOPPLER_TOKEN

RUN go mod tidy
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

EXPOSE 1323
RUN go build -o outagealert
