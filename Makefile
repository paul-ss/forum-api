.PHONY: dbinit
dbinit:
	sudo -u postgres psql -f configs/sql/init.sql;


# Create all tables
.PHONY: dbsetup
dbsetup:
	PGPASSWORD='662f2710-4e08' psql -U forum_user -h 127.0.0.1 -d forum_db -f configs/sql/forum.sql -w

# Drop all created tables
.PHONY: dbclear
dbclear:
	echo  "select 'drop table if exists \"' || tablename || '\" cascade;' from pg_tables where schemaname = 'public';" > configs/sql/1.sql;
	PGPASSWORD='662f2710-4e08' psql -U forum_user -h 127.0.0.1 -d forum_db -f configs/sql/1.sql -w | grep drop > configs/sql/2.sql;
	PGPASSWORD='662f2710-4e08' psql -U forum_user -h 127.0.0.1 -d forum_db -f configs/sql/2.sql -w;
	rm configs/sql/1.sql configs/sql/2.sql;


# Build target api
.PHONY: build
build:
	go build -o build/bin/forum ./cmd/forum/main.go;

