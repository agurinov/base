#!/bin/sh
set -e

ACTION="$0"

if [ -z "$BMP_DEBUG_MODE" ] ; then
	BMP_DEBUG_MODE="--debug"
fi

echo "$BMP_DEBUG_MODE $ACTION"

exec ./app \
	--debug run


# exec "$JAVA" $JAVA_OPTS \
# 	$ORIENTDB_OPTS_MEMORY \
#     $JAVA_OPTS_SCRIPT \
#     $ORIENTDB_SETTINGS \
#     $DEBUG_OPTS \
#     -Ddistributed=true \
#     -Djava.util.logging.config.file="$ORIENTDB_LOG_CONF" \
#     -Dorientdb.config.file="$CONFIG_FILE" \
#     -Dorientdb.www.path="$ORIENTDB_WWW_PATH" \
#     -Dorientdb.build.number="develop@re47e693f1470a7a642461be26983d4eca70777fd; 2018-06-06 11:18:56+0000" \
#     -cp "$ORIENTDB_HOME/lib/orientdb-server-3.0.2.jar:$ORIENTDB_HOME/lib/*:$ORIENTDB_HOME/plugins/*" \
#     $ARGS com.orientechnologies.orient.server.OServerMain
