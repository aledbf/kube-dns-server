FROM alpine:3.3

RUN apk add -U bash curl nsd unbound && \
  rm -rf /var/cache/apk/*
  
COPY kube-dns-server /kube-dns-server

ENTRYPOINT ["/kube-dns-server"]
