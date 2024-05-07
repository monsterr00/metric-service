CREATE TABLE IF NOT EXISTS metrics (
	ID varchar(255),
	MType varchar(255),
	Delta bigint,
	Value double precision,
PRIMARY KEY (ID, MType));

          