CREATE DATABASE calendardb;
CREATE ROLE db_user WITH LOGIN PASSWORD 'super_secret_password_42';
GRANT ALL PRIVILEGES ON DATABASE calendardb TO db_user;

\c calendardb;

CREATE TABLE events (
  id varchar(50) NOT NULL,
  title varchar(20) NOT NULL,
  description text NULL,
  owner_id varchar(50) NOT NULL,
  start_date timestamp,
  finish_date timestamp,
  notification_time timestamp
);

GRANT ALL PRIVILEGES ON SCHEMA public TO db_user;
GRANT ALL PRIVILEGES ON TABLE events TO db_user;
