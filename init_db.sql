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
    nao_requer_login BOOL NOT NULL ,    -- É necessário login para coleta dos dados?
    nao_requer_captcha BOOL NOT NULL ,    -- É necessário captcha para coleta dos dados?
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
 
    CONSTRAINT coleta_orgao_fk FOREIGN KEY (id_orgao) REFERENCES orgaos(id) ON DELETE CASCADE
);

CREATE INDEX coletas_indice ON coletas(id_orgao,mes,ano);

CREATE TABLE contracheques(
    id INT NOT NULL,    -- Identificador da folha de pagamento dentro de uma coleta. Exemplos: 1, 2, 3...
    id_coleta VARCHAR(25) NOT NULL,    -- Identificador da coleta associada a folha de pagamento
    id_orgao VARCHAR(10) NOT NULL,    -- A sigla do órgão em minúsculo. Exemplos tjal, mpms, mpam...
    mes INT NOT NULL,    -- O mês que os dados coletados se referem. 
    ano INT NOT NULL,    -- O ano que os dados coletados se referem. 
    nome TEXT,    -- Nome do servidor público
    matricula VARCHAR(50),    -- Matrícula do servidor público
    cargo TEXT,    -- Cargo do servidor público. Exemplos : PROMOTOR DE JUSTICA DE 1ª ENTRÂNCIA, PROCURADOR DE JUSTICA...
    lotacao TEXT,    -- O local onde o membro está lotado.

    CONSTRAINT contracheque_orgao_fk FOREIGN KEY (id_orgao) REFERENCES orgaos(id) ON DELETE CASCADE,
    CONSTRAINT contracheques_pk PRIMARY KEY (id, id_coleta),
    CONSTRAINT contracheque_coleta_fk FOREIGN KEY (id_coleta) REFERENCES coletas(id) ON DELETE CASCADE
);

CREATE INDEX contracheques_indice ON contracheques(id_orgao,mes,ano);

CREATE TABLE remuneracoes(
    id INT NOT NULL, -- Identificador da remuneração dentro de uma folha de pagamento. Exemplos: 1, 2, 3...,
    id_contracheque INT NOT NULL,    -- Identificador da folha de pagamento. Exemplos : 1, 2, 3...
    id_coleta VARCHAR(25) NOT NULL,  -- Identificador da coleta: id_orgao/mes/ano
    detalhamento_contracheque TEXT NOT NULL,    -- Descrição do ítem de remuneração. Exemplos: diárias, auxílio-alimentação, auxílio moradia...
    valor DECIMAL NOT NULL,     -- Valor associado ao item de remuneração
    categoria_contracheque VARCHAR(15) NOT NULL,

    CONSTRAINT remuneracoes_pk PRIMARY KEY (id, id_contracheque, id_coleta),
    CONSTRAINT remuneracao_contracheque_fk FOREIGN KEY (id_contracheque,id_coleta) REFERENCES contracheques(id,id_coleta) ON DELETE CASCADE
);

CREATE INDEX remuneracoes_categoria ON remuneracoes(categoria_contracheque);

