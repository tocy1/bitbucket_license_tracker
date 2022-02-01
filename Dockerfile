FROM golang:1.16-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN go build -o  bb-license-tracker ./src
CMD ["./bb-license-tracker"]
