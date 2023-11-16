CREATE TABLE IF NOT EXISTS "users" (
    "name" VARCHAR(100) PRIMARY KEY,
    "password" VARCHAR(200) NOT NULL,
    "role" VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS "customers" (
    "name" VARCHAR(100) PRIMARY KEY REFERENCES "users"("name") ON DELETE CASCADE,
    "money" INTEGER NOT NULL CHECK (money BETWEEN 0 AND 100000),
    "lost" BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS "workers" (
    "name" VARCHAR(100) PRIMARY KEY REFERENCES "users"("name") ON DELETE CASCADE,
    "fatigue" INTEGER NOT NULL CHECK (fatigue BETWEEN 0 AND 100),
    "salary" INTEGER NOT NULL CHECK (salary BETWEEN 10000 AND 30000),
    "carryweight" INTEGER NOT NULL CHECK (carryweight BETWEEN 5 AND 30),
    "drunk" INTEGER NOT NULL DEFAULT 1,
    "choosen" BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS "items" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "maxweight" INTEGER NOT NULL CHECK (maxweight <= 80),
    "minweight" INTEGER NOT NULL CHECK (minweight >= 10)
);

CREATE TABLE IF NOT EXISTS "tasks" (
    "id" SERIAL PRIMARY KEY,
    "itemname" VARCHAR(100) NOT NULL,
    "weight" INTEGER NOT NULL CHECK (weight BETWEEN 10 AND 80)
);

CREATE TABLE IF NOT EXISTS "completetasks" (
    "id" SERIAL PRIMARY KEY,
    "workername" VARCHAR(100) NOT NULL REFERENCES "users"("name") ON DELETE CASCADE,
    "itemname" VARCHAR(100) NOT NULL,
    "weight" INTEGER NOT NULL CHECK (weight >= 0)
);

INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('лампа', 15, 10);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('тарелки', 17, 10);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('компьютер', 20, 10);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('люстра', 20, 15);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('стул', 30, 10);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('холодильник', 60, 30);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('стол', 55, 20);
INSERT INTO "items" ("name", "maxweight", "minweight") VALUES ('кровать', 80, 60);

-- docker run --rm -d --link cd81110a4fc9:postgres -p 8081:8080 --network movers_mynetwork adminer