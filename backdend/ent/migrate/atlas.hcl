env "local" {
  url = "postgres://postgres:password@localhost:15432/postgres?search_path=public"

  dev = "docker+postgres://timescale/timescaledb:2.27.0-pg17/atlas_dev?search_path=public"


  migration {
    dir = "file://migrations"
  }

  exclude = [
    "_timescaledb_*",
    "timescaledb_information",
    "*[type=extension].version",
  ]
}
