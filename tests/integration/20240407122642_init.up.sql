CREATE TABLE banners
(
  id serial not null unique,
  tag_ids integer[] not null,
  feature_id int not null,
  content_title varchar(255) not null,
  content_text varchar(255) not null,
  content_url varchar(255) not null,
  is_active boolean not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);



