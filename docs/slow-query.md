# Setting up logging of Slow Query

**Add line to postgressql.conf**
* shared_preload_libraries = 'pg_stat_statements'

Then run "CREATE EXTENSION pg_stat_statements" in your database

