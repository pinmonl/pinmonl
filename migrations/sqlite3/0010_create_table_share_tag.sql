-- +migration Up
CREATE TABLE IF NOT EXISTS "share_tag" (
	"share_id" VARCHAR(50) NOT NULL,
	"tag_id" VARCHAR(50) NOT NULL,
	"kind" VARCHAR(100) NOT NULL
);

-- +migration Down
DROP TABLE IF EXISTS "share_tag";
