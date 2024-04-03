FROM node:20.11.1 as NODE_BUILDER
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY tsconfig.app.json tsconfig.json tsconfig.node.json vite.config.ts index.html env.d.ts ./
COPY frontend ./frontend
COPY public ./public
RUN npm run build

FROM golang:1.22.1 as GO_BUILDER
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY pkg ./pkg
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /game-server

FROM alpine:3.19
COPY --from=GO_BUILDER /game-server /game-server
COPY --from=NODE_BUILDER /app/.output /public
COPY --from=GO_BUILDER /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
CMD ["/game-server"]