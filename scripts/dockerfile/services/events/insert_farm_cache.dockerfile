FROM golang:1.24.4-bookworm

WORKDIR /services

COPY services/Events/farm/builds/insert_farm_cache /usr/bin/insert_farm_cache 

EXPOSE 50051

CMD ["insert_farm_cache"]
