-- +goose Up
-- +goose StatementBegin
create table "STEPS"
(
    id         uuid        not null,
    user_id    uuid        not null,
    device_id  uuid        not null,
    step_count integer     not null,
    created_at timestamptz not null,
    CONSTRAINT "STEPS_PK" PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES "DEVICES" (id) ON DELETE CASCADE
);

create materialized view "STEPS_DAYLE_AGGREGATED" AS
SELECT user_id,
       device_id,
       DATE_TRUNC('day', created_at) AS day,
       SUM(step_count)               AS total_steps
FROM "STEPS"
GROUP BY user_id, device_id, day;

CREATE UNIQUE INDEX idx_daily_device_agg
    ON "STEPS_DAYLE_AGGREGATED" (user_id, device_id, day);

create table "HEART_RATES"
(
    id         uuid        not null,
    user_id    uuid        not null,
    device_id  uuid        not null,
    bpm        integer     not null,
    created_at timestamptz not null,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES "DEVICES" (id) ON DELETE CASCADE
);


create table "TEMPERATURES"
(
    id                  uuid        not null,
    user_id             uuid        not null,
    device_id           uuid        not null,
    temperature_celsius float       not null,
    created_at          timestamptz not null,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES "DEVICES" (id) ON DELETE CASCADE
);

create table "WEIGHTS"
(
    id         uuid        not null,
    user_id    uuid        not null,
    device_id  uuid        not null,
    weight_kg  float       not null,
    created_at timestamptz not null,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES "DEVICES" (id) ON DELETE CASCADE
);


CREATE TABLE "SLEEPS"
(
    id         UUID        NOT NULL,
    user_id    UUID        NOT NULL,
    device_id  uuid        NOT NULL,
    started_at timestamptz NOT NULL,
    ended_at   timestamptz NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES "USERS" (id) ON DELETE CASCADE,
    FOREIGN KEY (device_id) REFERENCES "DEVICES" (id) ON DELETE CASCADE
);

CREATE INDEX idx_slept_user_created ON "SLEEPS" (user_id, started_at);
CREATE INDEX idx_heart_rates_user_created ON "HEART_RATES" (user_id, created_at);
CREATE INDEX idx_temperatures_user_created ON "TEMPERATURES" (user_id, created_at);
CREATE INDEX idx_weights_user_created ON "WEIGHTS" (user_id, created_at);
CREATE INDEX idx_steps_user_created ON "STEPS" (user_id, created_at);

CREATE UNIQUE INDEX "UNIQUE_STEPS"
    ON "STEPS" (
                user_id,
                device_id,
                step_count,
                created_at
        );
CREATE UNIQUE INDEX "UNIQUE_HEART_RATES"
    ON "HEART_RATES" (
                      user_id,
                      device_id,
                      bpm,
                      created_at
        );
CREATE UNIQUE INDEX "UNIQUE_TEMPERATURES"
    ON "TEMPERATURES" (
                       user_id,
                       device_id,
                       temperature_celsius,
                       created_at
        );
CREATE UNIQUE INDEX "UNIQUE_WEIGHTS"
    ON "WEIGHTS" (
                  user_id,
                  device_id,
                  weight_kg,
                  created_at
        );
CREATE UNIQUE INDEX "UNIQUE_SLEEPS"
    ON "SLEEPS" (
                 user_id,
                 device_id,
                 started_at,
                 ended_at
        );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP materialized view "STEPS_DAYLE_AGGREGATED";
drop table "SLEEPS";
drop table "WEIGHTS";
drop table "TEMPERATURES";
drop table "HEART_RATES";
drop table "STEPS";
-- +goose StatementEnd
