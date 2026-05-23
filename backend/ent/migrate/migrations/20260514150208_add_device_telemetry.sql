-- Create "device_telemetry" table
CREATE TABLE "device_telemetry" (
  "time" timestamptz NOT NULL,
  "device_identifier" character varying NOT NULL,
  "property_name" character varying NOT NULL,
  "value_type" character varying NOT NULL,
  "value_number" double precision NULL,
  "value_string" character varying NULL,
  "value_bool" boolean NULL,
  "value_json" jsonb NULL,
  "received_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT "device_telemetry_pkey" PRIMARY KEY ("time", "device_identifier", "property_name")
) WITH (
  tsdb.hypertable,
  tsdb.partition_column = 'time',
  tsdb.segmentby = 'device_identifier',
  tsdb.orderby = 'time DESC'
);

-- Create index "devicetelemetry_device_identifier_property_name_time" to table: "device_telemetry"
CREATE INDEX "devicetelemetry_device_identifier_property_name_time" ON "device_telemetry" ("device_identifier", "property_name", "time");
-- Create index "devicetelemetry_property_name_time" to table: "device_telemetry"
CREATE INDEX "devicetelemetry_property_name_time" ON "device_telemetry" ("property_name", "time");
-- Create index "devicetelemetry_received_at" to table: "device_telemetry"
CREATE INDEX "devicetelemetry_received_at" ON "device_telemetry" ("received_at");
