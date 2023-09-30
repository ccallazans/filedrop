create table if not exists roles (
    id serial not null,
    role varchar(255) unique not null,

    constraint pk_roles primary key(id)
);

insert into roles(role) values ('ADMIN'), ('USER'), ('GUEST');

create table if not exists users (
    id UUID not null,
    first_name varchar(255) not null,
    email varchar(255) unique not null,
    password varchar(255) not null,
    role_id serial not null,
    created_at timestamp not null,
    updated_at timestamp not null,

    constraint pk_users primary key(id),
    constraint fk_users_roles foreign key (role_id) references roles(id)
);

create table if not exists files (
    id uuid not null,
    filename varchar(255) not null,
    password varchar(255) not null,
    hash varchar(255) not null,
    location varchar(255) not null,
    is_active boolean not null,
    user_id UUID not null,
    created_at timestamp not null,
    updated_at timestamp not null,

    constraint pk_files primary key(id),
    constraint fk_files_users foreign key (user_id) references users(id)
);