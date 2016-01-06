package storages

const CREATE_TABLES_SQL = `
--Table that stores the url to short_url mapping
CREATE TABLE urls (
	id SERIAL PRIMARY KEY NOT NULL,
	created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (statement_timestamp() AT TIME ZONE 'utc'),
	url VARCHAR(512) NOT NULL,
	token VARCHAR(64) NOT NULL
);

CREATE UNIQUE INDEX index_urls_token ON urls(token);

--Table that stores all accesses to a short_url
CREATE TABLE stats (
	id SERIAL PRIMARY KEY NOT NULL,
	created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (statement_timestamp() AT TIME ZONE 'utc'),
	token VARCHAR(64) NOT NULL
);
CREATE INDEX index_stats_token ON stats(token);
`

const DROP_TABLES_SQL = `
DROP TABLE IF EXISTS urls;
DROP TABLE IF EXISTS stats;
`

const INSERT_URL_SQL = `INSERT INTO urls (url, token) VALUES ($1, $2)`
const SELECT_URL_SQL = `SELECT url FROM urls WHERE token=$1`

const INSERT_STATS_SQL = `INSERT INTO stats (token) VALUES ($1)`
const SELECT_STATS_SQL = `SELECT count(*) as hits FROM stats WHERE token=$1`
