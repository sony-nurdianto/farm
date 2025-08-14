FROM golang:1.24.4-bookworm

COPY services/Grpc/auth/build/auth_server  /usr/bin/auth_server

EXPOSE 50051

CMD ["auth_server"]

