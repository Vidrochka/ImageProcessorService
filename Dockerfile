FROM golang

ADD . /go/src/github.com/anthill-com/ImageProcessorService

RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/nfnt/resize
RUN go get github.com/pelletier/go-toml

RUN go install github.com/anthill-com/ImageProcessorService/ImageProcessorService...

COPY ./ImageProcessorService/config.toml /go/bin/

WORKDIR /go/bin/

ENTRYPOINT ./ImageProcessorService