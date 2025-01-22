alter table user_invitation
add COLUMN expiry timestamp(0) with time zone not null;