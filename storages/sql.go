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
`

const DROP_TABLES_SQL = `
DROP TABLE IF EXISTS urls;
`

const INSERT_URL_SQL = `INSERT INTO urls (url, token) VALUES ($1, $2)`

const SELECT_URL_SQL = `SELECT url FROM urls WHERE token=$1`
