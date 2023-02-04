CREATE TABLE short_links
(
    id        serial primary key,
    full_url  varchar(255) not null ,
    short_url varchar(255) not null ,
    user_id   varchar(255),
    created_at timestamp not null default now()
)