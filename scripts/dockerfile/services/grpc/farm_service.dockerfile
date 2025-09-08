FROM golang:1.24.4-bookworm


COPY  services/Grpc/farm/build/farm_service /usr/bin/farm_service

EXPOSE 50051

CMD ["farm_service"]

