# build stage
FROM golang AS build-env
ADD . /src
RUN cd /src && GO111MODULE=on go build -o /daemon cmd/cmd.go

# final stage
FROM ubuntu
WORKDIR /app
RUN apt update && apt install -y ca-certificates
COPY --from=build-env /daemon /daemon
CMD ['/daemon']
