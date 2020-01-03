-- +migration Up
CREATE TABLE IF NOT EXISTS "pinmon" (
	"pinl_id" VARCHAR(50) NOT NULL,
	"monl_id" VARCHAR(50) NOT NULL,
	"user_id" VARCHAR(50) NOT NULL,
	"sort" INTEGER
);

-- +migration Down
DROP TABLE IF EXISTS "pinmon";
