( set -a; source .env.local.production; set +a; psql $POSTGRES_URL )

