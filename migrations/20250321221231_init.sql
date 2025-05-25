-- +goose Up
-- +goose StatementBegin
create table "USERS"
(
    id                       uuid        not null,
    email                    text,
    first_name               text        not null,
    last_name                text        not null,
    is_registration_complete boolean     not null,
    created_at               timestamptz not null,
    PRIMARY KEY (id),
    CONSTRAINT "USERS_EMAIL_UNIQUE" unique (email)
);

create table "DEVICES"
(
    id            uuid        not null,
    user_id       uuid        not null,
    device_name   text        not null,
    created_at    timestamptz not null,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE
);

CREATE TABLE "USER_PASSWORDS"
(
    user_id       UUID not null,
    password_hash TEXT NOT NULL,
    salt          TEXT NOT NULL,
    CONSTRAINT "USER_PASSWORDS_USER_ID_UNIQUE" PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE
);


create table "AUTH_PROVIDERS"
(
    id               uuid        not null,
    user_id          uuid        not null,
    provider_name    text        not null,
    created_at       timestamptz not null,
    provider_user_id BIGINT      NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    -- один и тот же Telegram-аккаунт не сможет быть привязан к нескольким пользователям.
    constraint "AUTH_PROVIDERS_PROVIDER_USER_ID_UNIQUE" unique (provider_name, provider_user_id),
    -- один пользователь может иметь не более одного аккаунта от одного провайдера
    constraint "AUTH_PROVIDERS_USER_ID_PROVIDER_NAME_UNIQUE" unique (user_id, provider_name)
);

CREATE TABLE "EXTERNAL_APPS"
(
    id            UUID        not null,
    name          TEXT        not null,
    owner_user_id UUID        not null,
    api_key_hash  TEXT        not null,
    created_at    timestamptz not null,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    CONSTRAINT "EXTERNAL_APPS_OWNER_USER_ID_NAME" UNIQUE (owner_user_id, name)
);

CREATE TABLE "ROLES"
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE "PERMISSIONS"
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE "ROLE_PERMISSIONS"
(
    role_id       INT NOT NULL,
    permission_id INT NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES "ROLES" (id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES "PERMISSIONS" (id) ON DELETE CASCADE
);

CREATE TABLE "USER_ROLES"
(
    user_id UUID NOT NULL,
    role_id INT  NOT NULL,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES "ROLES" (id) ON DELETE CASCADE
);

CREATE TABLE "EXTERNAL_APPS_ROLES"
(
    external_app_id UUID NOT NULL,
    role_id         INT  NOT NULL,
    PRIMARY KEY (external_app_id, role_id),
    FOREIGN KEY (external_app_id) REFERENCES "EXTERNAL_APPS" (id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES "ROLES" (id) ON DELETE CASCADE
);

INSERT INTO "ROLES" (name)
VALUES ('admin');
INSERT INTO "ROLES" (name)
VALUES ('user');
INSERT INTO "ROLES" (name)
VALUES ('external_app_reader');
INSERT INTO "ROLES" (name)
VALUES ('external_app_writer');

INSERT INTO "PERMISSIONS" (name)
VALUES ('read_own_profile'),
       ('update_own_profile'),

-- Работа с устройствами
       ('register_device'),
       ('remove_device'),
       ('read_own_devices'),


-- Работа с данными о здоровье (только свои)
       ('read_own_steps'),
       ('write_own_steps'),
       ('read_own_heart_rates'),
       ('write_own_heart_rates'),
       ('read_own_temperatures'),
       ('write_own_temperatures'),
       ('read_own_weights'),
       ('write_own_weights'),
       ('read_own_sleeps'),
       ('write_own_sleeps'),

-- Административные
       ('read_all_users'),
       ('delete_user'),
       ('assign_roles'),
       ('read_all_devices'),
       ('read_all_health_data'),


--  Работа с внешними приложениями
       ('read_own_external_apps'),
       ('update_own_external_apps');

-- Роль: user
INSERT INTO "ROLE_PERMISSIONS" (role_id, permission_id)
SELECT r.id, p.id
FROM "ROLES" r,
     "PERMISSIONS" p
WHERE r.name = 'user'
  AND p.name IN (
                 'read_own_profile',
                 'update_own_profile',
                 'register_device',
                 'remove_device',
                 'read_own_devices',
                 'read_own_steps',
                 'write_own_steps',
                 'read_own_heart_rates',
                 'write_own_heart_rates',
                 'read_own_temperatures',
                 'write_own_temperatures',
                 'read_own_weights',
                 'write_own_weights',
                 'read_own_sleeps',
                 'write_own_sleeps',
                 'read_own_external_apps',
                 'update_own_external_apps'
    );

-- Роль: external_app_read
INSERT INTO "ROLE_PERMISSIONS" (role_id, permission_id)
SELECT r.id, p.id
FROM "ROLES" r,
     "PERMISSIONS" p
WHERE r.name = 'external_app_reader'
  AND p.name IN (
                 'read_own_profile',
                 'read_own_devices',
                 'read_own_steps',
                 'read_own_heart_rates',
                 'read_own_temperatures',
                 'read_own_weights',
                 'read_own_sleeps'
    );


-- Роль: external_app_write
INSERT INTO "ROLE_PERMISSIONS" (role_id, permission_id)
SELECT r.id, p.id
FROM "ROLES" r,
     "PERMISSIONS" p
WHERE r.name = 'external_app_writer'
  AND p.name IN (
                 'write_own_steps',
                 'write_own_heart_rates',
                 'write_own_temperatures',
                 'write_own_weights',
                 'write_own_sleeps',
                 'update_own_external_apps'
    );

-- Роль: admin
INSERT INTO "ROLE_PERMISSIONS" (role_id, permission_id)
SELECT r.id, p.id
FROM "ROLES" r,
     "PERMISSIONS" p
WHERE r.name = 'admin';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "EXTERNAL_APPS_ROLES";
drop table "USER_ROLES";
drop table "ROLE_PERMISSIONS";
drop table "PERMISSIONS";
drop table "ROLES";
drop table "EXTERNAL_APPS";
drop table "USER_PASSWORDS";
drop table "AUTH_PROVIDERS";
drop table "DEVICES";
drop table "USERS";
-- +goose StatementEnd
