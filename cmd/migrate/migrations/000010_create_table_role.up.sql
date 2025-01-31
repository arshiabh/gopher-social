create table if not exists roles (
    id bigserial primary key,
    name varchar(255) not null unique,
    description text,
    level int not null default 0
);

insert into roles (name, description, level)
values (
    'user',
    'user can create posts and comment',
    1
);

insert into roles (name, description, level)
values (
    'moderator',
    'moderator can update other users posts ',
    2
);

insert into roles (name, description, level)
values (
    'admin',
    'admin can update and delete other user posts',
    3
);
