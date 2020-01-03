-- +migration Up
CREATE TABLE IF NOT EXISTS "taggable" (
	"tag_id" VARCHAR(50) NOT NULL,
	"target_id" VARCHAR(50) NOT NULL,
	"target_name" VARCHAR(100) NOT NULL,
	"sort" INTEGER
);

-- +migration Down
DROP TABLE IF EXISTS "taggable";
