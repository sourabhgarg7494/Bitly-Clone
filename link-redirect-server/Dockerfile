FROM golang:latest 
EXPOSE 4000
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
ENV GOPATH /app
RUN cd /app ; go install api
CMD ["/app/bin/api"]
