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

CREATE TABLE index_users_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;

CREATE TABLE index_users_id(
    id BINARY(16) NOT NULL,
    row_key BINARY(16) NOT NULL UNIQUE,
    PRIMARY KEY (id, row_key)
) ENGINE=InnoDB;
/*
CREATE TABLE companies(
   id BINARY(16) PRIMARY KEY,
   category TEXT,
   domain TEXT,
   name TEXT
);
*/


/*
CREATE TABLE schools(
   id BINARY(16) PRIMARY KEY,
   category TEXT,
   domain TEXT,
   name TEXT
);
*/
