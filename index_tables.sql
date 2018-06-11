
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
   username TEXT,
   email TEXT,
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
    username TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (username, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_email(
    email TEXT NOT NULL,
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
   category TEXT,
   domain TEXT,
   name TEXT
);
*/

CREATE TABLE index_companies_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_companies_category(
    category TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (category, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_companies_domain(
    domain TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (domain, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_companies_name(
    name TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (name, row_key)
) ENGINE=InnoDB;


/*
CREATE TABLE schools(
   id BINARY(16) PRIMARY KEY,
   category TEXT,
   domain TEXT,
   name TEXT
);
*/

CREATE TABLE index_schools_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_schools_category(
    category TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (category, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_schools_domain(
    domain TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (domain, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_schools_name(
    name TEXT NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (name, row_key)
) ENGINE=InnoDB;
