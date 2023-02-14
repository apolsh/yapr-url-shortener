BEGIN;

create table if not exists shortened_urls
(
    id           varchar(20)        not null,
    original_url varchar            not null,
    owner        uuid               not null,
    status       smallint default 0 not null
);

create unique index if not exists shortened_urls_id_uindex on shortened_urls (id);

create unique index if not exists shortened_urls_original_url_uindex on shortened_urls (original_url);

COMMIT;