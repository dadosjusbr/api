# DadosJusBr
Nesse repositório se encontra todos os arquivos e diretórios que são responsáveis pela renderização do nosso site.
## Tecnologias - o que utilizamos?

- Vue
- Bootstrap
- Material Design
- ApexChart
- Node

## Como rodar localmente? (modo desenvolvimento)

**obs**: é preciso ter o node instalado em sua máquina, saiba mais [aqui](https://nodejs.org/en/download/).
No diretório ```./``` execute os seguintes comandos:
Irá instalar todas as dependências necessárias:

```
npm install
```
Irá executar a aplicação que ficará disponível em: http://localhost:8080/.
```
npm run serve
```

## Usando a API

A aplicação é alimentada por uma API que é estrutura no arquivo ```./main.go```. Para usar esta API é necessário um arquivo .env com as credenciais do BD (veja o arquivo ```./.env.example```).
Com as credênciais preenchidas rode o comando:

```
go run main.go
```

A API estará disponível em: http://localhost:**PORT**/. (o valor de PORT é indica no arquivo .env)

Você também pode acessar essa API que está em produção e online através do link: https://dadosjusbr.org/uiapi/v1.

**Para ter acesso aos endpoint disponíveis veja a documentação disponível em ```./main.go```.**

##### Utilizando a API na aplicação:

Para que a aplicação consuma dados da API é preciso indicar sua URL no arquivo ```./ui/.env.development```. A URL default é http://localhost:8081/uiapi/v1, mas pode ser alterado por https://dadosjusbr.org/uiapi/v1 caso você não tenha acesso às credenciais do BD.

## Rodar modo de produção

Para rodar em modo de produção basta executar o comando:

```
npm run serve
```

Em seguida os arquivos compilados da aplicação estarão disponíveis no diretório ```./dist```.

## Use o lint

Para manter a qualidade do código este projeto utilizado o lint.
Rode o comando para ter acesso a um relatório sobre a qualidade e legibilidade do código.

```
npm run lint
```