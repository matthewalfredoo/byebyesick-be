# ByeByeSick Backend

## How to Run?

### Check Things First
Check these files
1. `.env` 
2. `.dockerimages` (we are going to use this, at least for now)

### Docker
Make sure you have installed docker.

### Makefile
You knew it, install `make`.

### Database (Postgres)
We are running postgres in a docker container. The command is available on the [Makefile](Makefile), so make use of it.
Change it accordingly as per your requirements.

### Database Migration
Read [this documentation](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md) to learn how
to install `migrate` on your machine.

If you have installed or already have the `migrate` on your machine you can run
```bash
make migrateup
```
to start migrating to your database (make sure you already created one). Or

```bash
make migratedown
```
to remove all the migrations you have done.