FROM golang:1.17.0-alpine AS builder

WORKDIR /build
COPY ./go.* ./
RUN go mod download
COPY . .
RUN go build -o api

FROM alpine

# Important when starting API locally.
RUN apk add --no-cache tzdata

# Copia binário a partir do ambiente builder
COPY --from=builder /build/api /

# É necessário expor ao menos uma porta
EXPOSE 8081

# Declara e exporta variáveis de ambiente necessárias.
ENV PORT=$PORT \
    MONGODB_URI=$MONGODB_URI \
    MONGODB_NAME=$MONGODB_NAME \
    MONGODB_MICOL=$MONGODB_MICOL \
    MONGODB_AGCOL=$MONGODB_AGCOL \
    MONGODB_PKGCOL=$MONGODB_PKGCOL \
    MONGODB_REVCOL=$MONGODB_REVCOL \
    DADOSJUSBR_ENV=$DADOSJUSBR_ENV \
    DADOSJUS_URL=$DADOSJUS_URL \
    PACKAGE_REPO_URL=$PACKAGE_REPO_URL \
    SEARCH_LIMIT=$SEARCH_LIMIT \
    DOWNLOAD_LIMIT=$DOWNLOAD_LIMIT \
    PG_PORT=$PG_PORT \
    PG_HOST=$PG_HOST \
    PG_DATABASE=$PG_DATABASE \
    PG_USER=$PG_USER \
    PG_PASSWORD=$PG_PASSWORD \
    NEWRELIC_APP_NAME=$NEWRELIC_APP_NAME \
    NEWRELIC_LICENSE=$NEWRELIC_LICENSE 

# Inicia a API
CMD ["/api"]