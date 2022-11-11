\connect dadosjusbr;

CREATE TABLE orgaos(
    id VARCHAR(10) PRIMARY KEY,    -- A sigla do órgão em minúsculo. Exemplos tjal, mpms, mpam...
    nome VARCHAR(100) NOT NULL,    -- O nome do órgão. Exemplos: Tribunal de Justiça do Estado de Alagoas,Tribunal de Justiça de Pernambuco...
    jurisdicao VARCHAR(25) NOT NULL,    -- A jurisdição do órgão. Exemplos: Estadual (e), Federal(f) ou Distrito Federal e Territórios (d).
    entidade VARCHAR(25) NOT NULL,    -- Se é tribunal de justiça (tj), ministério público(mp), defensoria(d) ou procuradoria (p).
    uf VARCHAR(25) NOT NULL    -- Unidade federativa do órgão. Ficará vazio para o caso de órgãos com jurisdição federal ou distrito federal. Exemplos: PB, RS, PE,BA...
);
 
CREATE TABLE coletas(
    id VARCHAR(25) PRIMARY KEY,  -- chave primária da coleta: id_orgao/mes/ano
    id_orgao VARCHAR(10) NOT NULL,    -- A sigla do órgão em minúsculo. Exemplos tjal, mpms, mpam...
    mes INT NOT NULL,    -- O mês que os dados coletados se referem. 
    ano INT NOT NULL,    -- O ano que os dados coletados se referem. 
    timestamp TIMESTAMP NOT NULL,    -- Marca temporal em que o dado foi coletado.
    repositorio_coletor VARCHAR(150) NOT NULL,    -- URL do repositório do coletor dos dados relacionados a coleta.
    versao_coletor VARCHAR(25) NOT NULL,    -- Versão (identificador do commit) do repositório do coletor dos dados relacionados a coleta
    repositorio_parser VARCHAR(150), -- URL do repositório do parser dos dados relacionados a coleta. Somente preenchido quando o parser é diferente do coletor
    versao_parser VARCHAR(25), -- Versão (identificador do commit) do repositório do parser dos dados relacionados a coleta. Somente preenchido quando o parser é diferente do coletor
    estritamente_tabular BOOL NOT NULL ,    -- Órgãos que disponibilizam dados limpos (tidy data)
    formato_consistente BOOL NOT NULL ,    -- Órgão alterou a forma de expor seus dados entre o mês em questão e o mês anterior?
    tem_matricula BOOL NOT NULL ,    -- Órgão disponibiliza matrícula do servidor?
    tem_lotacao BOOL NOT NULL ,    -- Órgão disponibiliza lotação do servidor?
    tem_cargo BOOL NOT NULL ,    -- Órgão disponibiliza a função do servidor?
    acesso VARCHAR(50) NOT NULL,    -- Conseguimos prever/construir uma URL com base no órgão/mês/ano que leve ao download do dado? Exemplos : NECESSITA_SIMULAÇÃO_DO_USUÁRIO, AMIGAVEL_PARA_RASPAGEM, RASPAGEM_DIFICULTADA...
    extensao VARCHAR(25) NOT NULL,    -- Extensao do arquivo de dados. Exemplos: CSV, JSON, XLS, etc
    detalhamento_receita_base VARCHAR(25) NOT NULL,    -- Quão detalhado é a publicação da receita base. Exemplos: DETALHADO, SUMARIZADO...
    detalhamento_outras_receitas VARCHAR(25) NOT NULL,    -- Quão detalhado é a publicação das demais receitas. Exemplos: DETALHADO, SUMARIZADO...
    detalhamento_descontos VARCHAR(25) NOT NULL,    -- Quão detalhado é a publicação dos descontos. Exemplos: DETALHADO, SUMARIZADO...
    indice_completude DECIMAL NOT NULL,    -- Componente do índice de transparência resultante da análise dos metadados relacionados a disponibilidade dos dados
    indice_facilidade DECIMAL NOT NULL,    -- Componente do índice de transparência resultante da análise dos metadados relacionados a dificuldade para acessar os dados que estão disponíveis
    indice_transparencia DECIMAL NOT NULL,    -- Nota final, calculada utilizada os componentes de disponibilidade e dificuldade
    sumario JSON, -- JSON com algumas estatísticas referentes aos membros e remunerações de um órgão

    /*
    -- O formato de um sumário é parecido com isso:
    "sumario": {
        "membros":  214, -- Quantidade de membros ativos
        "remuneracao_base": {
            "maximo": 35462.22, -- Valor máximo de uma remuneração recebida por um membro
            "minimo": 7473.09, -- Valor mínimo de uma remuneração recebida por um membro
            "media":  33084.63186915894, -- Média das remunerações
            "total":  7080111.220000014, -- Total das remunerações
        },
        "outras_remuneracoes": {
            "maximo":  78348.23000000001, -- Valor máximo de uma outra remuneração recebida por um membro
            "media":  9524.795887850474, -- Média das outras remunerações
            "total":  2038306.3200000015, -- Total das outras remunerações
        },
        "histograma_renda": {
            "10000":  1,  -- Quantidade de membros que recebem até 10 mil reais
            "20000":  1,  -- Quantidade de membros que recebem entre 10 mil e 20 mil reais
            "30000":  3,  -- Quantidade de membros que recebem entre 20 mil e 30 mil reais
            "40000":  3,  -- Quantidade de membros que recebem entre 30 mil e 40 mil reais
            "50000":  0,  -- Quantidade de membros que recebem entre 40 mil e 50 mil reais
            "-1":     1,  -- Quantidade de membros que recebem mais de 50 mil reais
        }
    }
    */
    CONSTRAINT coleta_orgao_fk FOREIGN KEY (id_orgao) REFERENCES orgaos(id) ON DELETE CASCADE
);

CREATE INDEX coletas_indice ON coletas(id_orgao,mes,ano);

CREATE TABLE remuneracoes(
    id_orgao VARCHAR(10) NOT NULL,    -- A sigla do órgão em minúsculo. Exemplos tjal, mpms, mpam...
    mes INT NOT NULL,    -- O mês que os dados coletados se referem. 
    ano INT NOT NULL,    -- O ano que os dados coletados se referem. 
    linhas_descontos INT NOT NULL, -- Número de linhas para o determinado contracheque
    linhas_base INT NOT NULL, -- Número de linhas para o determinado contracheque
    linhas_outras INT NOT NULL, -- Número de linhas para o determinado contracheque
    zip_url TEXT NOT NULL, -- Link para o zip que contém os dados da remuneração

    CONSTRAINT remuneracoes_pk PRIMARY KEY (id_orgao, mes, ano )
);

