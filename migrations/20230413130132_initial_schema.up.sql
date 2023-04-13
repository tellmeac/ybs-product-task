-- create "couriers" table
CREATE TABLE "public"."couriers" ("id" serial NOT NULL, "courier_type" character varying(32) NOT NULL, "regions" json NOT NULL, "working_hours" json NOT NULL, PRIMARY KEY ("id"));
-- create "orders" table
CREATE TABLE "public"."orders" ("id" serial NOT NULL, "weight" numeric NOT NULL, "region" integer NOT NULL, "delivery_hours" json NOT NULL, "cost" integer NOT NULL, "completed_time" timestamptz NULL, PRIMARY KEY ("id"));
