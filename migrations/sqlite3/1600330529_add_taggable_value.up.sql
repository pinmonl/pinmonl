ALTER TABLE taggables
ADD COLUMN value VARCHAR(250);

ALTER TABLE taggables
ADD COLUMN value_type INTEGER;

ALTER TABLE taggables
ADD COLUMN value_prefix VARCHAR(200);

ALTER TABLE taggables
ADD COLUMN value_suffix VARCHAR(200);

UPDATE taggables SET
  value = '',
  value_type = 0,
  value_prefix = '',
  value_suffix = '';
