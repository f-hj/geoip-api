# Compile Golang
FROM golang as builder

COPY . /data
WORKDIR /data
RUN CGO_ENABLED=0 go build

# Get alpine
FROM alpine

COPY --from=builder /data/geoip-api /bin/geoip-api

# Download maxmind db
RUN mkdir -p /data
RUN wget https://raw.githubusercontent.com/gonet2/geoip/master/GeoIP2-City.mmdb -O /data/GeoIP2-City.mmdb

ENTRYPOINT [ "/bin/geoip-api" ]