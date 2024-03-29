FROM golang:1.15-alpine3.12 as go-builder
WORKDIR /app
COPY . ./
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o merge-me cmd/main.go
EXPOSE 10080 10433
ENTRYPOINT ["./merge-me"]

FROM gcr.io/distroless/base-debian10
COPY --from=go-builder /app/merge-me /
COPY --from=go-builder /app/default.yaml /
CMD ["/merge-me"]
