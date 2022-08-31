# Api dadosjusbr.org

> Ao mudar o foco para o sistema de justiça (incluindo MPs, Procuradorias e Defensorias) tivemos que mudar o formato de dados, coletores e o site. Estamos trabalhando árduamente para chegar na versão 1.0, o que deve acontecer no primeiro semestre de 2020.

[![Go Report Card](https://goreportcard.com/badge/github.com/dadosjusbr/remuneracoes)](https://goreportcard.com/report/github.com/dadosjusbr/remuneracoes)

A Lei de Acesso à Informação [(Lei n. 12.527, de 2011)](http://www.planalto.gov.br/ccivil_03/_ato2011-2014/2011/lei/l12527.htm), regula a obrigatoriedade da disponibilização na internet dos dados de gastos público, porém esses dados não são padronizados e cada órgão tem sua própria formatação, podendo ser encontrado em diversas nomenclaturas e tipos diferentes de arquivos (pdf, html, planilhas eletrônicas, json e etc). Por esse motivo, esses arquivos não possuem um formato amigável para ser usado por ferramentas de análise e processamento de dados.

Pensando nisso, o projeto [dadosjusbr](https://github.com/dadosjusbr) tem como principal objetivo prover acesso às informações de remunerações do sistema judiciário de forma consolidada e em formato aberto. Para tal, utilizamos do framework Nextjs para criar as interfaces do usuários e alimentamos essas interfaces com um servidor ambientado em GoLang.

Com essas tecnologias como base, criamos sistemas computacionais que realizam a coleta, conversão, consolidação e validação dos dados de forma contínua. O DadosJusBr é conectado ao repositório de [coleta](https://github.com/dadosjusbr/coletores), que é responsável por adquirir os dados dos órgãos e padronizá-los. Já o repositório de [storage](https://github.com/dadosjusbr/storage), é responsável pelo armazenamento desses dados coletados.

Com o monitoramento contínuo, podemos cobrar a disponiblização ou correção de informações, caso necessário. Por fim, disponibilizamos o [DadosJusBr](https://dadosjusbr.org/), um portal onde os dados são publicados em um formato amplamente compatível com ferramentas de análise e processamento de dados e estão organizados em uma página por mês de referência. Mais informações [aqui.](https://dadosjusbr.org/#/sobre)

Esse projeto foi elaborado com o intuito de praticar a cidadania e tornar os dados mais acessíveis para o cidadão. Você cidadão/empresa pode fazer parte dessa jornada conosco, quer saber como?

- Informe se há alguma inconsistência ou erros na api.
- Atue como fiscal e cobre dos órgãos sobre a disponibilidade dos dados à população.
- Sugira novos órgãos para elaboração de robôs, se tiver conhecimento, desenvolva um.
- Sugerir coisas interessantes que você acha que irão contribuir para o projeto!

## Como rodar a aplicação localmente?

#

### Caso você não tenha [Docker](https://www.docker.com/get-started/) e o [Docker compose](https://docs.docker.com/compose/install/) instalados na sua máquina, é necessário as seguintes dependências:

- [MongoDb](https://docs.mongodb.com/guides/server/install/) Versão 3.6+

- [GoLang](https://golang.org/doc/install) Versão 1.14+

- [Postgresql](https://www.postgresql.org/download/) Versão 14.4+

#

### Fazer o download do repositório remuneraçoes:

```console
$ git clone https://github.com/dadosjusbr/api.git
```

### Após a instalação, Renomear o arquivo `.env.example` na raiz do projeto para `.env` e configurar suas variáveis de ambiente:

| Variável          | Descrição                                                                                     |
| ----------------- | --------------------------------------------------------------------------------------------- | -------------------- |
| API_PORT          | Porta que servirá a API                                                                       |
| MONGODB_URI       | URI de conexão com o mongobd                                                                  |
| MONGODB_NAME      | Nome do banco de dados mongodb                                                                |
| MONGODB_MICOL     | Nome da coleção de **informações de remunerações mensais**                                    |
| MONGODB_AGCOL     | Nome da coleção de **órgãos**                                                                 |
| MONGODB_PKGCOL    | Nome da coleção de **arquivos coletados**                                                     |
| DADOSJUSBR_ENV    | `Development                                                                                  | Production` Ambiente |
| DADOSJUS_URL      | URI utilizada para mapeamento dos arquivos para download para o site do dados jus             |
| PACKAGE_REPO_URL  | URI utilizada para mapeamento dos arquivos para download para o repositório de arquivos swift |
| SEARCH_LIMIT      | Número limite de dados que a rota de pesquisa irá trazer                                      |
| DOWNLOAD_LIMIT    | Número limite de dados que a rota de download irá baixar                                      |
| PG_PORT           | Porta de conexão com o banco de dados postgres                                                |
| PG_HOST           | Host do banco de dados postgres                                                               |
| PG_DATABASE       | Nome do banco de dados postgres postgres                                                      |
| NEWRELIC_APP_NAME | Nome do app New Relic                                                                         |
| NEWRELIC_LICENSE  | Licensa New Relic                                                                             |

> ## Atenção
>
> Caso você utilize o Docker para executar o script, saiba que o host do banco de dados postgres precisa ser o endereço IP da sua máquina, pois o Docker não consegue acessar o banco utilizando "localhost" como host local.
> Para saber o endereço IP da sua máquina, rode o seguinte comando no terminal:

```console
$ hostname -I | cut -d" " -f1
```

Será exibido algo parecido com isso:

> 123.456.789.10%

Copie tudo que estiver antes do "%" e cole na variável de ambiente PG_HOST

## Rodando o servidor sem Docker

Para rodar o servidor sem utilizar o Docker, execute o seguinte comando no terminal:

```console
$ go run .
```

## Rodando o servidor com Docker

Para rodar o servidor utilizando o Docker, execute o seguinte comando no terminal:

```console
$ docker-compose -f docker-compose-dev.yml up -d
```

Para exibir os logs do docker e saber se ocorreu tudo bem:

```console
$ docker-compose logs
```

Para parar o container

```console
$ docker-compose -f docker-compose-dev.yml down
```

## Testando servidor

Caso a execução tenha sido realizada com sucesso, você pode utilizar o seu cliente de api REST para acessar o servidor local, que está localizado em http://{HOST}:{API_PORT}/v1/orgaos

