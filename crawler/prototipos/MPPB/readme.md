# Ministério Público da Paraíba - Crawler
Este crawler tem como objetivo a recuperação de informações sobre folhas de pagamentos dos funcionários do Ministério Público da Paraíba. O site com as informações pode ser acessado [aqui](http://pitagoras.mppb.mp.br/PTMP/FolhaListar) 
O crawler está estruturado como uma CLI. Você passa dois argumentos (mês e ano) e serão baixadas seis planilhas no formato ODS, cada planilha é referente a uma destas categorias: Membros Ativos, Membros Inativos, Servidores Ativos, Servidores Inativos, Servidores à Disposição e Aposentados/Pensionistas
.

### Como usar?
- É preciso ter o compilador de Go instalado em sua máquina!
- No diretório [**remuneracoes/crawler/prototipos/MPPB**](https://github.com/dadosjusbr/remuneracoes/tree/primeiros-crawlers/crawler/prototipos/MPPB) você encontrará o arquivo **crawler-mppb.go**, nesse arquivo está o código do crawler.
- Rode o comando abaixo, com mês e ano que você quer
#### go run crawler-mppb.go --mes=?? --ano=????

### Sobre os dados

As planilhas possuem as seguintes colunas:

- **Nome**: Nome completo do funcionário (String)
- **Matrícula**: Matrícula do funcionário (String)  
- **Cargo**: Cargo do funcionário dentro do MP (String)
- **Lotação**: Local (cidade, departamento, promotoria) em que o funcionário trabalha (String)
- **Remuneração do cargo efetivo**: Vencimento, GAMPU, V.P.I, Adicionais de Qualificação, G.A.E e G.A.S, além de outras desta natureza. Soma de todas essas remunerações. (Number) 
- **Outras Verbas Remuneratórias, Legais ou Judiciais**: V.P.N.I., Adicional por tempo de serviço, quintos, décimos e vantagens decorrentes de sentença judicial ou extensão administrativa. (Number) 
- **Função de Confiança ou Cargo em Comissão**: Rubricas que representam a retribuição paga pelo exercício de função (servidor efetivo) ou remuneração de cargo em comissão (servidor sem vínculo ou requisitado). (Number) 
- **Gratificação Natalina**: Parcelas da Gratificação Natalina (13º) pagas no mês corrente, ou no caso de vacância ou exoneração do servidor. (Number)  
- **Férias (⅓ Constitucional)**: Adicional correspondente a 1/3 (um terço) da remuneração, pago ao servidor por ocasião das férias. (Number) 
- **Abono de Permanência**:  Valor equivalente ao da contribuição previdenciária, devido ao funcionário público que esteja em condição de aposentar-se, mas que optou por continuar em atividade (instituído pela Emenda Constitucional nº 41, de 16 de dezembro de 2003. (Number) 
- **Total de Rendimentos Brutos**: Total dos rendimentos brutos pagos no mês. (Number) 
- **Contribuição Previdenciária**: Contribuição Previdenciária Oficial (Plano de Seguridade Social do Servidor Público e Regime Geral de Previdência Social). (Number) 
- **Imposto de Renda**: Imposto de Renda Retido na Fonte. (Number) 
- **Retenção por Teto Constitucional**: Valor deduzido da remuneração básica bruta, quando esta ultrapassa o teto constitucional, nos termos da legislação correspondente.
- **Total de Descontos**:  Total dos descontos efetuados no mês. (Number) 
- **Rendimento Líquido Total**: Rendimento líquido após os descontos referidos nos itens anteriores. (Number) 
- **Indenizações**: Auxílio-alimentação, Auxílio-transporte, Auxílio-Moradia, Ajuda de Custo e outras dessa natureza, exceto diárias, que serão divulgadas no Portal da Transparência. Soma de todas essas remunerações. (Number) 
- **Outras Remunerações Temporárias**: Valores pagos a título de Adicional de Insalubridade ou de Periculosidade, Adicional Noturno, Serviço Extraordinário, Substituição de Função, 'Atrasados'. Soma de todas essas remunerações. (Number) 






