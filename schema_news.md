# Below are index tables for news

```
CREATE TABLE news(
    id BINARY(16) PRIMARY KEY,
    domain VARCHAR(63) NOT NULL,
    timestamp BIGINT NOT NULL,
    author VARCHAR(255) NOT NULL,
    title VARCHAR(1024) NOT NULL, 
);

CREATE TABLE index_news_domain(
    domain VARCHAR(63) NOT NULL, 
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (domain, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_news_timestamp(
    timestamp BIGINT NOT NULL, 
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (timestamp, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_news_author(
    author VARCHAR(255) NOT NULL, 
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (author, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_news_title(
    title VARCHAR(1024) NOT NULL, 
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (title, row_key)
) ENGINE=InnoDB;
```
* Note that CREATE TABLE news() would not be used in the database, it would be a virtual table for reference only, all index tables are in database.
* Note that index_news_domain's (domain VARCHAR) can be school, company domain or topics (e.g. world cup, NBA) 
* all of these should be unique
