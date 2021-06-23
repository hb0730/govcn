FROM golang:1.16-alpine
WORKDIR /app
COPY subdomain .
ENTRYPOINT ["/app/subdomain"]