DROP TABLE IF EXISTS `assignments`;
CREATE TABLE `assignments` (
  `id` varchar(200) NOT NULL,
  `course_id` varchar(200) DEFAULT NULL,
  `title` varchar(225) DEFAULT NULL,
  `description` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_assignments_deleted_at` (`deleted_at`),
  KEY `fk_assignments_course` (`course_id`),
  CONSTRAINT `fk_assignments_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `categories`;
CREATE TABLE `categories` (
  `id` varchar(200) NOT NULL,
  `name` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;


DROP TABLE IF EXISTS `courses`;
CREATE TABLE `courses` (
  `id` varchar(200) NOT NULL,
  `mentor_id` varchar(200) DEFAULT NULL,
  `category_id` varchar(200) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `description` longtext,
  `thumbnail` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_courses_deleted_at` (`deleted_at`),
  KEY `fk_courses_mentor` (`mentor_id`),
  KEY `fk_courses_category` (`category_id`),
  CONSTRAINT `fk_courses_category` FOREIGN KEY (`category_id`) REFERENCES `categories` (`id`),
  CONSTRAINT `fk_courses_mentor` FOREIGN KEY (`mentor_id`) REFERENCES `mentors` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `materials`;
CREATE TABLE `materials` (
  `id` varchar(200) NOT NULL,
  `module_id` varchar(200) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `url` longtext,
  `description` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_materials_deleted_at` (`deleted_at`),
  KEY `fk_materials_module` (`module_id`),
  CONSTRAINT `fk_materials_module` FOREIGN KEY (`module_id`) REFERENCES `modules` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mentee_assignments`;
CREATE TABLE `mentee_assignments` (
  `id` varchar(200) NOT NULL,
  `mentee_id` varchar(200) DEFAULT NULL,
  `assignment_id` varchar(200) DEFAULT NULL,
  `assignment_url` longtext,
  `grade` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_mentee_assignments_mentee` (`mentee_id`),
  KEY `fk_mentee_assignments_assignment` (`assignment_id`),
  CONSTRAINT `fk_mentee_assignments_assignment` FOREIGN KEY (`assignment_id`) REFERENCES `assignments` (`id`),
  CONSTRAINT `fk_mentee_assignments_mentee` FOREIGN KEY (`mentee_id`) REFERENCES `mentees` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mentee_courses`;
CREATE TABLE `mentee_courses` (
  `id` varchar(200) NOT NULL,
  `mentee_id` varchar(200) DEFAULT NULL,
  `course_id` varchar(200) DEFAULT NULL,
  `status` varchar(50) DEFAULT NULL,
  `reviewed` varchar(1) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_mentee_courses_mentee` (`mentee_id`),
  KEY `fk_mentee_courses_course` (`course_id`),
  CONSTRAINT `fk_mentee_courses_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`),
  CONSTRAINT `fk_mentee_courses_mentee` FOREIGN KEY (`mentee_id`) REFERENCES `mentees` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mentee_progresses`;
CREATE TABLE `mentee_progresses` (
  `id` varchar(200) NOT NULL,
  `mentee_id` varchar(200) DEFAULT NULL,
  `course_id` varchar(200) DEFAULT NULL,
  `material_id` varchar(200) DEFAULT NULL,
  `completed` varchar(1) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_mentee_progresses_course` (`course_id`),
  KEY `fk_mentee_progresses_material` (`material_id`),
  KEY `fk_mentee_progresses_mentee` (`mentee_id`),
  CONSTRAINT `fk_mentee_progresses_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`),
  CONSTRAINT `fk_mentee_progresses_material` FOREIGN KEY (`material_id`) REFERENCES `materials` (`id`),
  CONSTRAINT `fk_mentee_progresses_mentee` FOREIGN KEY (`mentee_id`) REFERENCES `mentees` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mentees`;
CREATE TABLE `mentees` (
  `id` varchar(200) NOT NULL,
  `user_id` varchar(200) DEFAULT NULL,
  `fullname` varchar(255) DEFAULT NULL,
  `phone` varchar(15) DEFAULT NULL,
  `role` varchar(50) DEFAULT NULL,
  `birth_date` longtext,
  `profile_picture` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_mentees_user` (`user_id`),
  CONSTRAINT `fk_mentees_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `mentors`;
CREATE TABLE `mentors` (
  `id` varchar(191) NOT NULL,
  `user_id` varchar(200) DEFAULT NULL,
  `fullname` varchar(255) DEFAULT NULL,
  `phone` varchar(15) DEFAULT NULL,
  `role` varchar(50) DEFAULT NULL,
  `jobs` longtext,
  `gender` longtext,
  `birth_place` longtext,
  `birth_date` datetime(3) DEFAULT NULL,
  `address` longtext,
  `profile_picture` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_mentors_user` (`user_id`),
  CONSTRAINT `fk_mentors_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `modules`;
CREATE TABLE `modules` (
  `id` varchar(200) NOT NULL,
  `course_id` varchar(200) DEFAULT NULL,
  `title` varchar(255) DEFAULT NULL,
  `description` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_modules_deleted_at` (`deleted_at`),
  KEY `fk_modules_course` (`course_id`),
  CONSTRAINT `fk_modules_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `reviews`;
CREATE TABLE `reviews` (
  `id` varchar(200) NOT NULL,
  `mentee_id` varchar(200) DEFAULT NULL,
  `course_id` varchar(200) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `rating` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_reviews_mentee` (`mentee_id`),
  KEY `fk_reviews_course` (`course_id`),
  CONSTRAINT `fk_reviews_course` FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`),
  CONSTRAINT `fk_reviews_mentee` FOREIGN KEY (`mentee_id`) REFERENCES `mentees` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` varchar(200) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;