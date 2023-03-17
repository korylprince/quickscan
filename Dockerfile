FROM ubuntu:latest

ARG GO_PROJECT_NAME
ENV GO_PROJECT_NAME=${GO_PROJECT_NAME}

COPY docker-entrypoint.sh /
COPY ${GO_PROJECT_NAME} /

CMD ["/docker-entrypoint.sh"]
