FROM golang:1.20.1

## UPDATE THE OS
RUN apt-get update && \
    apt-get install -y tzdata 

WORKDIR /go/src

## SET ENVIRONMENT
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV TZ America/Sao_Paulo

## START A PROJECT
RUN go mod init github.com/renatospaka/cotacao-dolar

# ## COPY NECESSARY FILES
# COPY go.* ./

## INSTALL MY STANDARD LIBRARIES 
RUN go get github.com/satori/go.uuid && \
    go get github.com/stretchr/testify

# ## TIDY THE PROJECT
RUN go mod download && \
    go mod tidy

## KEEP THE CONTAINER RUNNiNG
CMD ["tail", "-f", "/dev/null"]