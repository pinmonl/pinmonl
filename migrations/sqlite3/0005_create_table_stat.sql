-- +migration Up
CREATE TABLE IF NOT EXISTS "stat" (
	"id" VARCHAR(50) NOT NULL PRIMARY KEY,
	"monl_id" VARCHAR(50) NOT NULL,
	"recorded_at" TIMESTAMP NULL,
	"kind" VARCHAR(100) NOT NULL,
	"value" VARCHAR(255) NOT NULL,
	"is_latest" TINYINT,
	"manifest" TEXT NOT NULL
);

-- +migration Down
DROP TABLE IF EXISTS "stat";
