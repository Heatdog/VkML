CREATE TABLE IF NOT EXISTS documents(
    url TEXT NOT NULL,
    pub_date INT8 NOT NULL,
    fetch_time INT8 NOT NULL,
    text TEXT NOT NULL,
    CONSTRAINT url_fetch_time_documents_pk PRIMARY KEY (url, fetch_time)
);

CREATE INDEX documents_fetch_time ON documents(fetch_time);