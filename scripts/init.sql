CREATE TABLE IF NOT EXISTS urls (
    id BIGINT NOT NULL AUTO_INCREMENT,
    long_url TEXT NOT NULL,
    short_url VARCHAR(255) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (short_url)
);
