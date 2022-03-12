FROM golang:latest
COPY ./go/app /app/go

WORKDIR /app/go/handler

RUN rm go.mod

# 初期化
RUN go mod init github.com/Shinya0714/manamana/go/app/handler

RUN go get github.com/sclevine/agouti \
 && go get github.com/labstack/echo/v4 \
 && go get github.com/labstack/echo/v4/middleware \
 && go get github.com/joho/godotenv \
 && go get github.com/lib/pq \
 && go get github.com/jinzhu/gorm

WORKDIR /app/go

RUN rm go.mod

# 初期化
RUN go mod init github.com/Shinya0714/manamana

# モジュール取得
RUN go get github.com/Shinya0714/manamana/go/app/handler
 
# ビルド
RUN go build -o main . 

# unzip取得
RUN apt-get update && apt-get install -y unzip

# chrome取得
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add \
 && wget http://dl.google.com/linux/deb/pool/main/g/google-chrome-unstable/google-chrome-unstable_93.0.4577.18-1_amd64.deb \
 && apt-get install -y -f ./google-chrome-unstable_93.0.4577.18-1_amd64.deb

# chromedriver取得
ADD https://chromedriver.storage.googleapis.com/93.0.4577.15/chromedriver_linux64.zip /usr/local/go/bin/
RUN cd /usr/local/go/bin/ \
 && unzip chromedriver_linux64.zip

CMD ["/app/go/main"]