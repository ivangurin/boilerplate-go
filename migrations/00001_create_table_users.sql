-- +goose Up
-- +goose StatementBegin
create table users (
    id bigserial primary key,
    name text,
    email text,
    password text,
    is_admin boolean not null default false,
    deleted boolean not null default false,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);

insert into users (name, email, password, is_admin, created_at, updated_at) values
    ('Admin', 'admin@example.com', '$2a$14$ZMeqfBsFJGv63xq5AKRwheUa.Aidv6/PWeXgyZemuEv68VwD24kBW', true, now(), now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd