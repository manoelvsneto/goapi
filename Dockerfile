# Usar uma imagem base oficial do Golang
#FROM golang:1.22-alpine AS builder
FROM golang:1.22.3-alpine AS builder
# Definir o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copiar os arquivos go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixar as dependências
RUN go mod download

# Copiar o restante dos arquivos da aplicação para o diretório de trabalho
COPY . .

# Compilar a aplicação
RUN go build -o main .

RUN rm -f configmap.yaml
RUN rm -f deployment.yaml
RUN rm -f service.yaml
RUN rm -f ingress.yaml


# Usar uma imagem base menor para rodar a aplicação
FROM alpine:latest

# Definir o diretório de trabalho
WORKDIR /root/

# Copiar o binário compilado da imagem de build
COPY --from=builder /app/main .

# Definir a porta que a aplicação vai expor
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]
