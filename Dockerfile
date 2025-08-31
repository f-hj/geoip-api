# Compile Golang
FROM golang as builder

COPY . /data
WORKDIR /data
RUN CGO_ENABLED=0 go build

# Download maxmind db
ARG MAXMIND_LICENSE

RUN wget "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&suffix=tar.gz&license_key=${MAXMIND_LICENSE}" -O GeoLite2-City.tar.gz && tar -xvf GeoLite2-City.tar.gz  --strip-components 1 --wildcards */GeoLite2-City.mmdb
RUN wget "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-ASN&suffix=tar.gz&license_key=${MAXMIND_LICENSE}" -O GeoLite2-ASN.tar.gz && tar -xvf GeoLite2-ASN.tar.gz  --strip-components 1 --wildcards */GeoLite2-ASN.mmdb

# Get alpine
FROM alpine:3.13.0
RUN apk add --no-cache tzdata

COPY --from=builder /data/geoip-api /bin/geoip-api

RUN mkdir -p /data
COPY --from=builder /data/GeoLite2-City.mmdb /data/GeoLite2-City.mmdb
COPY --from=builder /data/GeoLite2-ASN.mmdb /data/GeoLite2-ASN.mmdb

ENTRYPOINT [ "/bin/geoip-api" ]