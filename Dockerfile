FROM golang:latest AS builder

WORKDIR /webAETZ

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /webAETZ/cmd/web

# Статическая сборка без зависимостей на C-библиотеки
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /webAETZ
# Копируем скомпилированный бинарный файл
COPY --from=builder /webAETZ/cmd/web/main .
COPY --from=builder /webAETZ/urls.txt .
COPY tls /webAETZ/tls
# Даем права на выполнение
RUN chmod +x ./main

 

EXPOSE 4001

CMD ["./main"]
