drop table if exists "game_score" cascade;
drop table if exists "account" cascade;
drop table if exists "game" cascade;

create table "account" (
    "id" bigserial primary key,
    "email" character varying(255) not null,
    "password" character varying(255) not null,
    "created_at" timestamp without time zone not null default now(),
    "updated_at" timestamp without time zone not null default now()
);

insert into "account" ( "email", "password" ) values
        ( 'jean.hamond@corp.online.com', 'totoTOTO89' ),
        ( 'lise.remeur@corp.onlone.com', 'totoTOTO89' );

create table "game" (
    "id" bigserial primary key,
    "type" character varying(255) not null,
	"status" character varying(255) not null,
    "created_at" timestamp without time zone not null default now(),
    "updated_at" timestamp without time zone not null default now()
);

insert into "game" ("id","type","status") values ( 99,'babyfoot','finished');
insert into "game" ("id","type","status") values (100,'babyfoot','finished');

create table "game_score" (
    "id" bigserial primary key,
    "game_id" bigint not null,
    "account_id" bigint not null,
    "score" bigint not null default 0,
    "winner" boolean not null default false,
    "created_at" timestamp without time zone not null default now(),
    "updated_at" timestamp without time zone not null default now()
);


--insert into "game_score" ("id","game_id","type","status") values ( 99,99,'babyfoot','finished');
--insert into "game_score" ("id","game_id","type","status") values (100,99,'babyfoot','finished');
--
--create table "game_score" (
--    "id" bigserial primary key,
--    "account_id" bigint not null,
--    "score" bigint not null default 0,
--    "winner" boolean not null default false,
--    "created_at" timestamp without time zone not null default now(),
--    "updated_at" timestamp without time zone not null default now()
--);
--
--insert into "game_score" ("type","status") values ('babyfoot','finished');
--insert into "game_score" ("type","status") values ('babyfoot','finished');
