FROM alpine:latest

COPY linkly /linkly
WORKDIR /
RUN chmod +x linkly

ENTRYPOINT ["/linkly"]
