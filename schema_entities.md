Schema entities for users, companies and schools

users is a table storing all jogchat users

```
CREATE TABLE users(
   id BINARY(16) PRIMARY KEY,
   username VARCHAR(20),
   email VARCHAR(254),
   phone INT(10),
   password BINARY(60),
   token BINARY(60),
   activate boolean
);

Below are index tables for users

```
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
    password BINARY(60) NOT NULL,
    row_key  BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (password, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_activate(
    activate BOOLEAN NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (activate, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_token(
    token BINARY(60) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (token, row_key)
) ENGINE=InnoDB;
```

Below are company table and index tables for companies

```
CREATE TABLE companies(
   id BINARY(16) PRIMARY KEY,
   category VARCHAR(255),
   domain VARCHAR(63),
   name VARCHAR(255)
);

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
```

Below are schools table and index tables for schools

```
CREATE TABLE schools(
   id BINARY(16) PRIMARY KEY,
   category VARCHAR(255),
   domain VARCHAR(63),
   name VARCHAR(255)
);

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
```

