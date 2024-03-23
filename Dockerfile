FROM node:lts as NODE_BUILDER
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM golang:latest as GO_BUILDER
WORKDIR /app
COPY go.sum go.mod ./
RUN go mod download
COPY . .
COPY --from=NODE_BUILDER /app/.output/* ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /game-server

FROM scratch
COPY --from=GO_BUILDER /game-server /game-server
COPY --from=GO_BUILDER /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080
CMD ["/game-server"]