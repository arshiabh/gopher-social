create table if not exists user_invitation (
    token bytea primary key,
    user_id bigint not null
);