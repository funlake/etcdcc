CREATE TABLE `center_config` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `key` varchar(190) COLLATE utf8mb4_unicode_ci NOT NULL,
  `val` text COLLATE utf8mb4_unicode_ci,
  `mod` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `env` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `key_uniq` (`key`,`env`,`mod`) USING BTREE,
  KEY `key_inx` (`key`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;