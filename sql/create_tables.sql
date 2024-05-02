create table if not exists apps_activityhistory
(
    timestamps    timestamp,
    location      varchar(255),
    latitude      varchar(255),
    longitude     varchar(255),
    status        varchar(255),
    user_id       integer,
    message       varchar(255),
    status_report varchar(255),
    id            serial
        primary key
);

create table if not exists apps_activityhistory_rescuers
(
    id                    serial
        primary key,
    activityhistory_id    integer,
    rescueteamcontacts_id integer
);

create table if not exists apps_rescueteamcontacts
(
    id      serial
        primary key,
    name    varchar(255),
    address varchar(255),
    contact varchar(255)
);

create table if not exists apps_trustedcontacts
(
    id       serial
        primary key,
    name     varchar(255),
    address  varchar(255),
    contact  varchar(255),
    owner_id integer
);

create table if not exists apps_user
(
    id                 serial
        primary key,
    password           varchar(255),
    username           varchar(255),
    email              varchar(255) unique ,
    first_name         varchar(255),
    last_name          varchar(255),
    address            varchar(255),
    contact            varchar(255),
    device_id          varchar(255),
    profile_picture    varchar(255),
    role    varchar(255),
    is_onboarding_done boolean,
    date_joined        timestamp
);

create table if not exists apps_vehicle
(
    id         serial
        primary key,
    brand      varchar(255),
    model      varchar(255),
    year_model integer,
    plate_no   varchar(255),
    owner_id   integer
);

create table if not exists apps_accidentalert
(
    id        serial
        primary key,
    latitude  varchar(255) not null,
    longitude varchar(255) not null,
    is_active boolean      not null,
    device_id varchar(255) not null
);

create table if not exists apps_accident_rescuer
(
    id                    serial
        primary key,
    activity_history_id   integer       not null,
    rescueteamcontacts_id integer       not null,
    responders_name       varchar(255)  not null,
    notes                 varchar(1000) not null
);