create table if not exists permissions
(
    group        String,
    code         String,
    description  String,
    created_time DATETIME('UTC'),
    created_by   UUID,
    updated_time DATETIME('UTC'),
    updated_by   UUID
)
    ENGINE = MergeTree
        PRIMARY KEY (code);


create table if not exists organize
(
    uuid          UUID,
    name          String,
    contact       String,
    website       String,
    status        String,
    verify_status String,
    created_time  DATETIME('UTC'),
    created_by    UUID,
    updated_time  DATETIME('UTC'),
    updated_by    UUID
)
    ENGINE = MergeTree
        PRIMARY KEY (uuid);

create table if not exists application
(
    organize_uuid   UUID,
    app_uuid        UUID,
    app_name        String,
    app_description String,
    default_app     FixedString(1),
    status          String,
    created_time    DATETIME('UTC'),
    created_by      UUID,
    updated_time    DATETIME('UTC'),
    updated_by      UUID
)
    ENGINE = MergeTree
        PARTITION BY (organize_uuid)
        PRIMARY KEY (organize_uuid, app_uuid);

create table if not exists api_keys
(
    app_uuid                UUID,
    key_uuid                UUID,
    access_key              String,
    access_key_expire_time  DATETIME('UTC'),
    refresh_key             String,
    refresh_key_expire_time DATETIME('UTC'),
    permissions             Array(String),
    status                  String,
    created_time            DATETIME('UTC'),
    created_by              UUID,
    updated_time            DATETIME('UTC'),
    updated_by              UUID
)
    ENGINE = MergeTree
        PARTITION BY (app_uuid)
        PRIMARY KEY (app_uuid, key_uuid);

create table users
(
    uuid          UUID,
    email         varchar(255),
    secret        String,
    vital_status  String,
    verify_status String,
    created_time  DATETIME('UTC'),
    created_by    UUID,
    updated_time  DATETIME('UTC'),
    updated_by    UUID
)
    engine = MergeTree
        PRIMARY KEY (uuid);

create table user_information
(
    user_uuid    UUID,
    language     varchar(5),
    tag_name     String,
    tag_value    String,
    created_time DATETIME('UTC'),
    created_by   UUID,
    updated_time DATETIME('UTC'),
    updated_by   UUID
)
    ENGINE = MergeTree
        PRIMARY KEY (user_uuid, language, tag_name);

create table if not exists user_sessions
(
    user_uuid    UUID,
    mac_address  String,
    user_agent   String,
    created_time DATETIME('UTC'),
    created_by   UUID
)
    ENGINE = MergeTree
        PARTITION BY (user_uuid)
        PRIMARY KEY (user_uuid, mac_address, user_agent);
