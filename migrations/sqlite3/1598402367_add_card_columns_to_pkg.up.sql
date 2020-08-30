ALTER TABLE pkgs
ADD COLUMN title VARCHAR(250);

ALTER TABLE pkgs
ADD COLUMN description TEXT;

ALTER TABLE pkgs
ADD COLUMN image_id VARCHAR(50);

ALTER TABLE pkgs
ADD COLUMN custom_uri VARCHAR(1000);

UPDATE pkgs SET 
  title = '',
  description = '',
  image_id = '',
  custom_uri = '';
