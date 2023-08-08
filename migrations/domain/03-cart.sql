CREATE TABLE `user` (
  `id` char(36) PRIMARY KEY,
  `email` varchar(255) UNIQUE NOT NULL,
  `username` varchar(255) UNIQUE NOT NULL,
  `name` varchar(255),
  `password` varchar(255) NOT NULL,
  `role` ENUM ('trainee', 'admin') NOT NULL,
  `cart_id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` char(36) NOT NULL,
  `updated_by` char(36) NOT NULL,
  `deleted_by` char(36) NULL DEFAULT NULL 
);

CREATE TABLE `cart` (
  `id` char(36) PRIMARY KEY,
  `user_id` char(36) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` char(36) NOT NULL,
  `updated_by` char(36) NOT NULL,
  `deleted_by` char(36) NULL DEFAULT NULL 
);

CREATE TABLE `product` (
  `id` char(36) PRIMARY KEY,
  `name` varchar(255) NOT NULL,
  `stock` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` char(36) NOT NULL,
  `updated_by` char(36) NOT NULL,
  `deleted_by` char(36) NULL DEFAULT NULL 
);

CREATE TABLE `cart_item` (
  `id` char(36) PRIMARY KEY,
  `cart_id` char(36) NOT NULL,
  `product_id` char(36) NOT NULL,
  `quantity` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` char(36) NOT NULL,
  `updated_by` char(36) NOT NULL,
  `deleted_by` char(36) NULL DEFAULT NULL 
);

CREATE TABLE `order` (
  `id` char(36) PRIMARY KEY,
  `user_id` char(36) NOT NULL,
  `total_price` decimal(10,2) NOT NULL,
  `address` varchar(255) NOT NULL,
  `status` ENUM ('pending', 'shipped', 'paid') NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` char(36) NOT NULL,
  `updated_by` char(36) NOT NULL,
  `deleted_by` char(36) NULL DEFAULT NULL 
);

CREATE TABLE `order_item` (
  `id` char(36) PRIMARY KEY,
  `order_id` char(36) NOT NULL,
  `product_id` char(36) NOT NULL,
  `quantity` int NOT NULL,
  `price` decimal(10,2) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_by` char(36) NOT NULL,
  `updated_by` char(36) NOT NULL,
  `deleted_by` char(36) NULL DEFAULT NULL 
);


ALTER TABLE `cart` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE;
ALTER TABLE `cart_item` ADD FOREIGN KEY (`cart_id`) REFERENCES `cart` (`id`) ON DELETE CASCADE;
ALTER TABLE `cart_item` ADD FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE;
ALTER TABLE `order` ADD FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE;
ALTER TABLE `order_item` ADD FOREIGN KEY (`order_id`) REFERENCES `order` (`id`) ON DELETE CASCADE;
ALTER TABLE `order_item` ADD FOREIGN KEY (`product_id`) REFERENCES `product` (`id`) ON DELETE CASCADE;



