FROM ubuntu:20.10

LABEL org.opencontainers.image.source https://github.com/ovrclk/ismyaccountfucked

EXPOSE 8080

ADD ismyaccountfucked ismyaccountfucked

ENTRYPOINT ["/ismyaccountfucked"]

CMD ["server"]
