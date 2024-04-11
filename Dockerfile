FROM golang:1.20

WORKDIR app

COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum
COPY example.csv example.csv
COPY template template
RUN go install github.com/golang/mock/mockgen@v1.6.0
RUN go generate ./...

RUN go mod tidy

RUN cd cmd && go build main.go && cp main ./../
RUN cd -
RUN chmod +x ./main
CMD ["./main", ""]
