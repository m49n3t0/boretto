
--+Migrate Down

drop table if exists "robot";
drop table if exists "header";
drop table if exists "endpoint";
drop table if exists "task";

--+Migrate Up

create table "task" (
    "id" bigserial primary key not null,
    "version" bigint not null,
    "context" character varying(255) not null,
    "function" character varying(255) not null,
    "step" character varying(255) not null,
    "status" character varying(255) not null,
    "retry" bigint not null default 8,
    "creation_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now(),
    "todo_date" timestamp with time zone not null default now(),
    "done_date" timestamp with time zone,
    "arguments" jsonb not null default '{}',
    "buffer" jsonb not null default '{}',
    "comment" character varying(255)
);

create table "endpoint" (
    "id" bigserial primary key not null,
    "version" bigint not null,
    "name" character varying(255) not null,
    "method" character varying(255) not null,
    "url" character varying(255) not null,
    "creation_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now(),
    constraint "endpoint_unique_name" unique("name")
);

create table "header" (
    "id" bigserial primary key not null,
    "name" character varying(255) not null,
    "value" character varying(255) not null,
    "creation_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now()
);

create table "robot" (
    "id" bigserial primary key not null,
    "function" character varying(255) not null,
    "version" bigint not null,
    "status" character varying(255) not null,
    "definition" jsonb not null default '{}',
    "creation_date" timestamp with time zone not null default now(),
    "last_update" timestamp with time zone not null default now()
);

--insert into "http_endpoint" ( "type", "version", "name", "method", "url" ) values
--                            ( 'http', '1', 'starting', 'get', 'https://api.ovh.com/1.0/' ),
--                            ( 'http', '1', 'checking', 'get', 'https://api.ovh.com/1.0/checking' ),
--                            ( 'http', '1', 'onserver', 'get', 'https://api.ovh.com/1.0/onserver' ),
--                            ( 'http', '1', 'oninterne', 'get', 'https://api.ovh.com/1.0/oninterne' ),
--                            ( 'http', '1', 'ending', 'get', 'https://api.ovh.com/1.0/ending' );
--
--insert into "robot" ( "function", "version", "status", "definition" ) values
--                    ( 'database/create', '1', 'active', '{"sequence":[{"name":"starting","endpoint_type":"http","endpoint_id":1},{"name":"checking","endpoint_type":"http","endpoint_id":2},{"name":"on_server","endpoint_type":"http","endpoint_id":3},{"name":"on_interne","endpoint_type":"http","endpoint_id":4},{"name":"ending","endpoint_type":"http","endpoint_id":5}]}' );
--
--
--insert into "task" ( "version", "context", "function", "step", "status", "retry" ) values
--                    ( '1', 'toto', 'database/create', 'starting', 'todo', '8' );








