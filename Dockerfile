FROM golang:1.15-alpine3.12 as go-builder
COPY . ./
RUN apk add --update git
RUN go build -o merge-me .
EXPOSE 10080 10433
ENTRYPOINT ["./merge-me"]
