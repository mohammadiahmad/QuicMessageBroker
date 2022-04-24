FROM golang:1.17-alpine as builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build  -o /app/main ./cmd


#################
FROM scratch
WORKDIR /app
COPY --from=builder /app/main /app/main
ENTRYPOINT ["/app/main"]