FROM alpine:3.22.1

RUN apk --no-cache add ca-certificates

COPY  services/Rest/farm_gateway/build/farm_gateway /bin/farm_gateway
EXPOSE 50051
CMD ["/bin/farm_gateway"]
