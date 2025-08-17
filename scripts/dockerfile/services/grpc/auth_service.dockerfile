FROM golang:1.24.4-bookworm

COPY services/Grpc/auth/build/auth_server  /usr/bin/auth_server
# RUN curl -L https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.20/grpc_health_probe-linux-amd64 \
#   -o /usr/local/bin/grpc_health_probe && \
#   chmod +x /usr/local/bin/grpc_health_probe

EXPOSE 50051

CMD ["auth_server"]

