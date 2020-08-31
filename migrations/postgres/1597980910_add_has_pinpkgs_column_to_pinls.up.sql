ALTER TABLE pinls
ADD COLUMN has_pinpkgs BOOLEAN;

UPDATE pinls
  SET has_pinpkgs = false;
