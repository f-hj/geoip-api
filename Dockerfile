# Compile Golang
FROM golang as builder

COPY . /data
WORKDIR /data
RUN CGO_ENABLED=0 go build

# Download maxmind db
RUN wget https://raw.githubusercontent.com/gonet2/geoip/master/GeoIP2-City.mmdb -O /data/GeoIP2-City.mmdb

# Get alpine
FROM alpine:3.13.0

COPY --from=builder /data/geoip-api /bin/geoip-api

RUN mkdir -p /data
COPY --from=builder /data/GeoIP2-City.mmdb /data/GeoIP2-City.mmdb

ENTRYPOINT [ "/bin/geoip-api" ]