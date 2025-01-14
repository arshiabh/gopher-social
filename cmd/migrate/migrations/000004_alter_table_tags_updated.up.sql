alter table posts
  add column tags varchar(100) [];
alter table posts
  add column updated_at timestamp(0) with time zone NOT NULL default now();