CREATE TYPE if not exists "user_role" AS ENUM (
  'user',
  'admin',
  'superadmin'
);

CREATE TYPE if not exists "gender" AS ENUM (
  'male',
  'female'
);

CREATE TYPE if not exists "room_category" AS ENUM (
  'single',
  'double',
  '3xroom',
  '4xroom',
  '5xroom'
);

CREATE TYPE if not exists "room_status" AS ENUM (
  'available',
  'booked',
  'maintenance'
);

CREATE TYPE if not exists "room_type" AS ENUM (
  'economy',
  'standard',
  'comfort',
  'deluxe',
  'suite'
);

CREATE TYPE if not exists "user_status" AS ENUM (
  'active',
  'blocked',
  'inverify'
);

CREATE TABLE if not exists "users" (
  "id" UUID PRIMARY KEY not NULL,
  "fullname" VARCHAR(100) not NULL,
  "username" VARCHAR(100) not NULL UNIQUE,
  "email" VARCHAR(100) not NULL UNIQUE,
  "password" VARCHAR(255) not NULL,
  "phone" VARCHAR(20) not NULL,
  "user_status" user_status NOT NULL DEFAULT 'inverify',
  "gender" gender NOT NULL DEFAULT 'male',
  "role" user_role NOT NULL DEFAULT 'user',
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE if not exists "rooms" (
  "id" UUID PRIMARY KEY,
  "type" room_type NOT NULL,
  "category" room_category NOT NULL,
  "status" room_status NOT NULL DEFAULT 'available',
  "price" DECIMAL(10,2),
  "availability" BOOLEAN,
  "rating" FLOAT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE if not exists "room_images" (
  "id" UUID PRIMARY KEY,
  "room_id" UUID NOT NULL,
  "image_url" TEXT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE if not exists "room_reviews" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "room_id" UUID NOT NULL,
  "rating" FLOAT,
  "comment" TEXT,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE if not exists "bookings" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "room_id" UUID NOT NULL,
  "check_in_date" DATE,
  "check_out_date" DATE,
  "status" VARCHAR(20),
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE if not exists "payments" (
  "id" UUID PRIMARY KEY,
  "booking_id" UUID NOT NULL,
  "amount" DECIMAL(10,2),
  "status" VARCHAR(20),
  "payment_date" TIMESTAMP,
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

CREATE TABLE if not exists "complaints" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "booking_id" UUID NOT NULL,
  "message" TEXT NOT NULL,
  "status" VARCHAR(20) DEFAULT 'pending',
  "created_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP),
  "updated_at" TIMESTAMP DEFAULT (CURRENT_TIMESTAMP)
);

ALTER TABLE "room_images" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "room_reviews" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "room_reviews" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "bookings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bookings" ADD FOREIGN KEY ("room_id") REFERENCES "rooms" ("id");

ALTER TABLE "payments" ADD FOREIGN KEY ("booking_id") REFERENCES "bookings" ("id");

ALTER TABLE "complaints" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "complaints" ADD FOREIGN KEY ("booking_id") REFERENCES "bookings" ("id");
