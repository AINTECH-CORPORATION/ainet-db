# AINET DB (`ainet-db`)

Motor de datos y runtime `.db` de **AINET** (Corporación **AINTECH**), basado en `pocketbase/pocketbase` (Licencia MIT).

## Clasificación Taxonómica
- **Nivel**: `AINET` (motores, modelos, runtime, datos y conocimiento)
- **Repositorio**: `AINTECH-CORPORATION/ainet-db`
- **Canon de Marca**: `branding-v2.20261103.1` (`guia_sha256: 54c489c7659f79e4205ee7f7fbc965e39d525e7d0ed4d827985c5e0e8956664a`)

## Capacidades SOTA Integradas (Roadmap & Registro)
| Módulo | Inspiración OSS | Función en `ainet-db` |
|---|---|---|
| `vector-vec0` | `asg017/sqlite-vec` | Búsqueda vectorial nativa dentro del mismo `.db` |
| `wal-stream` | `benbjohnson/litestream` | Replicación continua de WAL (`RPO ~ 1s`) |
| `db-explorer` | `simonw/datasette` | Explorador instantáneo solo-lectura + facetas + export CSV/JSON |
| `edge-branch` | `tursodatabase/libsql` | *Database branching* y sincronización edge |
| `bench` | Propia (`bench/bench_test.go`) | Benchmark reproducible WAL (`63.152 ns/op` write · `14.470 ns/op` read) |

## Verificación rápida
```bash
go test ./ainetdb/...
go test -bench=BenchmarkAINETDB -benchmem ./bench
```
