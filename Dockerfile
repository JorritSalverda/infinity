FROM docker:20.10.7-dind

RUN mkdir -p /work
WORKDIR /work

COPY ./docker-entrypoint.sh /
RUN chmod 500 /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

COPY infinity usr/local/bin