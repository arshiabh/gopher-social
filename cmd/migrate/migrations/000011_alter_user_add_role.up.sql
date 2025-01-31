alter table if exists users
add COLUMN role_id int references roles(id) default 1;


--default 1 change to users here then drop default and set not null for it
update users 
set role_id = (select id from roles where name = 'user');

alter table users
    alter COLUMN role_id drop default;

alter table users
    alter COLUMN role_id set not null;
