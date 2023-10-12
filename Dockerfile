FROM ubuntu:latest
RUN apt-get update && apt-get upgrade
RUN apt-get install -y git

# Organization name
ENV PING_PONG_ORG=''
# Bucket name
ENV PING_PONG_BUCKET=''
# Token
ENV PING_PONG_TOKEN=''
# DB URL
ENV PING_PONG_URL=''
# Target URL
ENV PING_PONG_TARGET=''

EXPOSE 8086
COPY ./ping-pong-api /root
CMD /root/ping-pong-api
