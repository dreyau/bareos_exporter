FROM golang as builder
RUN go get -d -v github.com/dreyau/bareos_exporter
WORKDIR /go/src/github.com/dreyau/bareos_exporter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bareos_exporter .

FROM busybox:latest

ENV mysql_port 3306
ENV mysql_server 192.168.3.70
ENV mysql_username monty

WORKDIR /bareos_exporter
COPY --from=builder /go/src/github.com/dreyau/bareos_exporter/bareos_exporter bareos_exporter

CMD ./bareos_exporter -p pw/auth -u $mysql_username -h $mysql_server -P $mysql_port
EXPOSE 8080
