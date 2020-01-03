FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go get -d -v
RUN go build -o /go/bin/app

FROM scratch
COPY --from=builder /go/bin/app /go/bin/app
ENTRYPOINT ["/go/bin/app"]