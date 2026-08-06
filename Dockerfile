FROM node:24.15.0 AS node_builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend .
RUN npm run build

FROM golang:1.26.2 AS go_builder
WORKDIR /app
COPY backend/go.sum backend/go.mod .
RUN go mod download
COPY backend .
RUN CGO_ENABLED=0 GOOS=linux go build -o /game-server

FROM alpine:3.23.4
COPY --from=go_builder /game-server /game-server
COPY --from=node_builder /app/.output /public
COPY --from=go_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
CMD ["/game-server"]