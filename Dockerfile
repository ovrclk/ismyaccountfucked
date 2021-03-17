FROM ghcr.io/ovrclk/akash:0.10.1
LABEL org.opencontainers.image.source https://github.com/ovrclk/ismyaccountfucked

EXPOSE 8080

CMD /ismyaccountfucked server
