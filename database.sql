CREATE TABLE "trials" (
  "id" SERIAL PRIMARY KEY,
  "start_time" int
);

CREATE TABLE "data" (
  "trialid" int,
  "time" int,
  "a1" int,
  "a2" int,
  "a3" int,
  "a4" int,
  "a5" int
);

ALTER TABLE "data" ADD FOREIGN KEY ("trialid") REFERENCES "trials" ("id");
