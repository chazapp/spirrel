FROM golang:1.22 AS builder
ARG VERSION
RUN test -n "$VERSION"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags="-X main.Version=$VERSION" -o spirrel

FROM alpine:3.20.2
WORKDIR /app
COPY --from=builder /app /app
RUN adduser -D -H -h /app appuser
RUN chown -R appuser:appuser /app
EXPOSE 8080
USER appuser
ENTRYPOINT ["/app/spirrel"]