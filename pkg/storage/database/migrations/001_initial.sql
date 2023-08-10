-- +migrate Up
CREATE TABLE entries (
	hash varchar(64) not null,
	source_name varchar(255) not null,
	text varchar(255) not null,
	url varchar not null,
	image_url varchar,
	categories varchar,
  status char(1) not null,
	created_at datetime not null,
  ttl datetime not null,
	CONSTRAINT hash PRIMARY KEY (hash)
);

CREATE INDEX entries_source_name_IDX ON entries (source_name);
CREATE INDEX entries_ttl_IDX ON entries (ttl);
CREATE INDEX entries_created_at_IDX ON entries (created_at);

-- +migrate Down
DROP TABLE entries;
