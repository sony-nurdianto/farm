FROM golang:1.24.4-bookworm

COPY services/Events/farmer/builds/users_mutations  /usr/bin/users_mutations

EXPOSE 50051

CMD ["users_mutations"]

