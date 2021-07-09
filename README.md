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

### Configuraçoes necessárias:

[MongoDb](https://docs.mongodb.com/guides/server/install/) Versão 3.6+

[GoLang](https://golang.org/doc/install) Versão 1.14+

[Node](https://nodejs.org/en/download/) Versão 13.12+

### Para rodar o servidor:

Fazer o download do repositório remuneraçoes:

```console
$ git clone https://github.com/dadosjusbr/api.git
```

Após a instalação, Renomear o arquivo `.env.example` na raiz do projeto para `.env` e configurar suas variáveis de ambiente.  
Agora a aplicação está pronta para ter um servidor local funcionando, para isso, fazemos:

```console
$ go run main.go
```
