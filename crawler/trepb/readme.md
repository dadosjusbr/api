# Tribunal Regional Eleitoral da Paraíba

Este crawler tem como objetivo a recuperação de informações sobre folhas de pagamentos dos funcionários do Tribunal Regional Eleitoral da Paraíba. O site com as informações pode ser acessado [aqui](http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?acao=Anexo_VIII).

O crawler está estruturado como uma CLI. Você passa quatro argumentos (mês, ano, cpf e nome) e é baixado um arquivo no formato **HTML** representando a folha de pagamento da instituição. O arquivo contêm dois elementos table, sendo a primeira um dicionário dos dados e a segunda a tabela de remunerações.

Os dados de nome e cpf são necessários para obter uma chave de acesso para a api fornecida pelo TRE-PB. O crawler criará um cache dessa chave para diminuir a necessidade de requisições de chave para cada uso.

## Legislação

Os dados devem estar de acordo com a [Resolução 102 do CNJ](https://atos.cnj.jus.br/atos/detalhar/69).

## Como usar

- É preciso ter o compilador de Go instalado em sua máquina. Mais informações [aqui](https://golang.org/dl/).
- Rode o comando abaixo, com mês, ano, cpf e nome que você quer ter acesso as informações
- O nome deve ser completo. O cpf deve ter o formato xxx.xxx.xxx-xx. Informações para obter chave de acesso para a API.

```sh
cd crawler/trepb
go build
./trepb --mes=${MES} --ano=${ANO} --nome="${NOME COMPLETO}" --cpf="${CPF}"
```

## Dicionário de Dados

Para cada funcionário, o JSON possui os seguintes campos:

- **Nome (String)**: Nome completo do funcionário
- **Lotação (String)**: Local (cidade, departamento, promotoria) em que o funcionário trabalha
- **Cargo (String)**: Cargo do funcionário dentro do MP
- **Rendimentos:**
	- **Remuneração Paradigma (Number)**: Remuneração do cargo efetivo - Vencimento, G.A.J., V.P.I, Adicionais de Qualificação, G.A.E e G.A.S, além de outras desta natureza.
	- **Vantagens Pessoais (Number)**: V.P.N.I., Adicional por tempo de serviço, quintos, décimos e vantagens decorrentes de sentença judicial ou extensão administrativa, abono de permanência.
	- **Subsídio, FC e CJ (Number)**: Subsídios, diferença de subsídios, função de confiança e cargo em comissão.
	- **Indenizações (Number)**: Auxílio-alimentação, Auxílio-transporte, Auxílio Pré-escolar, Auxílio Saúde, Auxílio Natalidade, Auxílio Moradia, Ajuda de Custo, além de outras desta natureza.
	- **Vantagens Eventuais (Number)**: Abono constitucional de 1/3 de férias, indenização de férias, antecipação de férias, gratificação natalina, antecipação de gratificação natalina, serviço extraordinário, substituição, pagamentos retroativos, além de outras desta natureza.
	- **Gratificações (Number)**
	- **Total de Créditos (Number)**: Total dos rendimentos pagos no mês.
- **Descontos:**
	- **Contribuição Previdenciária (Number)**: Contribuição Previdenciária Oficial (Plano de Seguridade Social do Servidor Público e Regime Geral de Previdência Social).
	- **Imposto de Renda (Number)**: Imposto de Renda Retido na Fonte
	- **Descontos Diversos**: Cotas de participação de auxílio pré-escolar, auxílio transporte e demais descontos extraordinários de caráter não pessoal. 
	- **Retenção por Teto Constitucional (Number)**: Valores retidos por excederem ao teto remuneratório constitucional conforme Resoluções nº 13 e 14, do CNJ.
	- **Total de Debitos (Number)**:  Total dos descontos efetuados no mês
- **Rendimento Líquido (Number)**: endimento líquido após os descontos referidos nos itens anteriores.
- **Remuneração do órgão de origem (Number)**: Remuneração percebida no órgão de origem por magistrados e servidores, cedidos ou requisitados, optantes por aquela remuneração.
- 
## Arquivos
  
### Remunerações ###

- **URL Base**: [http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?](http://apps.tre-pb.jus.br/transparenciaDadosServidores/infoServidores?)
    - **Parâmetros da URL**: mes=${MES}&ano=${ANO}&chaveDeAcesso=${CHAVE}
    - **Parâmetros fixos**: acao=AnexoVIII&folha=&valida=true&toExcel=false
- **Formato**: Table html