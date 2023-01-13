FROM golang:1.19-buster

COPY . /app/cross-clipboard
WORKDIR /app/cross-clipboard

# install libx11-dev abd Xvfb
RUN apt update && apt install -y libx11-dev xvfb

RUN go build .

ENTRYPOINT [ "sh", "run.sh" ]
