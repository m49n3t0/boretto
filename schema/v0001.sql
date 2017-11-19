
--+Migrate Down

DROP TABLE IF EXISTS "robot";

DROP TABLE IF EXISTS "http_endpoint";
DROP TABLE IF EXISTS "endpoint";

DROP TABLE IF EXISTS "task";

--+Migrate Up

CREATE TABLE "task" (
    "id" BIGSERIAL PRIMARY KEY NOT NULL,
    "version" BIGINT NOT NULL,
    "context" CHARACTER VARYING(255) NOT NULL,
    "function" CHARACTER VARYING(255) NOT NULL,
    "step" CHARACTER VARYING(255) NOT NULL,
    "status" CHARACTER VARYING(255) NOT NULL,
    "retry" BIGINT NOT NULL DEFAULT 8,
    "arguments" JSONB NOT NULL DEFAULT '{}',
    "buffer" JSONB NOT NULL DEFAULT '{}',
    "todo_date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "creation_date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_update" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE "endpoint" (
    "id" BIGSERIAL PRIMARY KEY NOT NULL,
    "type" CHARACTER VARYING(255) NOT NULL,
    "version" BIGINT NOT NULL,
    "name" CHARACTER VARYING(255) NOT NULL,
    "creation_date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_update" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT "endpoint_unique_name" UNIQUE("name")
);

CREATE TABLE "http_endpoint" (
    "method" CHARACTER VARYING(255) NOT NULL,
    "url" CHARACTER VARYING(255) NOT NULL,
    CONSTRAINT "endpoint_http_check" CHECK ("type"='HTTP')
) INHERITS ("endpoint");

CREATE TABLE "robot" (
    "id" BIGSERIAL PRIMARY KEY NOT NULL,
    "function" CHARACTER VARYING(255) NOT NULL,
    "version" BIGINT NOT NULL,
    "status" CHARACTER VARYING(255) NOT NULL,
    "definition" JSONB NOT NULL DEFAULT '{}',
    "creation_date" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_update" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

insert into "http_endpoint" ( "type", "version", "name", "method", "url" ) values
                            ( 'HTTP', '1', 'starting', 'GET', 'https://api.com/starting' ),
                            ( 'HTTP', '1', 'checking', 'GET', 'https://api.com/checking' ),
                            ( 'HTTP', '1', 'onServer', 'GET', 'https://api.com/onServer' ),
                            ( 'HTTP', '1', 'onInterne', 'GET', 'https://api.com/onInterne' ),
                            ( 'HTTP', '1', 'ending', 'GET', 'https://api.com/ending' );

insert into "robot" ( "function", "version", "status", "definition" ) values
                    ( 'database/create', '1', 'ACTIVE', '{"sequence":[{"name":"STARTING","endpoint_type":"HTTP","endpoint_id":1},{"name":"CHECKING","endpoint_type":"HTTP","endpoint_id":2},{"name":"ON_SERVER","endpoint_type":"HTTP","endpoint_id":3},{"name":"ON_INTERNE","endpoint_type":"HTTP","endpoint_id":4},{"name":"ENDING","endpoint_type":"HTTP","endpoint_id":5}]}' );


insert into "task" ( "version", "context", "function", "step", "status", "retry" ) values
                    ( '1', 'toto', 'database/create', 'starting', 'TODO', '8' );







