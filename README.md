# dadosjusbr.org

> Ao mudar o foco para o sistema de justiça (incluindo MPs, Procuradorias e Defensorias) tivemos que mudar o formato de dados, coletores e o site. Estamos trabalhando árduamente para chegar na versão 1.0, o que deve acontecer no primeiro semestre de 2020.

[![Build Status](https://travis-ci.org/dadosjusbr/remuneracoes.svg?branch=master)](https://travis-ci.org/dadosjusbr/remuneracoes) [![codecov](https://codecov.io/gh/dadosjusbr/remuneracoes/branch/master/graph/badge.svg)](https://codecov.io/gh/dadosjusbr/remuneracoes) [![Go Report Card](https://goreportcard.com/badge/github.com/dadosjusbr/remuneracoes)](https://goreportcard.com/report/github.com/dadosjusbr/remuneracoes)

Para garantir o cumprimento da Lei de Acesso à Informação [(Lei n. 12.527, de 2011)](http://www.planalto.gov.br/ccivil_03/_ato2011-2014/2011/lei/l12527.htm) e observando o decisão do STF no [ARE 652777](http://www.stf.jus.br/portal/jurisprudenciaRepercussao/verAndamentoProcesso.asp?incidente=4121428&numeroProcesso=652777&classeProcesso=ARE&numeroTema=483#), o Conselho Nacional de Justiça (CNJ) publicou a [Resolução n. 151](http://www.cnj.jus.br/busca-atos-adm?documento=2537), que determina a divulgação nominal da remuneração recebida por membros, servidores e colaboradores do Judiciário na Internet. Em 20 de outubro de 2017, o CNJ determinou que as informações fossem encaminhadas por [documento padrão](http://cnj.jus.br/files/conteudo/arquivo/2017/11/becada0200f03cb5a129ce57513f8ff3.xls). As informações apresentadas são de responsabilidades dos tribunais. Na medida em que os dados padronizados são enviados, o CNJ faz a consolidação e a publicação das planilhas [nesta página](http://www.cnj.jus.br/transparencia/remuneracao-dos-magistrados).

Esse modelo de planilha é um arquivo XLS, um formato proprietário do Microsoft Excel, que é composto por cinco abas com os respectivos nomes: contracheque, subsídio, indenizações, vantagens eventuais e dados cadastrais. Tais dados possuem um valor imenso para o controle social das políticas públicas.

No entanto, essas publicações são realizadas mensalmente através de um conjunto de 93 planilhas, uma para cada tribunal, em um formato que não é amigável para ser usado por ferramentas de análise e processamento de dados. O próprio Tribunal de Contas da União, no Processo [TC 017.368/2016-2](https://portal.tcu.gov.br/fiscalizacao-de-tecnologia-da-informacao/atuacao/avaliacao-de-transparencia/), constatou em 2018 que muitos órgãos do judiciário federal "apresentam limitações das ferramentas de pesquisa dos portais dificultam ou impossibilitam a extração de relatórios em formatos abertos, e não são aderentes a requisitos de acessibilidade, bem como não obedecem a padronização na apresentação dos dados e da nomenclatura de termos, dificultando o acesso às informações contidas nos portais".

Pensando nisso, o projeto [dadosjusbr](https://github.com/dadosjusbr) tem como principal objetivo prover acesso às informações de remunerações do sistema judiciário de forma consolidada e em formato aberto. Para tal, criamos sistemas computacionais que realizam a conversão, consolidação e validação dos dados de forma contínua. Com o monitoramento contínuo, podemos cobrar a disponiblização ou correção de informações, caso necessário. Por fim, disponibilizamos o [dadosjusbr.online](https://dadosjusbr.online), um portal onde os dados são publicados em um formato amplamente compatível com ferramentas de análise e processamento de dados e estão organizados em uma página por mês de referência.

### Como rodar a aplicação localmente?

### Configuraçoes necessárias:

[MongoDb](https://docs.mongodb.com/guides/server/install/)  Versão 3.6+  

[GoLang](https://golang.org/doc/install)  Versão 1.14+  

[Node](https://nodejs.org/en/download/) Versão 13.12+  

### Como rodar a cli:
Fazer o download do repositório remuneraçoes:

```console
$ git clone https://github.com/dadosjusbr/remuneracoes.git
```

Dentro do diretório remuneraçoes usar o auxiliador de pacotes para instalar as dependências:
```console
$ cd remuneracoes
$ npm i
```
Renomear o arquivo .env.sample na raiz do projeto para .env e configurar suas variáveis de ambiente
 
### Para rodar o servidor:
```console
$ go run main.go 
$ npm run serve
```
