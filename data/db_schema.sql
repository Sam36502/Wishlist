PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE META_INFO (FIELD TEXT PRIMARY KEY, VALUE TEXT);
INSERT INTO META_INFO VALUES('schema_version','1.1.0');
INSERT INTO META_INFO VALUES('created_date','2025-11-11');
CREATE TABLE IF NOT EXISTS "tbl_user" (
  "id_user" INTEGER NOT NULL PRIMARY KEY,
  "email" TEXT NOT NULL UNIQUE,
  "password" TEXT NOT NULL,
  "name" TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS "tbl_status" (
  "id_status" INTEGER NOT NULL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "desc" text
);
INSERT INTO tbl_status VALUES(1,'Available','No one is planning to get this yet.');
INSERT INTO tbl_status VALUES(2,'Reserved','Someone is planning to get this already.');
INSERT INTO tbl_status VALUES(3,'Received','This was on the wishlist, but has since been received.');
CREATE TABLE IF NOT EXISTS "tbl_item" (
  "id_item" INTEGER NOT NULL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "desc" TEXT NOT NULL,
  "status_id" INTEGER NOT NULL REFERENCES "tbl_status" ("id_status") ON DELETE CASCADE ON UPDATE CASCADE,
  "user_id" INTEGER NOT NULL REFERENCES "tbl_user" ("id_user") ON DELETE CASCADE ON UPDATE CASCADE,
  "reserved_by_user_id" INTEGER DEFAULT NULL REFERENCES "tbl_user" ("user_id") ON DELETE CASCADE ON UPDATE CASCADE,
  "price" INTEGER NOT NULL DEFAULT '0'
);
CREATE TABLE IF NOT EXISTS "tbl_link" (
  "id_link" INTEGER NOT NULL PRIMARY KEY,
  "text" TEXT NOT NULL,
  "hyperlink" TEXT NOT NULL,
  "item_id" INTEGER NOT NULL REFERENCES "tbl_item" ("id_item") ON DELETE CASCADE ON UPDATE CASCADE
);
COMMIT;
