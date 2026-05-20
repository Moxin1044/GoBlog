-- 确保数据库存在
CREATE DATABASE IF NOT EXISTS goblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE goblog;

-- 初始化站点配置
INSERT INTO site_configs (`key`, `value`, created_at, updated_at) VALUES
('site_name', 'GoBlog', NOW(), NOW()),
('site_logo', '', NOW(), NOW()),
('copyright', 'Copyright © 2024 GoBlog. All rights reserved.', NOW(), NOW()),
('icp', '', NOW(), NOW()),
('register_enabled', 'true', NOW(), NOW()),
('site_url', 'http://localhost:8080', NOW(), NOW()),
('site_description', 'A personal blog powered by GoBlog', NOW(), NOW()),
('ai_enabled', 'true', NOW(), NOW()),
('ai_max_daily_per_user', '50', NOW(), NOW()),
('monitor_enabled', 'true', NOW(), NOW()),
('monitor_map_enabled', 'true', NOW(), NOW()),
('auto_backup_enabled', 'false', NOW(), NOW()),
('auto_backup_cron', '0 3 * * *', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at=NOW();

-- 初始化默认分类
INSERT INTO categories (name, name_en, sort, created_at) VALUES
('技术', 'Technology', 1, NOW()),
('生活', 'Life', 2, NOW()),
('随笔', 'Essay', 3, NOW());

-- 初始化默认标签
INSERT INTO tags (name, name_en, created_at) VALUES
('Go', 'Go', NOW()),
('Vue', 'Vue', NOW()),
('Docker', 'Docker', NOW()),
('前端', 'Frontend', NOW()),
('后端', 'Backend', NOW());
