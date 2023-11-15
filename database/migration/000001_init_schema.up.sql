CREATE TABLE IF NOT EXISTS "users" (
    "name" VARCHAR(100) PRIMARY KEY,
    "password" VARCHAR(200) NOT NULL,
    "role" VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS "customers" (
    "name" VARCHAR(100) PRIMARY KEY REFERENCES "users"("name"),
    "money" INTEGER NOT NULL CHECK (salary BETWEEN 10000 AND 100000)
);

CREATE TABLE IF NOT EXISTS "workers" (
    "name" VARCHAR(100) PRIMARY KEY REFERENCES "users"("name"),
    "fatigue" INTEGER NOT NULL CHECK (fatigue BETWEEN 0 AND 100),
    "salary" INTEGER NOT NULL CHECK (salary BETWEEN 10000 AND 30000),
    "carryweight" INTEGER NOT NULL CHECK (carryweight BETWEEN 5 AND 30),
    "drunk" INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE IF NOT EXISTS "items" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "maxweight" INTEGER NOT NULL CHECK (maxweight >= 0),
    "minweight" INTEGER NOT NULL CHECK (minweight >= 0)
);

CREATE TABLE IF NOT EXISTS "tasks" (
    "id" SERIAL PRIMARY KEY,
    "itemname" VARCHAR(100) NOT NULL,
    "weight" INTEGER NOT NULL CHECK (weight >= 0)
);

CREATE TABLE IF NOT EXISTS "completetasks" (
    "id" SERIAL PRIMARY KEY,
    "workername" VARCHAR(100) NOT NULL REFERENCES "users"("name"),
    "itemname" VARCHAR(100) NOT NULL,
    "weight" INTEGER NOT NULL CHECK (weight >= 0)
);

INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('лампа', 1, 5);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('тарелки', 5, 15);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('компьютер', 10, 20);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('люстра', 15, 20);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('стул', 10, 30);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('холодильник', 30, 60);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('стол', 20, 55);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('кровать', 60, 100);