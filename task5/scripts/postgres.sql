DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
    "id" bigserial primary key,
    "name" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL,
    "password" varchar(255) NOT NULL,
    "token" text NOT NULL,
    "token_verification" varchar(20) DEFAULT NULL,
    "status" int2 DEFAULT 0,
    "created_date" timestamp default now()
);
ALTER TABLE "users" OWNER TO "postgres";

DROP TABLE IF EXISTS "messages";
CREATE TABLE "messages" (
  "id" bigserial primary key,
  "sender_id" bigint NOT NULL,
  "receiver_id" bigint NOT NULL,
  "content" text NOT NULL,
  "created_date" timestamp default now(),
  "receiver_status" int2 DEFAULT 0,
  "sender_status" int2 DEFAULT 0
)
;
ALTER TABLE "messages" OWNER TO "postgres";