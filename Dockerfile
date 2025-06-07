FROM golang:1.24-alpine3.21

WORKDIR /gocov-threshold
COPY . /gocov-threshold


RUN go mod download
RUN go mod tidy

RUN cd cmd/threshold && go install

CMD ["threshold"]