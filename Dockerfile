FROM golang:1.20

RUN git clone https://github.com/blablatov/restduples.git
WORKDIR restduples

RUN go mod download

COPY *.go ./
COPY *.conf ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /restduples

EXPOSE 12345

CMD ["/restduples"]