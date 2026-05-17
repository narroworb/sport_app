# Service-level compose workflow

This repository uses one `docker-compose.yaml` per service, located under `services/<service>/docker-compose.yaml`.

The root `docker-compose.yml` has been removed — instead use the helper scripts to start services sequentially.

Provided scripts:

- `start-services.sh` — Bash script for Linux/macOS. Iterates `services/*/docker-compose.yaml` and runs `docker compose -f <file> up -d --build` for each file in turn. Ensures `sport_network` exists and verifies that common external volumes exist — it will NOT create missing volumes automatically.

- `start-services.bat` — Windows CMD equivalent. Verifies that required external volumes exist and will NOT create them automatically.

- `compose-all.sh` / `compose-all.bat` — legacy wrappers kept for compatibility; they call `docker compose` with multiple `-f` flags in a single command.

Usage (Linux/macOS):

```bash
chmod +x start-services.sh
./start-services.sh        # starts all services sequentially (up -d --build)
```

Usage (Windows CMD):

```cmd
start-services.bat         # starts all services sequentially
```

Notes:

- The scripts use `docker compose` (the newer plugin). If your environment uses the legacy `docker-compose` binary, edit the scripts and replace `docker compose` with `docker-compose`.
- The scripts create a Docker network named `sport_network` if it does not exist. They will verify the presence of common external volumes and will refuse to run if any required volumes are missing — this prevents accidental creation of empty volumes that could lead to data loss. Create missing volumes manually or restore backups before running.
- To stop or view logs for a specific service, run `docker compose -f services/<service>/docker-compose.yaml down` or `docker compose -f services/<service>/docker-compose.yaml logs -f`.

If you'd like, I can also add `stop-services.sh`/`.bat` that iterate the service compose files and bring them down in reverse order.


