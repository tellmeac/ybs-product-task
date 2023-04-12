-- Create "couriers" table
CREATE TABLE "public"."couriers" ("id" serial NOT NULL, "type" character varying(32) NULL, "regions" json NULL, "working_hours" json NULL, PRIMARY KEY ("id"));
-- Create "orders" table
CREATE TABLE "public"."orders" ("id" serial NOT NULL, "weight" integer NULL, "region" integer NULL, "delivery_hours" json NULL, "cost" integer NULL, "completed_time" timestamptz NULL, PRIMARY KEY ("id"));
