DROP TABLE IF EXISTS cell;

SHOW WARNINGS;

CREATE TABLE cell
(
	added_at      BIGINT PRIMARY KEY AUTO_INCREMENT,
	row_key       BINARY(16) NOT NULL,
	column_name	  VARCHAR(64) NOT NULL,
	ref_key		    BIGINT NOT NULL,
	body		      BLOB,
	created_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
	UNIQUE `cell_idx`(`row_key`, `column_name`, `ref_key`)
) ENGINE=InnoDB;

SHOW WARNINGS;
