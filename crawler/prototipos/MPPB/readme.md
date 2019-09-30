# Ministério Público da Paraíba - Crawler
Este crawler tem como objetivo a recuperação de informações sobre folhas de pagamentos dos funcionários do Ministério Público da Paraíba. O site com as informações pode ser acessado [aqui](http://pitagoras.mppb.mp.br/PTMP/FolhaListar) 
O crawler está estruturado como uma CLI. Você passa dois argumentos (mês e ano) e serão baixadas seis planilhas no formato ODS.

### Como usar?
- É preciso ter o compilador de Go instalado em sua máquina!
- No diretório [**remuneracoes/crawler/prototipos/MPPB**](https://github.com/dadosjusbr/remuneracoes/tree/primeiros-crawlers/crawler/prototipos/MPPB) você encontrará o arquivo **crawler-mppb.go**, nesse arquivo está o código do crawler.
- Rode o comando abaixo, com mês e ano que você quer
#### go run crawler-mppb.go --mes=?? --ano=????
