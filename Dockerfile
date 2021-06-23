FROM scratch
WORKDIR /app
COPY subdomain .
ENTRYPOINT ["/app/subdomain"]