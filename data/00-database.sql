-- DATABASE STRUCTURE SCRIPT

--DROP DATABASE IF EXISTS "wishlist";
--CREATE DATABASE "wishlist" CHARACTER SET utf8;

--GRANT INSERT, SELECT, UPDATE, DELETE ON "wishlist".* TO "wishlist_user";
--USE "wishlist";

DROP TABLE IF EXISTS "tbl_user";
CREATE TABLE "tbl_user" (
  "id_user" INTEGER NOT NULL PRIMARY KEY,
  "email" TEXT NOT NULL UNIQUE,
  "password" TEXT NOT NULL,
  "name" TEXT NOT NULL
);

DROP TABLE IF EXISTS "tbl_status";
CREATE TABLE "tbl_status" (
  "id_status" INTEGER NOT NULL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "desc" text
);

INSERT INTO "tbl_status" VALUES (1,'Available','No one is planning to get this, yet.'),(2,'Reserved','Someone is planning to get this already.'),(3,'Received','This was on the wishlist, but has since been received.');

DROP TABLE IF EXISTS "tbl_item";
CREATE TABLE "tbl_item" (
  "id_item" INTEGER NOT NULL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "desc" TEXT NOT NULL,
  "status_id" INTEGER NOT NULL REFERENCES "tbl_status" ("id_status") ON DELETE CASCADE ON UPDATE CASCADE,
  "user_id" INTEGER NOT NULL REFERENCES "tbl_user" ("id_user") ON DELETE CASCADE ON UPDATE CASCADE,
  "reserved_by_user_id" INTEGER DEFAULT NULL REFERENCES "tbl_user" ("user_id") ON DELETE CASCADE ON UPDATE CASCADE,
  "price" INTEGER NOT NULL DEFAULT '0'
);

DROP TABLE IF EXISTS "tbl_link";
CREATE TABLE "tbl_link" (
  "id_link" INTEGER NOT NULL PRIMARY KEY,
  "text" TEXT NOT NULL,
  "hyperlink" TEXT NOT NULL,
  "item_id" INTEGER NOT NULL REFERENCES "tbl_item" ("id_item") ON DELETE CASCADE ON UPDATE CASCADE
);
