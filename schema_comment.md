# Table for comment

```
CREATE TABLE comment(
    id VARCHAR(36) PRIMARY KEY,
    newsId BINARY(16) NOT NULL,
    content VARCHAR(300) NOT NULL,
    timestamp BIGINT NOT NULL,
    parentCommentId BINARY(16) NOT NULL,
);
```
Note that summary, originalUrl should not be indexed
# Below are index tables for comment

```
CREATE TABLE index_comment_id(
    id VARCHAR(36) NOT NULL, 
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_comment_newsId(
    newsId VARCHAR(36) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (newsId, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_comment_content(
    content VARCHAR(300) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (content, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_comment_timestamp(
    timestamp BIGINT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (timestamp, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_comment_parentCommentId(
    parentCommentId BIGINT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE, 
    PRIMARY KEY (parentCommentId, row_key)
) ENGINE=InnoDB;
```

* Note that all should be indexed including content, to prevent hackers keep saying spam words (e.g. thousands of fuck)
