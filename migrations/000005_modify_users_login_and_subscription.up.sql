ALTER TABLE users
  ADD COLUMN login_method VARCHAR(50) DEFAULT NULL,
  DROP COLUMN subscription_plan,
  DROP COLUMN subscription_status;