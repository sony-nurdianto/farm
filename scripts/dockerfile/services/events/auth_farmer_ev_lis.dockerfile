FROM golang:1.24.4-bookworm

COPY services/Events/auth/builds/auth_farmer_ev_lis  /usr/bin/auth_farmer_ev_lis
EXPOSE 50051

CMD ["auth_farmer_ev_lis"]

