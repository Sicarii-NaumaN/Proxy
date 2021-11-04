FROM ubuntu:latest
USER root

ENV TZ=Europe/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update && apt-get -y install software-properties-common
RUN apt-get update && add-apt-repository -y ppa:longsleep/golang-backports
RUN apt -y install golang
#RUN export PATH=$PATH:/opt/go/bin
RUN go version

WORKDIR /
COPY . ./

RUN go mod tidy

EXPOSE 8080 8081 8082

CMD go run .