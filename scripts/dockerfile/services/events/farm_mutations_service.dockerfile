FROM golang:1.24.4-bookworm

WORKDIR /services

COPY services/Events/farm/builds/farm_mutations /services/farm_mutations 

COPY services/Events/farm/state .


EXPOSE 50051

CMD ["./farm_mutations"]

