-- Checking block
EXPLAIN SELECT * FROM users WHERE role_id = "d8faffca-a00e-11f0-94f9-362064ef513e";
EXPLAIN SELECT * FROM access_tokens WHERE token = "8ef033a06c2035ad5ba7b585918b7455a899e35da086ee0e84c98303d01ba9fc" AND status = true;
SELECT * FROM roles;
SELECT * FROM profiles;
SELECT * FROM role_permissions;
SELECT * FROM access_tokens WHERE status = 1;
SELECT * FROM permissions;
SELECT users.*, profiles.* FROM `users` INNER JOIN profiles ON profiles.user_id = users.id WHERE users.id = '79fb3083-a010-11f0-94f9-362064ef513e';
SELECT * FROM media;
SELECT * FROM posts WHERE id = 2;
SELECT * FROM schema_migrations;
SELECT `id` FROM `roles` WHERE name = 'admin';
SELECT posts.id, posts.content, posts.created_at, posts.updated_at,
media.id AS media_id, media.path AS media_path, media.type AS media_type, media.status AS media_status,
       categories.id AS categories_id, categories.name AS categories_name
FROM `posts`
    LEFT JOIN media ON media.post_id = posts.id
    LEFT JOIN post_categories ON post_categories.post_id = posts.id
    LEFT JOIN categories ON categories.id = post_categories.category_id
WHERE LOWER(posts.content) LIKE LOWER('%%') ORDER BY posts.created_at desc LIMIT 10;

SELECT posts.*, media.* FROM posts INNER JOIN media ON posts.id = media.post_id
                        WHERE posts.content LIKE '%%' ORDER BY posts.created_at DESC LIMIT 10;

SELECT GROUP_CONCAT(DISTINCT media.id ORDER BY media.id ASC SEPARATOR ',') AS media_ids,
       GROUP_CONCAT(DISTINCT media.path ORDER BY media.id ASC SEPARATOR ',') AS media_paths,
       GROUP_CONCAT(DISTINCT media.type ORDER BY media.id ASC SEPARATOR ',') AS media_types,
       GROUP_CONCAT(DISTINCT media.status ORDER BY media.id ASC SEPARATOR ',') AS media_statuses FROM media LEFT JOIN posts ON media.post_id = posts.id WHERE post_id = 6;


SELECT
    posts.id, posts.content, posts.created_at, posts.updated_at,
    GROUP_CONCAT(DISTINCT categories.id ORDER BY categories.id ASC SEPARATOR ',') AS category_ids,
    GROUP_CONCAT(DISTINCT categories.name ORDER BY categories.id ASC SEPARATOR ',') AS category_names,
    GROUP_CONCAT(DISTINCT media.id ORDER BY media.id ASC SEPARATOR ',') AS media_ids,
    GROUP_CONCAT(DISTINCT media.path ORDER BY media.id ASC SEPARATOR ',') AS media_paths,
    GROUP_CONCAT(DISTINCT media.type ORDER BY media.id ASC SEPARATOR ',') AS media_types,
    GROUP_CONCAT(DISTINCT media.status ORDER BY media.id ASC SEPARATOR ',') AS media_statuses FROM `posts`
    LEFT JOIN post_categories ON post_categories.post_id = posts.id
    LEFT JOIN categories ON categories.id = post_categories.category_id
    LEFT JOIN media ON media.post_id = posts.id WHERE posts.id = '6' GROUP BY `posts`.`id`

SELECT * FROM profiles;
SELECT `profile_pic` FROM `profiles` WHERE user_id = 10;
SELECT * FROM notifications;
UPDATE notifications SET noti_type = "personal_post";
SELECT *,
       COUNT(notifications.id) OVER() AS total_count
FROM `notifications` WHERE LOWER(notifications.message) LIKE LOWER('%%') AND receiver_id = 5 ORDER BY notifications.created_at desc LIMIT 10;
