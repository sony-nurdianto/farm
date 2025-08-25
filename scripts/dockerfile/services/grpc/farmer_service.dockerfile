FROM golang:1.24.4-bookworm

COPY services/Grpc/farmer/build/farmer_service  /usr/bin/farmer_service

EXPOSE 50051

CMD ["farmer_service"]

