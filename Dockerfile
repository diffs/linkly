FROM golang:1.23

COPY linkly /linkly
WORKDIR /
RUN chmod +x linkly

ENTRYPOINT ["/linkly"]
