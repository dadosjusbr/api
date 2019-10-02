# Ministério Público da Paraíba - Crawler

Este crawler tem como objetivo a recuperação de informações sobre folhas de pagamentos dos funcionários do Ministério Público da Paraíba. O site com as informações pode ser acessado [aqui](http://pitagoras.mppb.mp.br/PTMP/FolhaListar).
O crawler está estruturado como uma CLI. Você passa dois argumentos (mês e ano) e serão baixadas seis planilhas no formato ODS, cada planilha é referente a uma destas categorias: Membros Ativos, Membros Inativos, Servidores Ativos, Servidores Inativos, Servidores à Disposição e Aposentados/Pensionistas
.

## Como usar

- É preciso ter o compilador de Go instalado em sua máquina. Mais informações [aqui](https://golang.org/dl/).
- Rode o comando abaixo, com mês e ano que você quer ter acesso as informações

```sh
cd crawler/mppb
go run crawler-mppb.go --mes=${MES} --ano=${ANO}
```

## Dicionário de Dados

As planilhas possuem as seguintes colunas:

- **Nome (String)**: Nome completo do funcionário
- **Matrícula (String)**: Matrícula do funcionário  
- **Cargo (String)**: Cargo do funcionário dentro do MP
- **Lotação (String)**: Local (cidade, departamento, promotoria) em que o funcionário trabalha
- **Remuneração do cargo efetivo (Number)**: Vencimento, GAMPU, V.P.I, Adicionais de Qualificação, G.A.E e G.A.S, além de outras desta natureza. Soma de todas essas remunerações
- **Outras Verbas Remuneratórias, Legais ou Judiciais (Number)**: V.P.N.I., Adicional por tempo de serviço, quintos, décimos e vantagens decorrentes de sentença judicial ou extensão administrativa
- **Função de Confiança ou Cargo em Comissão (Number)**: Rubricas que representam a retribuição paga pelo exercício de função (servidor efetivo) ou remuneração de cargo em comissão (servidor sem vínculo ou requisitado)
- **Gratificação Natalina (Number)**: Parcelas da Gratificação Natalina (13º) pagas no mês corrente, ou no caso de vacância ou exoneração do servidor
- **Férias - ⅓ Constitucional (Number)**: Adicional correspondente a 1/3 (um terço) da remuneração, pago ao servidor por ocasião das férias
- **Abono de Permanência (Number)**:  Valor equivalente ao da contribuição previdenciária, devido ao funcionário público que esteja em condição de aposentar-se, mas que optou por continuar em atividade (instituído pela Emenda Constitucional nº 41, de 16 de dezembro de 2003
- **Total de Rendimentos Brutos (Number)**: Total dos rendimentos brutos pagos no mês
- **Contribuição Previdenciária (Number)**: Contribuição Previdenciária Oficial (Plano de Seguridade Social do Servidor Público e Regime Geral de Previdência Social)
- **Imposto de Renda (Number)**: Imposto de Renda Retido na Fonte
- **Retenção por Teto Constitucional (Number)**: Valor deduzido da remuneração básica bruta, quando esta ultrapassa o teto constitucional, nos termos da legislação correspondente
- **Total de Descontos (Number)**:  Total dos descontos efetuados no mês
- **Rendimento Líquido Total (Number)**: Rendimento líquido após os descontos referidos nos itens anteriores
- **Indenizações (Number)**: Auxílio-alimentação, Auxílio-transporte, Auxílio-Moradia, Ajuda de Custo e outras dessa natureza, exceto diárias, que serão divulgadas no Portal da Transparência. Soma de todas essas remunerações
- **Outras Remunerações Temporárias (Number)**: Valores pagos a título de Adicional de Insalubridade ou de Periculosidade, Adicional Noturno, Serviço Extraordinário, Substituição de Função, 'Atrasados'. Soma de todas essas remunerações

## Planilhas

Seguem o formato do Anexo 1 - CPJ 17/2012 (Similar à Tabela I CNMP 200/2019)

- Lista de planilhas: [http://pitagoras.mppb.mp.br/PTMP/FolhaListar](http://pitagoras.mppb.mp.br/PTMP/FolhaListar)
  
### Remunerações ###

- URL Base: [http://pitagoras.mppb.mp.br/PTMP/FolhaPagamentoExercicioMesNewOds](http://pitagoras.mppb.mp.br/PTMP/FolhaPagamentoExercicioMesNewOds)
- Parâmetros da url: exercicio=[ano]&mes=[mes]&tipo=[Membros Ativos (*1*), Membros Inativos (*2*), Servidores Ativos (*3*), Servidores Inativos (*4*), Servidores à Disposição (*5*), Aposentados/Pensionistas (*6*)]
- [Exemplo](http://pitagoras.mppb.mp.br/PTMP/FolhaPagamentoExercicioMesNewOds?mes=1&exercicio=2019&tipo=1)
  


- Exemplos:
  - 
  - [Verbas Referentes a exercícios anteriores](http://pitagoras.mppb.mp.br/PTMP/FolhaExercicioAnteriorMesNewOds?exercicio=2019&mes=1)
  - [Remuneração de estagiários](http://pitagoras.mppb.mp.br/PTMP/FolhaPagamentoEstagiarioExercicioMesOds?mes=1&exercicio=2018)


### **-  Tabela de remuneração de estagiários**:  
[http://pitagoras.mppb.mp.br/PTMP/FolhaPagamentoEstagiarioExercicioMesOds?mes=1&exercicio=2018](http://pitagoras.mppb.mp.br/PTMP/FolhaPagamentoEstagiarioExercicioMesOds?mes=1&exercicio=2018)
-   **Formato da tabela:**  Anexo 1 - CPJ 17/2012 (Similar à Tabela I CNMP 200/2019) com algumas diferenças de nomenclatura:
	- Remuneração de cargo efetivo: Remuneração
	- Função de Confiança ou Cargo em Comissão: Função de Confiança
	- Gratificação Natalina: 13º Vencimento 
	- Férias (⅓ constitucional): Adicional de Férias (Constitucional)
-  **Formato da url:** exercicio=ano, mes=mes, tipo=tipo de servidor

### **-  Tabela de Indenizações**:  
https://pitagoras.mppb.mp.br/PTMP/FolhaVerbaIndenizRemTemporariaOds?mes=1&exercicio=2019&tipo=

-   **Formato da tabela:**  Possui em comum com as tabelas anteriores apenas os campos: Matrícula, Nome, Cargo e lotação. É complemetada por duas supercolunas com cada indenização discriminadas, são elas:
	- Verbas Indenizatórias 1
	- Outras remunerações temporárias 2
-  **Formato da url:** exercicio=ano, mes=mes
- **Obs**: Indenizações disponíveis apenas para o ano de 2019.
