# Below are index tables for news
```
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
    author VARCHAR(256) NOT NULL, 
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (author, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_news_title(
    title VARCHAR(1024) NOT NULL, 
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (title, row_key)
) ENGINE=InnoDB;
```

* Note that index_news_domain's (domain VARCHAR) can be school, company domain or topics (e.g. world cup, NBA) 
* all of these should be unique
