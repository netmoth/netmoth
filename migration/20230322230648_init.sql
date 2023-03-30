-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS provider_signature_id_seq;
CREATE TABLE "provider_signature" (
    "id" int4 NOT NULL DEFAULT nextval('provider_signature_id_seq'::regclass),
    "name" varchar(64) NOT NULL,
    "type" varchar(10) NOT NULL,
    "source" varchar(255),
    "signature" varchar(255),
    "active" bool NOT NULL,
    "created" timestamp(6) NULL DEFAULT CURRENT_TIMESTAMP(6),
    "updated" timestamp(6) NULL DEFAULT NULL,
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS signature_botnet_id_seq;
CREATE TABLE "signature_botnet" (
    "id" int4 NOT NULL DEFAULT nextval('signature_botnet_id_seq'::regclass),
    "ip" inet NOT NULL,
    "port" int4 NOT NULL,
    "provider" int4 NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "provider_signature"("id")
);

CREATE SEQUENCE IF NOT EXISTS signature_cert_id_seq;
CREATE TABLE "signature_cert" (
    "id" int4 NOT NULL DEFAULT nextval('signature_cert_id_seq'::regclass),
    "sha1" varchar(40) NOT NULL,
    "name" varchar(64) NOT NULL,
    "provider" int4,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "provider_signature"("id")
);

CREATE SEQUENCE IF NOT EXISTS untitled_table_209_id_seq;
CREATE TABLE "signature_ip" (
    "id" int4 NOT NULL DEFAULT nextval('untitled_table_209_id_seq'::regclass),
    "ip" inet NOT NULL,
    "provider" int4 NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "provider_signature"("id")
);

CREATE SEQUENCE IF NOT EXISTS signature_tracker_id_seq;
CREATE TABLE "signature_tracker" (
    "id" int4 NOT NULL DEFAULT nextval('signature_tracker_id_seq'::regclass),
    "url" varchar(255) NOT NULL,
    "provider" int4 NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "provider_signature"("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "signature_ip";
DROP TABLE "signature_cert";
DROP TABLE "signature_botnet";
DROP TABLE "signature_tracker";
DROP TABLE "provider_signature";
-- +goose StatementEnd
