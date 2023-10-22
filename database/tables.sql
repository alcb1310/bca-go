-- REQUIRED TABLES

create table if not exists company (
     id uuid PRIMARY KEY  default gen_random_uuid(),
     ruc text not null,
     name text not null,
     employees smallint  default 1,
     is_active boolean  default true,
     created_at timestamp with time zone  default now()
);

alter table company drop constraint if exists company_name_key;
alter table company drop constraint if exists company_ruc_key;

alter table company add constraint company_name_key unique (name);
alter table company add constraint company_ruc_key unique (ruc);

create table if not exists role (
     id char(2) PRIMARY KEY,
     name varchar(255) not null,
     created_at timestamp with time zone default now()
);

create table if not exists "user" (
     id uuid PRIMARY KEY default gen_random_uuid(),
     email text not null, 
     name text not null, 
     password text not null,
     created_at timestamp with time zone  default now(),

     company_id uuid not null references company(id) on delete restrict,
     role_id char(2) references role(id) on delete restrict,

     unique (email)
);

create table if not exists logged_in_user (
     user_id uuid not null,
     token bytea not null,

     unique (user_id)
);

--  REQUIRED VIEWS

create or replace view user_without_password as
select c.id company_id, c.ruc ruc, c.name company_name, c.employees employees, u.id user_id,
u.email user_email, u.name user_name, r.name role_name, r.id role_id
from "user" u
inner join role r on u.role_id = r.id
inner join company c on u.company_id = c.id;
