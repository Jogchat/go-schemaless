
# Set database user name and privileges

```
CREATE USER 'root'@'%' IDENTIFIED BY 'Umiuni_jogchat_schemales_2018@';
GRANT ALL ON *.* TO 'root'@'%';
SET PASSWORD FOR 'root'@'localhost' = PASSWORD('Umiuni_jogchat_schemales_2018@');
FLUSH PRIVILEGES;
```

# Below are schemaless core level table
## cell is a schemaless cell that can store any nosql blob (schema flexibility)

```
CREATE TABLE cell
(
    added_at         BIGINT PRIMARY KEY AUTO_INCREMENT,
    row_key          BINARY(16) NOT NULL,
    column_name      VARCHAR(64) NOT NULL,
    ref_key          BIGINT NOT NULL,
    body             BLOB,
    created_at       DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT cell_idx UNIQUE(row_key, column_name,ref_key)
) ENGINE=InnoDB;
```


# Below are application level schema tables

Schema entities for users, companies and schools:
* https://github.com/Jogchat/go-schemaless/blob/master/schema_entities.md

## Below are index tables for news
```
CREATE TABLE index_news_domain(
    domain VARCHAR(64) NOT NULL, 
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

## Note that:
* Every email address is composed of two parts. The local part that comes before the '@' sign, and the domain part that follows it. In "user@example.com", the local part is "user", and the domain part is "example.com".

The local part must not exceed 64 characters and the domain part cannot be longer than 255 characters.
https://stackoverflow.com/questions/386294/what-is-the-maximum-length-of-a-valid-email-address

The combined length of the local + @ + domain parts of an email address must not exceed 254 characters. As described in RFC3696 Errata ID 1690.
* Username length, choose 20. Twitter uses 15, pearson uses 32, blind use 10.
http://help.pearsoncmg.com/rumba/b2c_self_reg/en/Content/b2c_signin_guidelines.html
https://help.twitter.com/en/managing-your-account#username-email-and-phone
* https://stackoverflow.com/questions/1885630/whats-the-difference-between-varchar-and-char
* label part 63 characters max: https://en.wikipedia.org/wiki/Domain_Name_System  
https://stackoverflow.com/questions/14402407/maximum-length-of-a-domain-name-without-the-http-www-com-parts
* Bcrypt hash length: https://stackoverflow.com/questions/5881169/what-column-type-length-should-i-use-for-storing-a-bcrypt-hashed-password-in-a-d
