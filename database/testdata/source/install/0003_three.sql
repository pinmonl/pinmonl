-- +migration Up
CREATE TABLE three (id INTEGER PRIMARY KEY);
-- +migration Down
DROP TABLE three;