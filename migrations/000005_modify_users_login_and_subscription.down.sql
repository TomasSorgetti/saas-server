ALTER TABLE users
  ADD COLUMN subscription_plan ENUM('free', 'pro', 'enterprise'),
  ADD COLUMN subscription_status ENUM('active', 'pending', 'canceled'),
  DROP COLUMN login_method;