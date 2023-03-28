-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS provider_signature_id_seq;
CREATE TABLE "public"."provider_signature" (
    "id" int4 NOT NULL DEFAULT nextval('provider_signature_id_seq'::regclass),
    "name" varchar(40) NOT NULL,
    "type" varchar(10) NOT NULL,
    "source" varchar,
    "signature" varchar,
    "created" timestamp NOT NULL DEFAULT now(),
    "updated" timestamp,
    PRIMARY KEY ("id")
);

CREATE SEQUENCE IF NOT EXISTS signature_botnet_id_seq;
CREATE TABLE "public"."signature_botnet" (
    "id" int4 NOT NULL DEFAULT nextval('signature_botnet_id_seq'::regclass),
    "ip" inet NOT NULL,
    "port" int4 NOT NULL,
    "provider" int4 NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "public"."provider_signature"("id")
);

CREATE SEQUENCE IF NOT EXISTS signature_cert_id_seq;
CREATE TABLE "public"."signature_cert" (
    "id" int4 NOT NULL DEFAULT nextval('signature_cert_id_seq'::regclass),
    "sha1" varchar(40) NOT NULL,
    "name" varchar(40) NOT NULL,
    "provider" int4,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "public"."provider_signature"("id")
);

CREATE SEQUENCE IF NOT EXISTS untitled_table_209_id_seq;
CREATE TABLE "public"."signature_ip" (
    "id" int4 NOT NULL DEFAULT nextval('untitled_table_209_id_seq'::regclass),
    "ip" inet NOT NULL,
    "provider" int4 NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "public"."provider_signature"("id")
);

CREATE SEQUENCE IF NOT EXISTS signature_tracker_id_seq;
CREATE TABLE "public"."signature_tracker" (
    "id" int4 NOT NULL DEFAULT nextval('signature_tracker_id_seq'::regclass),
    "url" varchar NOT NULL,
    "provider" int4 NOT NULL,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("provider") REFERENCES "public"."provider_signature"("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "public"."signature_ip";
DROP TABLE "public"."signature_cert";
DROP TABLE "public"."signature_botnet";
DROP TABLE "public"."signature_tracker";
DROP TABLE "public"."provider_signature";
-- +goose StatementEnd
