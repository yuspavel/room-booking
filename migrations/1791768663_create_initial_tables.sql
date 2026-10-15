-- +goose Up
CREATE TABLE room
(
    id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       VARCHAR(30) NOT NULL,
    capacity   INTEGER NOT NULL,
    hour_price NUMERIC(10,2) NOT NULL,
    currency   VARCHAR(4) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users
(
id         TEXT PRIMARY KEY DEFAULT gen_random_uuid()::TEXT,
name       VARCHAR(20) NOT NULL,
pass_hash   VARCHAR(100) NOT NULL,
email      VARCHAR(30) NOT NULL UNIQUE,
phone      VARCHAR(20),
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TYPE booking_status AS ENUM
(
    'pending',
    'confirmed',
    'cancelled',
    'completed'
);

CREATE TABLE booking
(
id         BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
room_id    INTEGER NOT NULL,
user_id    TEXT NOT NULL,
title      VARCHAR(50) NOT NULL,
description VARCHAR(200),
participants_count INTEGER NOT NULL,
from_time  TIMESTAMPTZ NOT NULL,
to_time    TIMESTAMPTZ NOT NULL,
status     booking_status NOT NULL DEFAULT 'pending',
paid       BOOLEAN,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE room;
DROP TABLE users;
DROP TABLE booking;

