-- +migration Up
CREATE TABLE IF NOT EXISTS "share" (
	"id" VARCHAR(50) NOT NULL,
	"user_id" VARCHAR(50) NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	"description" TEXT NOT NULL,
	"image_id" VARCHAR(50) NOT NULL,
	"created_at" TIMESTAMP NULL,
	"updated_at" TIMESTAMP NULL
);

-- +migration Down
DROP TABLE IF EXISTS "share";
