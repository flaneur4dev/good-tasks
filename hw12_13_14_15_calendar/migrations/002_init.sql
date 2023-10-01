\c calendardb db_user;

INSERT INTO events (id, title, description, owner_id, start_date, finish_date, notification_time)
VALUES
  ('1', 'event 1', 'description 1', '11', '2023-09-01 04:05:00', '2023-09-01 14:00:00', '2023-09-01 13:00:00'),
  ('2', 'event 2', 'description 2', '22', '2023-09-04 10:05:00', '2023-09-05 10:05:00', '2023-09-05 08:00:00'),
  ('3', 'event 3', 'description 3', '33', '2023-08-28 06:30:00', '2023-08-28 18:00:00', '2023-08-28 17:30:00'),
  ('4', 'event 4', 'description 4', '44', '2023-09-15 01:15:00', '2023-09-16 11:00:00', '2023-09-16 07:00:00');
