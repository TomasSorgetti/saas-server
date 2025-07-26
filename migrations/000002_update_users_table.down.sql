ALTER TABLE users
  MODIFY COLUMN role ENUM('luthier', 'admin', 'superadmin') NOT NULL;

ALTER TABLE users
  ADD COLUMN reset_password_token VARCHAR(255),
  ADD COLUMN reset_password_expires TIMESTAMP NULL;

DROP TABLE IF EXISTS password_resets;
DROP TABLE IF EXISTS email_verifications;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS subscription_plans;