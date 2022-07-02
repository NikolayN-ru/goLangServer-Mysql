FROM golang:1.12.6-stretch
WORKDIR /app
COPY . /app
# RUN go get -u github.com/go-sql-driver/mysql
# RUN go get -u github.com/gorilla/mux
RUN go mod vendor
# RUN go get .
RUN go build -o main .
CMD ["/app/main"]

# FROM golang:latest 
# RUN mkdir /app 
# ADD . /app/ 
# WORKDIR /app 
# RUN go build -o main . 
# CMD ["/app/main"]