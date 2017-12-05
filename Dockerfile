FROM convox/golang

WORKDIR /go/src/github.com/convox/dmz
COPY . .
RUN go install .

CMD ["dmz"]
