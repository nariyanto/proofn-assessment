DROP TABLE IF EXISTS "users";
CREATE TABLE "users" (
  "id" int4 NOT NULL DEFAULT nextval('user_id_seq'::regclass),
  "email" varchar(255) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "password" varchar(255) COLLATE "pg_catalog"."default" DEFAULT NULL::character varying,
  "token" text COLLATE "pg_catalog"."default" DEFAULT NULL,
  "token_verification" varchar(10) COLLATE "pg_catalog"."default" DEFAULT NULL,
  "status" int2 DEFAULT NULL
)
;
ALTER TABLE "users" OWNER TO "postgres";

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO "users" VALUES (13, 'reza.kaseptea@gmail.com', '$2a$14$iarifcwPrfHf9D1gpqpneuQ3hLeGZCxN30lUlgSvFm2KkpVrgU9Hi', 'cNSpRsjuNS', 'TmTKAdCSYS', 1);
INSERT INTO "users" VALUES (10, 'rezairwantoo@gmail.com', '$2a$14$iarifcwPrfHf9D1gpqpneuQ3hLeGZCxN30lUlgSvFm2KkpVrgU9Hi', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMCwicGFzc3dvcmQiOiIkMmEkMTQkaWFyaWZjd1ByZkhmOUQxZ3BxcG5ldVEzaExlR1pDeE4zMGxVbGdTdkZtMktrcFZyZ1U5SGkiLCJlbWFpbCI6InJlemFpcndhbnRvb0BnbWFpbC5jb20iLCJlbXBsb3llZV9pZCI6MCwiZXhwIjoxNTg2NDk5OTg5fQ.EDRvgPFa--T9QRpOLj6tJZpgvebZHKmK2s_y13ZaHKI', 'TmTKAdCSYW', 1);
COMMIT;

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");