FROM golang:latest

COPY ./go/app /app/go

# 切り替え
WORKDIR /app/go/general

RUN rm go.mod

# 初期化
RUN go mod init github.com/Shinya0714/manamana/go/app/general

# モジュール取得
RUN go mod tidy

RUN rm general

# ビルド
RUN go build -o general . 

RUN touch general

RUN pwd

RUN ls -lt

# 切り替え
WORKDIR /app/go/handler

RUN rm go.mod

# 初期化
RUN go mod init github.com/Shinya0714/manamana/go/app/handler

# モジュール取得
RUN go mod tidy

RUN rm handler

# ビルド
RUN go build -o handler . 

RUN touch handler

RUN pwd

RUN ls -lt
                        
# 切り替え
WORKDIR /app/go

RUN rm go.mod

# 初期化
RUN go mod init github.com/Shinya0714/manamana

# モジュール取得
RUN go mod tidy

RUN rm main

# ビルド
RUN go build -o main . 

RUN touch main

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

RUN pwd

RUN ls -lt

CMD ["/app/go/main"]