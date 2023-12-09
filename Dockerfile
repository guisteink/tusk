# Use a imagem oficial do Golang como base
FROM golang:latest

# Define o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos de configuração para o diretório de trabalho
COPY go.mod go.sum ./

# Baixa as dependências do módulo
RUN go mod download

# Copia o código-fonte para o diretório de trabalho
COPY . .

# Compila o código Go
RUN go build -o ./app ./cmd/main.go

# Expõe a porta que o aplicativo estará ouvindo
EXPOSE 3000

# Comando para iniciar o aplicativo após o contêiner ser iniciado
CMD ["./scripts/start-app.sh"]
