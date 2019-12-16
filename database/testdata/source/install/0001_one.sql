-- +migration Up
CREATE TABLE one (id INTEGER PRIMARY KEY);
-- +migration Down
DROP TABLE one;