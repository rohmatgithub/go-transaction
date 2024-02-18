-- +migrate Up

-- +migrate StatementBegin


CREATE TABLE IF NOT EXISTS sales_order (
	id              SERIAL PRIMARY KEY,
	company_id      INTEGER NOT NULL,
    order_number	VARCHAR(20) DEFAULT '',
    order_date		DATE,
    customer_id		INTEGER	DEFAULT 0,
    total_gross_amount		FLOAT	DEFAULT 0.0,
    total_net_amount		FLOAT	DEFAULT 0.0,
	created_by      INTEGER NULL DEFAULT 0,
	created_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by      INTEGER NULL DEFAULT 0,
	updated_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT uq_so_primary_key UNIQUE (company_id, order_number)
);

CREATE TABLE IF NOT EXISTS sales_order_item (
	id              SERIAL PRIMARY KEY,
	order_id		INTEGER NOT NULL,
	product_id     	INTEGER NOT NULL,
	qty				INTEGER DEFAULT 0,
    selling_price	FLOAT	DEFAULT 0.0,
    line_gross_amount	FLOAT	DEFAULT 0.0,
    line_net_amount		FLOAT	DEFAULT 0.0,
	created_by      INTEGER NULL DEFAULT 0,
	created_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by      INTEGER NULL DEFAULT 0,
	updated_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT uq_soi_primary_key UNIQUE (order_id, product_id)
);

CREATE TABLE IF NOT EXISTS sales_invoice (
	id              SERIAL PRIMARY KEY,
	company_id      INTEGER NOT NULL,
    invoice_number	VARCHAR(20) DEFAULT '',
    invoice_date		DATE,
    customer_id		INTEGER	DEFAULT 0,
    order_id		INTEGER	DEFAULT 0,
    total_gross_amount		FLOAT	DEFAULT 0.0,
    total_net_amount		FLOAT	DEFAULT 0.0,
	created_by      INTEGER NULL DEFAULT 0,
	created_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by      INTEGER NULL DEFAULT 0,
	updated_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT uq_si_primary_key UNIQUE (company_id, invoice_number)
);

CREATE TABLE IF NOT EXISTS sales_invoice_item (
	id              SERIAL PRIMARY KEY,
	invoice_id		INTEGER NOT NULL,
	product_id     	INTEGER NOT NULL,
	qty				INTEGER DEFAULT 0,
    selling_price	FLOAT	DEFAULT 0.0,
    line_gross_amount	FLOAT	DEFAULT 0.0,
    line_net_amount		FLOAT	DEFAULT 0.0,
	created_by      INTEGER NULL DEFAULT 0,
	created_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	updated_by      INTEGER NULL DEFAULT 0,
	updated_at      timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT uq_sii_primary_key UNIQUE (invoice_id, product_id)
);

-- +migrate StatementEnd