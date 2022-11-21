CREATE TABLE IF NOT EXISTS extract (
	id VARCHAR(40) PRIMARY KEY,
	fileName TEXT NOT NULL,
	dir TEXT NOT NULL,
	status TEXT NOT NULL,
	aligorithm TEXT NULL,
	filters TEXT NULL,
	partialExtraction TEXT NOT NULL CHECK (partialExtraction IN (0, 1)),
	background BOOLEAN NOT NULL CHECK (background IN (0, 1)),
	createdOn TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);