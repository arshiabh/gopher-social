create table if not exists comments (
    id bigserial primary key,
    user_id bigserial not null,
    post_id bigserial not null,
    content text not null,
    created_at timestamp(0) with time zone NOT NULL default now()
);