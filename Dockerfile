FROM golang:1.16-alpine
RUN apk add build-base
WORKDIR /app
COPY * ./
COPY assets/img ./assets/img
COPY templates ./templates
RUN go build .
EXPOSE 8080
CMD ["./psyme","--web"]
