# Tribunal de Justiça da Paraíba

Este crawler tem como objetivo a recuperação de informações sobre folhas de pagamentos dos funcionários do Tribunal de Justiça da Paraíba. O site com as informações pode ser acessado [aqui](https://www.tjpb.jus.br/transparencia/gestao-de-pessoas/folha-de-pagamento-de-pessoal).

O crawler está estruturado como uma CLI. Você passa dois argumentos (mês e ano) e serão baixados dois arquivo no formato **PDF** representando a folha de pagamento da instituição, sendo um referente a magistrados e outro referente a servidores.

## Legislação
Os dados devem estar de acordo com a [Resolução 102 do CNJ](https://atos.cnj.jus.br/atos/detalhar/69).

## Como usar

- É preciso ter o compilador de Go instalado em sua máquina. Mais informações [aqui](https://golang.org/dl/).
- Rode o comando abaixo, com mês e ano que você quer ter acesso as informações

```sh
cd crawler/tjpb
go run crawler-tjpb.go --mes=${MES} --ano=${ANO}
```

## Dicionário de Dados

As planilhas apresentadas nos pdfs contém as seguintes colunas:

### Magistrados: 

- **Magistrado(a) (String)**: Nome completo do funcionário
- **Lotação (String)**: Local (cidade, departamento, promotoria) em que o funcionário trabalha
- **Cargo (String)**: Cargo do funcionário dentro do MP
- **Rendimentos:**
	- **Subsídio (Number)**: Rubricas que representam a retribuição paga pelo exercício de função (servidor efetivo) ou remuneração de cargo em comissão (servidor sem vínculo ou requisitado)
	- **Abono Permanência Previdenciário**: Reembolso da contribuição previdenciária, devido ao Servidor público em regime contratual estatutário que esteja em condição de aposentar-se, mas que optou por continuar em atividade.
	- **Vantagens Pessoais (Number)**: V.P.N.I., Adicional por tempo de serviço, quintos, décimos e vantagens decorrentes de sentença judicial ou extensão administrativa, abono de permanência.
	- **Grat. Exercício Mesa Diretora,ESMA e Diretoria Fórum**: Gratificações.
	- **Indenizações (Number)**: Auxílio-alimentação, Auxílio-transporte, Auxílio Pré-escolar, Auxílio Saúde, Auxílio Natalidade, Auxílio Moradia, Ajuda de Custo, além de outras desta natureza.
	- **Vantagens Eventuais (Number)**: Abono constitucional de 1/3 de férias, indenização de férias, antecipação de férias, gratificação natalina, antecipação de gratificação natalina, serviço extraordinário, substituição, pagamentos retroativos, além de outras desta natureza.
- **Descontos:**
	- **Previdência Pública (Number)**: Contribuição Previdenciária Oficial (Plano de Seguridade Social do Servidor Público e Regime Geral de Previdência Social).
	- **Imposto de Renda (Number)**: Imposto de Renda Retido na Fonte
	- **Descontos Diversos**: Cotas de participação de auxílio pré-escolar, auxílio transporte e demais descontos extraordinários de caráter não pessoal. 
	- **Retenção por Teto Constitucional (Number)**: Valores retidos por excederem ao teto remuneratório constitucional conforme Resoluções nº 13 e 14, do CNJ.
- **Total de Debitos (Number)**:  Total dos descontos efetuados no mês
- **Total de Créditos (Number)**: Total dos rendimentos pagos no mês.
- **Total de Créditos Líquidos (Number)**: Rendimento líquido após os descontos referidos nos itens anteriores.
- **Diárias**:  Valor de diárias efetivamente pago no mês de referência, ainda que o período de afastamento se estenda para além deste.”

### Servidores:

- **Nome do Servidor (String)**: Nome completo do funcionário
- **Lotação (String)**: Local (cidade, departamento, promotoria) em que o funcionário trabalha
- **Cargo (String)**: Cargo do funcionário dentro do MP
- **Rendimentos:**
	- **Remuneração Paradigma (Number)**: Remuneração do cargo efetivo - Vencimento, G.A.J., V.P.I, Adicionais de Qualificação, G.A.E e G.A.S, além de outras desta natureza.
	- **Vantagens Pessoais (Number)**: V.P.N.I., Adicional por tempo de serviço, quintos, décimos e vantagens decorrentes de sentença judicial ou extensão administrativa, abono de permanência.
	- **Função de Confiança ou Cargo em Comissão**: Remuneração de cargo em comissão ou confiança (servidor sem vínculo ou requisitado)
	- **Indenizações (Number)**: Auxílio-alimentação, Auxílio-transporte, Auxílio Pré-escolar, Auxílio Saúde, Auxílio Natalidade, Auxílio Moradia, Ajuda de Custo, além de outras desta natureza.
	- **Vantagens Eventuais (Number)**: Abono constitucional de 1/3 de férias, indenização de férias, antecipação de férias, gratificação natalina, antecipação de gratificação natalina, serviço extraordinário, substituição, pagamentos retroativos, além de outras desta natureza.
	- **Total de Créditos (Number)**: Total dos rendimentos pagos no mês.
- **Descontos:**
	- **Previdência Pública (Number)**: Contribuição Previdenciária Oficial (Plano de Seguridade Social do Servidor Público e Regime Geral de Previdência Social).
	- **Imposto de Renda (Number)**: Imposto de Renda Retido na Fonte
	- **Descontos Diversos**: Cotas de participação de auxílio pré-escolar, auxílio transporte e demais descontos extraordinários de caráter não pessoal. 
	- **Retenção por Teto Constitucional (Number)**: Valores retidos por excederem ao teto remuneratório constitucional conforme Resoluções nº 13 e 14, do CNJ.
	- **Total de Debitos (Number)**:  Total dos descontos efetuados no mês
- **Rendimento Líquido (Number)**: Rendimento líquido após os descontos referidos nos itens anteriores.
- **Diárias**:  Valor de diárias efetivamente pago no mês de referência, ainda que o período de afastamento se estenda para além deste.”

## Arquivos
  
### Remunerações ###

- **URL Base**: [https://www.tjpb.jus.br/transparencia/gestao-de-pessoas/folha-de-pagamento-de-pessoal](https://www.tjpb.jus.br/transparencia/gestao-de-pessoas/folha-de-pagamento-de-pessoal)
- **Formato**: PDF
- **Tipos**: Existe um arquivo referente aos servidores e um referente aos magistrados.
	- **Obs**: Antes de outubro de 2012 pode-se ter apenas o arquivo de magistrados ou um arquivo geral que não discrimina por funcionário. Esses são nomeados apenas como "remuneracoes-tjpb-mes-ano.pdf" quando baixados.