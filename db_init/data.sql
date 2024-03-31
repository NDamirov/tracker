CREATE TABLE tasks (
  task_id bigserial primary key,
  author_id bigint not null,
  task_description text,
  task_status varchar(50),
  created_at bigint
);
