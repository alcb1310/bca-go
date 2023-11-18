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

create table if not exists project(
     id uuid PRIMARY KEY default gen_random_uuid(),
     name varchar(255) not null,
     is_active boolean not null default true,
     created_at timestamp with time zone default now(),

     company_id uuid not null references company(id) on delete restrict,

     unique(company_id, name)
);

create table if not exists supplier (
     id uuid PRIMARY KEY default gen_random_uuid(),
     name text not null,
     supplier_id text not null,
     contact_name text,
     contact_email text,
     contact_phone text,
     created_at timestamp with time zone  default now(),

     company_id uuid not null references company(id) on delete restrict,
     unique(name, company_id),
     unique(supplier_id, company_id)
);

create table if not exists budget_item (
     id uuid PRIMARY KEY default gen_random_uuid(),
     code text not null,
     name varchar(255) not null,
     accumulates boolean not null default true,
     level integer not null,
     created_at timestamp with time zone default now(),

     parent_id uuid references budget_item(id) on delete restrict,
     company_id uuid not null references company(id) on delete restrict,

     unique(code, company_id),
     unique(name, company_id)
);

create table if not exists budget (
     id uuid PRIMARY KEY default gen_random_uuid(),
     project_id uuid not null references project(id) on delete restrict,
     budget_item_id uuid not null references budget_item(id) on delete restrict,

     initial_quantity double precision,
     initial_cost double precision,
     initial_total double precision not null,

     spent_quantity double precision,
     spent_total double precision not null,

     to_spend_quantity double precision,
     to_spend_cost double precision,
     to_spend_total double precision not null,

     updated_budget_quantity double precision not null,
     created_at timestamp with time zone default now(),

     company_id uuid not null references company(id) on delete restrict,

     unique(project_id, budget_item_id, company_id)
);

--  REQUIRED VIEWS

create or replace view user_without_password as
select c.id company_id, c.ruc ruc, c.name company_name, c.employees employees, u.id user_id,
u.email user_email, u.name user_name, r.name role_name, r.id role_id
from "user" u
inner join role r on u.role_id = r.id
inner join company c on u.company_id = c.id;

create or replace view budget_item_with_parents as
select b.id id, b.code code, b.name as name, b.accumulates accumulates, b.level as level, p.code parent_code,
p.name parent_name, p.id parent_id, b.company_id company_id
from budget_item b
left outer join budget_item p on b.parent_id = p.id;

