FROM golang:1.23-bullseye
WORKDIR /build

COPY ./go.mod .

ENV GO111MODULE=on
ENV APP_HOME=/outagealert
ENV PGHOST=/var/run/postgresql
RUN mkdir -p "$APP_HOME"

ARG DOPPLER_TOKEN
ENV DOPPLER_TOKEN=${DOPPLER_TOKEN}

RUN go mod tidy
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN yes | apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg && \
curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | gpg --dearmor -o /usr/share/keyrings/doppler-archive-keyring.gpg && \
echo "deb [signed-by=/usr/share/keyrings/doppler-archive-keyring.gpg] https://packages.doppler.com/public/cli/deb/debian any-version main" | tee /etc/apt/sources.list.d/doppler-cli.list && \
apt-get update && apt-get install doppler
RUN doppler configure set token $DOPPLER_TOKEN

COPY . .

EXPOSE 1323
RUN go build -o outagealert
