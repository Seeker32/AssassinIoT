-- Modify "device_telemetry" table
ALTER TABLE "device_telemetry"
  ADD COLUMN "value_int" integer NULL,
  ADD COLUMN "value_float" double precision NULL,
  DROP COLUMN "value_number";
