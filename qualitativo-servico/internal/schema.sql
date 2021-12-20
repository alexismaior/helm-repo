attach ':memory:' as 'siop';

CREATE TABLE siop.localizador
(
    localizadorid              integer                                   NOT NULL primary key,
    acaoid                     integer                                   NOT NULL,
    municipioid                integer,
    ufid                       integer,
    regiaoid                   integer,
    localizador                character(4)                              NOT NULL,
    identificadorunico         integer                                   NOT NULL,
    descricao                  character varying(255)                    NOT NULL,
    repercussaofinanceira      text,
    valorrepercussaofinanceira numeric(16, 2),
    snexclusaologica           boolean                     DEFAULT false NOT NULL,
    indicadoralteracao         character(1),
    mesanotermino              date,
    mesanoinicio               date,
    valortotalfisico           numeric(16, 2),
    valortotalfinanceiro       numeric(16, 2),
    snnovo                     boolean                     DEFAULT false NOT NULL,
    datahoraalteracao          timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    sequencial                 integer,
    momentoid                  integer                                   NOT NULL,
    tipoinclusaolocalizadorid  integer                                   NOT NULL,
    previsaofinanceira         numeric(16, 2),
    metafisica                 numeric(16, 2),
    snatual                    boolean                     DEFAULT false NOT NULL,
    usuarioidenvio             integer,
    dataenvio                  timestamp without time zone,
    codigotemporario           character(4),
    recortegeograficoid        integer,
    snvalidado                 boolean                     DEFAULT false NOT NULL,
    versao                     integer                     DEFAULT 0
);

CREATE TABLE siop.esfera
(
    esferaid           integer not null primary key,
    esfera             varchar not null,
    descricao          varchar not null,
    descricaoabreviada varchar not null,
    snativo            boolean not null
);

CREATE TABLE siop.funcao
(
    funcaoid           integer not null primary key,
    funcao             varchar not null,
    exercicio          int     not null,
    descricao          varchar not null,
    descricaoabreviada varchar not null,
    snativo            boolean not null
);

CREATE TABLE siop.subfuncao
(
    subfuncaoid        integer not null primary key,
    subfuncao          varchar not null,
    funcaoid           integer not null,
    exercicio          int     not null,
    descricao          varchar not null,
    descricaoabreviada varchar not null,
    snativo            boolean not null
);

CREATE TABLE siop.programa
(
    programaid       integer not null primary key,
    programa         varchar not null,
    orgaoid          integer,
    exercicio        int     not null,
    titulo           varchar not null,
    snativo          boolean DEFAULT true not null,
    snexclusaologica boolean not null,
    snatual          boolean not null
);

CREATE TABLE siop.acao
(
    acaoid                          integer                                   NOT NULL primary key,
    funcaoid                        integer                                   NOT NULL,
    subfuncaoid                     integer                                   NOT NULL,
    orgaosiorgid                    integer,
    momentoid                       integer                                   NOT NULL,
    tipoacaoid                      integer                                   NOT NULL,
    programaid                      integer                                   NOT NULL,
    tipoinclusaoacaoid              integer                                   NOT NULL,
    esferaid                        integer                                   NOT NULL,
    orgaoid                         integer                                   NOT NULL,
    exercicio                       smallint                                  NOT NULL,
    anovaloracao                    smallint,
    anocadastramento                smallint,
    identificadorunico              integer                                   NOT NULL,
    acao                            character(4)                              NOT NULL,
    titulo                          character varying(255)                    NOT NULL,
    finalidade                      text,
    descricao                       text,
    baselegal                       text,
    unidaderesponsavel              character varying(120),
    repercussaofinanceira           text,
    valorrepercussaofinanceira      numeric(16, 2),
    sndireta                        boolean                     DEFAULT false NOT NULL,
    sndescentralizada               boolean                     DEFAULT false NOT NULL,
    snlinhacredito                  boolean                     DEFAULT false NOT NULL,
    sntransferenciaobrigatoria      boolean                     DEFAULT false NOT NULL,
    sntransferenciavoluntaria       boolean                     DEFAULT false NOT NULL,
    sntransferenciaoutras           boolean                     DEFAULT false NOT NULL,
    sndespesaobrigatoria            boolean,
    detalhamentoimplementacao       text,
    formaacompanhamento             text,
    identificacaosazonalidade       text,
    insumosutilizados               text,
    mesanotermino                   date,
    mesanoinicio                    date,
    snativo                         boolean                     DEFAULT true  NOT NULL,
    snexclusaologica                boolean                     DEFAULT false NOT NULL,
    usuarioidenvio                  integer,
    dataenvio                       timestamp without time zone,
    indicadoralteracao              character(1),
    sndetalhamentoetapaacao         boolean                     DEFAULT false NOT NULL,
    snatual                         boolean                     DEFAULT true  NOT NULL,
    snvalidado                      boolean                     DEFAULT false NOT NULL,
    valortotalfisico                numeric(16, 2),
    valortotalfinanceiro            numeric(16, 2),
    snempreendimentoppipac          boolean                     DEFAULT false NOT NULL,
    snnovo                          boolean                     DEFAULT false NOT NULL,
    datahoraalteracao               timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    sngrandevulto                   boolean                     DEFAULT false NOT NULL,
    subtipoacaoid                   integer,
    snplurianual                    boolean                     DEFAULT false NOT NULL,
    iniciativaid                    integer,
    acaopadronizadaid               integer,
    codigotemporario                character(4),
    versao                          integer                     DEFAULT 0,
    snaquisicaoinsumoestrategico    boolean                     DEFAULT false NOT NULL,
    sndetalhamentoplanoorcamentario boolean                     DEFAULT false NOT NULL,
    snregionalizarnaexecucao        boolean                     DEFAULT false NOT NULL,
    beneficiario                    text,
    snforapadronizacao              boolean                     DEFAULT false NOT NULL,
    produtoid                       integer,
    unidademedidaid                 integer,
    especificacaoproduto            text,
    snparticipacaosocial            boolean,
    detalhamentoparticipacaosocial  text
);