PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;

CREATE TABLE META_INFO (
	FIELD	TEXT PRIMARY KEY,
	VALUE	TEXT
);
INSERT INTO META_INFO VALUES('schema_version','1.1.1');

CREATE TABLE META_CHANGELOG (
	ID			INTEGER PRIMARY KEY,
	PUB_DATE	TEXT DEFAULT CURRENT_DATE,
	VERSION		TEXT,
	TITLE		TEXT,
	MARKDOWN	TEXT,
	RENDERED	TEXT
);
INSERT INTO META_CHANGELOG (PUB_DATE, VERSION, TITLE, MARKDOWN) VALUES
	( '2022-05-26', '1.2.0', 'Minor Feature Update', 'Added a couple convenience features, and no, I haven''t added image capability yet. that brings a lot of technical challenges with it.\n - "Keep me logged in" should now actually work, at least up to a month.\n - Item status is now shown on the list page instead of the description.\n' ),
	( '2022-05-21', '1.1.0', 'Security Update!', 'I''ve managed to actually get around to updating this. Sadly, I don''t have any really exciting new features to report, but the security between the website and API has been improved significantly.\n - Uses JWT instead of (incorrect) Oauth2 for API authentication\n - SSL Certificates are dynamically mapped so it won''t just stop working once they expire\n - Code was slightly cleaned up everywhere' ),
	( '2021-11-21', '1.0.0', 'First Full Release!', 'Today, the 21st of November 2021, I''ve finally gotten the whole app to a point where I''m happy enough with how finished it is. Below is a list of all the features available in the first version. Thank you for your patience and please enjoy:\n - Creating user accounts and wishlists\n - Adding/deleting items from your wishlist\n - Reserving other users'' items for purchase and\n - Changing your password' )
;

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
