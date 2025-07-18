# Build stage
FROM golang:1.24-alpine3.22 AS builder

# Instalar dependências necessárias e atualizar pacotes
RUN apk update && apk add --no-cache git ca-certificates tzdata

# Definir diretório de trabalho
WORKDIR /app

# Copiar go.mod e go.sum primeiro (para cache das dependências)
COPY go.mod go.sum ./

# Baixar dependências
RUN go mod download

# Copiar código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o main .

# Runtime stage
FROM alpine:3.19

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Criar usuário não-root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Definir diretório de trabalho
WORKDIR /app

# Copiar binário da stage de build
COPY --from=builder /app/main .

# Definir permissões
RUN chown -R appuser:appgroup /app
USER appuser

# Expor porta
EXPOSE 3000

# Comando para executar a aplicação
CMD ["./main"]