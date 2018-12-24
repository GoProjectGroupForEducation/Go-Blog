FROM golang:1.8
MAINTAINER WilsonZhong "1316628630@qq.com"
WORKDIR $GOPATH/src/github.com/GoProjectGroupForEducation/Go-Blog
ADD . $GOPATH/src/github.com/GoProjectGroupForEducation/Go-Blog
RUN go get -v github.com/GoProjectGroupForEducation/Go-Blog
EXPOSE 8081
ENTRYPOINT ["Go-Blog"]