
/*
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
*/

/*
CREATE TABLE users(
   id BINARY(16) PRIMARY KEY,
   username VARCHAR(20),
   email VARCHAR(254),
   phone INT(10),
   password TEXT,
   activate boolean
);
*/

CREATE TABLE index_users_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_username(
    username VARCHAR(20) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (username, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_email(
    email VARCHAR(254) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (email, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_phone(
    phone INT(10),
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (phone, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_password(
    password TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (password, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_activate(
    activate BOOLEAN NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (activate, row_key)
) ENGINE=InnoDB;


/*
CREATE TABLE companies(
   id BINARY(16) PRIMARY KEY,
   category VARCHAR(255),
   domain VARCHAR(63),
   name VARCHAR(255)
);
*/

CREATE TABLE index_companies_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_companies_category(
    category VARCHAR(255) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (category, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_companies_domain(
    domain VARCHAR(63) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (domain, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_companies_name(
    name VARCHAR(255) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (name, row_key)
) ENGINE=InnoDB;


/*
CREATE TABLE schools(
   id BINARY(16) PRIMARY KEY,
   category VARCHAR(255),
   domain VARCHAR(63),
   name VARCHAR(255)
);
*/

CREATE TABLE index_schools_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_schools_category(
    category VARCHAR(255) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (category, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_schools_domain(
    domain VARCHAR(63) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (domain, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_schools_name(
    name VARCHAR(255) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (name, row_key)
) ENGINE=InnoDB;


Note that:
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
