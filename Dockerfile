FROM ubuntu:latest
LABEL authors="gentle"

ENTRYPOINT ["top", "-b"]