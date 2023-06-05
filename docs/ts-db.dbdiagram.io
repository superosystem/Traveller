
Table "assignments" {
  "id" varchar(200) [pk, not null]
  "course_id" varchar(200) [default: NULL]
  "title" varchar(225) [default: NULL]
  "description" longtext
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]
  "deleted_at" datetime(3) [default: NULL]

Indexes {
  deleted_at [name: "idx_assignments_deleted_at"]
  course_id [name: "fk_assignments_course"]
}
}

Table "categories" {
  "id" varchar(200) [pk, not null]
  "name" varchar(255) [default: NULL]
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]
}

Table "courses" {
  "id" varchar(200) [pk, not null]
  "mentor_id" varchar(200) [default: NULL]
  "category_id" varchar(200) [default: NULL]
  "title" varchar(255) [default: NULL]
  "description" longtext
  "thumbnail" varchar(255) [default: NULL]
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]
  "deleted_at" datetime(3) [default: NULL]

Indexes {
  deleted_at [name: "idx_courses_deleted_at"]
  mentor_id [name: "fk_courses_mentor"]
  category_id [name: "fk_courses_category"]
}
}

Table "materials" {
  "id" varchar(200) [pk, not null]
  "module_id" varchar(200) [default: NULL]
  "title" varchar(255) [default: NULL]
  "url" longtext
  "description" longtext
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]
  "deleted_at" datetime(3) [default: NULL]

Indexes {
  deleted_at [name: "idx_materials_deleted_at"]
  module_id [name: "fk_materials_module"]
}
}

Table "mentee_assignments" {
  "id" varchar(200) [pk, not null]
  "mentee_id" varchar(200) [default: NULL]
  "assignment_id" varchar(200) [default: NULL]
  "assignment_url" longtext
  "grade" bigint [default: NULL]
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]

Indexes {
  mentee_id [name: "fk_mentee_assignments_mentee"]
  assignment_id [name: "fk_mentee_assignments_assignment"]
}
}

Table "mentee_courses" {
  "id" varchar(200) [pk, not null]
  "mentee_id" varchar(200) [default: NULL]
  "course_id" varchar(200) [default: NULL]
  "status" varchar(50) [default: NULL]
  "reviewed" varchar(1) [default: NULL]
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]

Indexes {
  mentee_id [name: "fk_mentee_courses_mentee"]
  course_id [name: "fk_mentee_courses_course"]
}
}

Table "mentee_progresses" {
  "id" varchar(200) [pk, not null]
  "mentee_id" varchar(200) [default: NULL]
  "course_id" varchar(200) [default: NULL]
  "material_id" varchar(200) [default: NULL]
  "completed" varchar(1) [default: NULL]
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]

Indexes {
  course_id [name: "fk_mentee_progresses_course"]
  material_id [name: "fk_mentee_progresses_material"]
  mentee_id [name: "fk_mentee_progresses_mentee"]
}
}

Table "mentees" {
  "id" varchar(200) [pk, not null]
  "user_id" varchar(200) [default: NULL]
  "fullname" varchar(255) [default: NULL]
  "phone" varchar(15) [default: NULL]
  "role" varchar(50) [default: NULL]
  "birth_date" longtext
  "profile_picture" longtext
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]

Indexes {
  user_id [name: "fk_mentees_user"]
}
}

Table "mentors" {
  "id" varchar(191) [pk, not null]
  "user_id" varchar(200) [default: NULL]
  "fullname" varchar(255) [default: NULL]
  "phone" varchar(15) [default: NULL]
  "role" varchar(50) [default: NULL]
  "jobs" longtext
  "gender" longtext
  "birth_place" longtext
  "birth_date" datetime(3) [default: NULL]
  "address" longtext
  "profile_picture" longtext
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]

Indexes {
  user_id [name: "fk_mentors_user"]
}
}

Table "modules" {
  "id" varchar(200) [pk, not null]
  "course_id" varchar(200) [default: NULL]
  "title" varchar(255) [default: NULL]
  "description" longtext
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]
  "deleted_at" datetime(3) [default: NULL]

Indexes {
  deleted_at [name: "idx_modules_deleted_at"]
  course_id [name: "fk_modules_course"]
}
}

Table "reviews" {
  "id" varchar(200) [pk, not null]
  "mentee_id" varchar(200) [default: NULL]
  "course_id" varchar(200) [default: NULL]
  "description" varchar(255) [default: NULL]
  "rating" bigint [default: NULL]
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]

Indexes {
  mentee_id [name: "fk_reviews_mentee"]
  course_id [name: "fk_reviews_course"]
}
}

Table "users" {
  "id" varchar(200) [pk, not null]
  "email" varchar(255) [default: NULL]
  "password" varchar(255) [default: NULL]
  "created_at" datetime(3) [default: NULL]
  "updated_at" datetime(3) [default: NULL]
}

Ref "fk_assignments_course":"courses"."id" < "assignments"."course_id"

Ref "fk_courses_category":"categories"."id" < "courses"."category_id"

Ref "fk_courses_mentor":"mentors"."id" < "courses"."mentor_id"

Ref "fk_materials_module":"modules"."id" < "materials"."module_id"

Ref "fk_mentee_assignments_assignment":"assignments"."id" < "mentee_assignments"."assignment_id"

Ref "fk_mentee_assignments_mentee":"mentees"."id" < "mentee_assignments"."mentee_id"

Ref "fk_mentee_courses_course":"courses"."id" < "mentee_courses"."course_id"

Ref "fk_mentee_courses_mentee":"mentees"."id" < "mentee_courses"."mentee_id"

Ref "fk_mentee_progresses_course":"courses"."id" < "mentee_progresses"."course_id"

Ref "fk_mentee_progresses_material":"materials"."id" < "mentee_progresses"."material_id"

Ref "fk_mentee_progresses_mentee":"mentees"."id" < "mentee_progresses"."mentee_id"

Ref "fk_mentees_user":"users"."id" < "mentees"."user_id"

Ref "fk_mentors_user":"users"."id" < "mentors"."user_id"

Ref "fk_modules_course":"courses"."id" < "modules"."course_id"

Ref "fk_reviews_course":"courses"."id" < "reviews"."course_id"

Ref "fk_reviews_mentee":"mentees"."id" < "reviews"."mentee_id"
