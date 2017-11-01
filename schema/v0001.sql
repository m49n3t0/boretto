drop table if exists "task";
drop table if exists "endpoint_http";
drop table if exists "robot";

create table "task" (
    --"id" uuid PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "id" bigserial primary key not null,
    "version" bigint not null,
    "context" character varying(255) not null,
    "function" character varying(255) not null,
    "step" character varying(255) not null,
    "status" character varying(255) not null,
    "retry" bigint not null default 8,
    "arguments" jsonb not null default '{}',
    "buffer" jsonb not null default '{}'
);

create table "endpoint_http" (
    "id" bigserial primary key not null,
    "name" character varying(255) not null,
    "method" character varying(255) not null,
    "url" character varying(255) not null
);

create table "robot" (
    "id" bigserial primary key not null,
    "function" character varying(255) not null,
    "version" bigint not null,
    "definition" jsonb not null default '{}'
);


































