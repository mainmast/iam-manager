FROM golang:1.13.5-alpine3.11 AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .
RUN ls -l

RUN export GO111MODULE=on
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/iam-manager

WORKDIR /app/cmd/iam-migrations
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/migrate

FROM alpine
WORKDIR /app

COPY --from=builder /app/startup.sh /app/startup.sh
RUN chmod +x /app/startup.sh

COPY --from=builder /app/iam-manager /app/iam-manager
COPY --from=builder /app/migrate /app/migrate
COPY --from=builder /app/cmd /app/cmd
COPY --from=builder /app/conf /app/conf

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd
USER nobody

ENTRYPOINT [ "/app/startup.sh" ]