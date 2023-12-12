FROM golang:1.21

RUN apt-get update && apt-get install -y \
    xorg-dev \
    libgl1-mesa-dev \
    libopenal1 \
    libopenal-dev \
    libvorbis0a \
    libvorbis-dev \
    libvorbisfile3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main

ENV HEADLESS=true

EXPOSE 8080

CMD ["./main"]
