# build stage
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN apk add --no-cache git
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/ldap-svc ./ 

# runtime stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/ldap-svc .
USER 1000:1000
EXPOSE 8080
ENTRYPOINT ["/app/ldap-svc"]
