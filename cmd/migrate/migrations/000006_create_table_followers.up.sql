CREATE TABLE if not exists followers ( 
    follower_id bigint not null,
    user_id bigint not null,
    created_at timestamp(0) with time zone default now(),
    primary key(user_id, follower_id),
    foreign key(user_id) references users(id) on delete cascade, 
    foreign key(user_id) references users(id) on delete cascade 
);