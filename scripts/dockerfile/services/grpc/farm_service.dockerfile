FROM alpine:3.22.1

RUN apk --no-cache add ca-certificates

COPY  services/Grpc/farm/build/farm_service /bin/farm_service
EXPOSE 50051
CMD ["/bin/farm_service"]

