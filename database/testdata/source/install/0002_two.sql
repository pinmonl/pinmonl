-- +migration Up
CREATE TABLE two (id INTEGER PRIMARY KEY);
-- +migration Down
DROP TABLE two;