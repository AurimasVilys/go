CREATE DATABASE scootin;
USE scootin;

CREATE TABLE scooter
(
    identifier               varchar(64) PRIMARY KEY NOT NULL DEFAULT (UUID_TO_BIN(UUID())),
    occupied_user_identifier varchar(64)
);

CREATE TABLE scooter_event
(
    scooter_identifier varchar(64)  NOT NULL,
    event              varchar(255) NOT NULL,
    timestamp          DATETIME     NOT NULL DEFAULT NOW(),
    latitude           varchar(255) NOT NULL,
    longitude          varchar(255) NOT NULL
);

CREATE TABLE user
(
    identifier               varchar(64) PRIMARY KEY NOT NULL DEFAULT (UUID_TO_BIN(UUID()))
);