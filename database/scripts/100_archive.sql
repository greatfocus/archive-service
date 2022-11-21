CREATE TABLE IF NOT EXISTS archive (
	id VARCHAR(40) PRIMARY KEY,
	fileName TEXT NOT NULL,
	dir TEXT NOT NULL,
	status TEXT NOT NULL,
	aligorithm TEXT NULL,
	filters TEXT NULL,
	background BOOLEAN NOT NULL CHECK (background IN (0, 1)),
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(fileName)
);