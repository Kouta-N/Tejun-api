CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `email` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `password` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
  `profile_image` varchar(255) COLLATE utf8mb4_general_ci DEFAULT '',
  `email_verified` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_email` (`email`)
);

CREATE TABLE `todos` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `content` json NOT NULL,
  `is_done` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
  -- KEY `fk_todo_user` (`user_id`),
  -- CONSTRAINT `fk_todo_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
);

INSERT INTO users (`id`,`email`,`password`,`name`,`created_at`) VALUES (1,'ichiro.com','12345678','一郎',now()), (2,'jiro.com','12345678','二郎',now());

INSERT INTO todos (`user_id`,`content`,`created_at`) VALUES (1, JSON_ARRAY('起きる', '歯を磨く', '顔を洗う', 'カレーを食べる'),now()), (2, JSON_ARRAY('本を開く', '読む', '栞を挟む', '本棚にしまう'),now());

-- docker-compose exec db bash
-- mysql -u root -p
-- show databases;
-- use test;
-- select * from users;