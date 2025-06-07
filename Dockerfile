FROM golang:1.24-alpine3.21

WORKDIR /gocov-threshold
COPY . /gocov-threshold

RUN go mod download

RUN cd cmd && go install

CMD ["gocov-threshold"]