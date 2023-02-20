# Api dadosjusbr.org

[![Go Report Card](https://goreportcard.com/badge/github.com/dadosjusbr/remuneracoes)](https://goreportcard.com/report/github.com/dadosjusbr/remuneracoes)

A Lei de Acesso à Informação [(Lei n. 12.527, de 2011)](http://www.planalto.gov.br/ccivil_03/_ato2011-2014/2011/lei/l12527.htm), regula a obrigatoriedade da disponibilização na internet dos dados de gastos público, porém esses dados não são padronizados e cada órgão tem sua própria formatação, podendo ser encontrado em diversas nomenclaturas e tipos diferentes de arquivos (pdf, html, planilhas eletrônicas, json e etc). Por esse motivo, esses arquivos não possuem um formato amigável para ser usado por ferramentas de análise e processamento de dados.

Pensando nisso, o projeto [dadosjusbr](https://github.com/dadosjusbr) tem como principal objetivo prover acesso às informações de remunerações do sistema judiciário de forma consolidada e em formato aberto. Para tal, utilizamos do framework Nextjs para criar as interfaces do usuários e alimentamos essas interfaces com um servidor ambientado em GoLang.

Com essas tecnologias como base, criamos sistemas computacionais que realizam a coleta, conversão, consolidação e validação dos dados de forma contínua. O DadosJusBr é conectado aos repositórios [coletores](https://github.com/orgs/dadosjusbr/repositories?q=coletor), que são responsáveis por adquirir os dados dos órgãos e padronizá-los. Já o repositório de [storage](https://github.com/dadosjusbr/storage), é responsável pelo armazenamento desses dados coletados.

Com o monitoramento contínuo, podemos cobrar a disponiblização ou correção de informações, caso necessário. Por fim, disponibilizamos o [DadosJusBr](https://dadosjusbr.org/), um portal onde os dados são publicados em um formato amplamente compatível com ferramentas de análise e processamento de dados e estão organizados em uma página por mês de referência. Mais informações [aqui.](https://dadosjusbr.org/#/sobre)

Esse projeto foi elaborado com o intuito de praticar a cidadania e tornar os dados mais acessíveis para o cidadão. Você cidadão/empresa pode fazer parte dessa jornada conosco, quer saber como?

- Informe se há alguma inconsistência ou erros na api.
- Atue como fiscal e cobre dos órgãos sobre a disponibilidade dos dados à população.
- Sugira novos órgãos para elaboração de robôs, se tiver conhecimento, desenvolva um.
- Sugerir coisas interessantes que você acha que irão contribuir para o projeto!

## Como rodar a aplicação localmente?

#

### Caso você não tenha [Docker](https://www.docker.com/get-started/) e o [Docker compose](https://docs.docker.com/compose/install/) instalados na sua máquina, é necessário as seguintes dependências:

- [GoLang](https://golang.org/doc/install) Versão 1.18+

- [Postgresql](https://www.postgresql.org/download/) Versão 14.4+

#

### Fazer o download do repositório remuneraçoes:

```console
$ git clone https://github.com/dadosjusbr/api.git
```

### Após a instalação, Renomear o arquivo `.env.example` na raiz do projeto para `.env` e configurar suas variáveis de ambiente:

| Variável              | Descrição                                                                                                                    | Exemplo                         |
| --------------------- | ---------------------------------------------------------------------------------------------------------------------------- | ------------------------------- |
| PORT                  | Porta que servirá a API                                                                                                      | 8081                            |
| DADOSJUSBR_ENV        | O ambiente a ser executado                                                                                                   | 'Development' ou 'Production'   |
| DADOSJUS_URL          | URI utilizada para mapeamento dos arquivos para download para o site do DadosJusBr                                           | https://dadosjusbr.org/download |
| PACKAGE_REPO_URL      | URI utilizada para mapeamento dos arquivos para download para o repositório de arquivos AWS S3                               | https://example.amazonaws.com   |
| SEARCH_LIMIT          | Número limite de dados que a rota de pesquisa irá trazer                                                                     | 100                             |
| DOWNLOAD_LIMIT        | Número limite de dados que a rota de download irá baixar                                                                     | 10000                           |
| PG_DATABASE           | Nome do banco de dados postgres                                                                                              | dadosjusbr                      |
| PG_USER               | Nome do usuário do banco de dados postgres                                                                                   | dadosjusbr                      |
| PG_PORT               | Porta de conexão com o banco de dados postgres                                                                               | 5432                            |
| PG_HOST               | Host do banco de dados postgres                                                                                              | localhost                       |
| PG_PASSWORD           | Senha do banco de dados postgres                                                                                             | dadosjusbr                      |
| NEWRELIC_APP_NAME     | Nome do app New Relic                                                                                                        |                                 |
| NEWRELIC_LICENSE      | Licensa New Relic                                                                                                            |                                 |
| AWS_S3_BUCKET         | Nome do bucket localizado no AWS S3                                                                                          | dadosjusbr                      |
| AWS_REGION            | Região da AWS onde o bucket do AWS S3 está localizado                                                                        | us-east-1                       |
| AWS_ACCESS_KEY_ID     | Chave de acesso da aws, essa variável de ambiente é utilizada quando queremos rodar a aplicação utilizando elastic beanstalk |                                 |
| AWS_SECRET_ACCESS_KEY | Chave secreta da aws, essa variável de ambiente é utilizada quando queremos rodar a aplicação utilizando elastic beanstalk   |

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

## Rodando o servidor utilizando o elastic beanstalk cli

Primeiramente, é necessário realizar a [instalação do eb cli](https://docs.aws.amazon.com/pt_br/elasticbeanstalk/latest/dg/eb-cli3-install.html)

Após ter realizado a instalação, configure o eb cli para que ele utilize as variáveis de ambiente localizadas no .env

```console
eb local setenv `cat .env | sed '/^#/ d' | sed '?^$/ d'`
```

Lembrando que é necessário configurar as variáveis de ambiente AWS_ACCESS_KEY_ID e AWS_SECRET_ACCESS_KEY no .env

Após estar tudo configurado, basta executar:

```console
eb local run
```

## Testando servidor

Caso a execução tenha sido realizada com sucesso, você pode utilizar o seu cliente de api REST para acessar o servidor local, que está localizado em http://{HOST}:{PORT}/v1/orgaos

## Documentando as rotas da API utilizando o swagger

O swagger é uma ferramenta que ajuda no processo de documentar rotas de API's. Na API do DadosJusBr, utilizamos a biblioteca [swaggo](https://github.com/swaggo/swag) para criar as documentações. Com essa biblioteca, basta que a gente adicione comentários no nosso código e a documentação será gerada.

Para utilizar o swaggo, é necessário primeiramente ter o binário instalado. Para instala-lo, basta executar o seguinte comando:

```console
$ go install github.com/swaggo/swag/cmd/swag@latest
```

Com o binário instalado, podemos iniciar a criação da documentação.

Como exemplo iremos documentar a rota da api `/v1/orgao/:orgao`. Os passos são:

### 1 - Ir até o método chamado pela rota(nesse caso será o `GetAgencyById`)

```
func (h handler) GetAgencyById(c echo.Context) error {
  ...
}
```

### 2 - Adicionar, logo acima da declaração do método, os campos que a documentação da rota terá. O swaggo cria a documentação através dos comentários seguidos de anotações.

```
//	@ID				GetAgencyById
//	@Description	Busca um órgão específico utilizando seu ID.
//	@Produce		json
//	@Param			orgao				path		string	true	"ID do órgão. Exemplos: tjal, tjba, mppb."
//	@Success		200					{object}	agency	"Requisição bem sucedida."
//	@Failure		404					{string}	string	"Órgão não encontrado."
//	@Router			/v1/orgao/{orgao} 	[get]
func (h handler) GetAgencyById(c echo.Context) error {
...
}
```

O swaggo possui diversos atributos que podem ser adicionados na documentação, você pode ver quais são esses atributos na [documentação do swaggo](https://github.com/swaggo/swag#readme).

### 3 - Com os atributos adicionados, basta executar os seguintes comandos, para formatar o código e gerar a documentação:

```
$ swag fmt
$ swag init
```

Isso fará com que a pasta `docs` seja modificada. `Não modifique ela manualmente, pois ela é autogenerada pelo swaggo.`

### 4- Com isso, basta rodar a api e acessar a rota /swagger/index.html, lá estará toda a documentação criada.
