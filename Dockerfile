FROM scratch
WORKDIR /app
ENV SERVER=":80"
COPY subdomain .
ENTRYPOINT ["/app/subdomain"]
CMD ["-server=${SERVER}"]