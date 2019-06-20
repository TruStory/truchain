FROM alpine:3.7
LABEL Shane Vitarana <shane@trustory.io>

RUN apk update && \
    apk upgrade && \
    apk --no-cache add curl jq file

VOLUME [ /truchaind ]
WORKDIR /truchaind
EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start", "--log_level", "main:info,state:info,*:error,app:info,argument:info,backing:info,category:info,challenge:info,expiration:info,stake:info,stories:info"]
STOPSIGNAL SIGTERM

COPY wrapper.sh /usr/bin/wrapper.sh

