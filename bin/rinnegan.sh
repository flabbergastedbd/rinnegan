#!/usr/bin/env bash

if [ -z "$RINNEGAN_INFLUX_HOST" ]; then
	echo "Variables not set I guess, RINNEGAN_INFLUX_HOST not found"
	exit 1
fi

if [ -z "$RINNEGAN_HOSTS_FILE" ]; then
	echo "Variables not set I guess, RINNEGAN_HOSTS_FILE not found"
	exit 1
fi

function debug {
	printf "\e[34m$@\e[0m\n"
}

function info {
	printf "\e[32m$@\e[0m\n"
}

HOST_REGEX="$1"
REMOTE_AGENT_DIR="/tmp/rinnegan"
REMOTE_AGENT="$REMOTE_AGENT_DIR/agent"
DOCKER_ARGS=""

if [ "x$RINNEGAN_DEBUG" == "xtrue" ]; then
	REMOTE_AGENT="$REMOTE_AGENT -v"
fi

if [ "x$RINNEGAN_DOCKER" != "xtrue" ]; then
	RINNEGAN_DOCKER=""
fi

case "$2" in
        agent)
		shift
		shift
		HOST_EXECUTE_COMMAND=''
		REMOTE_EXECUTE_COMMAND="$REMOTE_AGENT $@"
		;;
        deploy)
		if [ -n "$RINNEGAN_DOCKER" ]; then
			HOST_EXECUTE_COMMAND='docker cp build $h:$REMOTE_AGENT_DIR'
		else
			HOST_EXECUTE_COMMAND='scp -r build/. $h:$REMOTE_AGENT_DIR/'
		fi
		if [ -n "$RINNEGAN_DOCKER" ]; then
			DOCKER_ARGS="-d"
			REMOTE_EXECUTE_COMMAND="$REMOTE_AGENT daemon start --influxdb $RINNEGAN_INFLUX_HOST"
			# echo "Deploy docker agents by daemon execing: docker exec --privileged -t -u root -d \$h $REMOTE_EXECUTE_COMMAND"
			# REMOTE_EXECUTE_COMMAND=""
		else
			REMOTE_EXECUTE_COMMAND="nohup $REMOTE_AGENT daemon start --influxdb $RINNEGAN_INFLUX_HOST > $REMOTE_AGENT_DIR/agent.log 2>&1 &"
		fi
		;;
        stop)
		HOST_EXECUTE_COMMAND=""
		REMOTE_EXECUTE_COMMAND="$REMOTE_AGENT daemon stop"
		;;
        list)
		HOST_EXECUTE_COMMAND=""
		REMOTE_EXECUTE_COMMAND="$REMOTE_AGENT module list"
		;;
        wipe)
		HOST_EXECUTE_COMMAND=""
		REMOTE_EXECUTE_COMMAND="rm -rf $REMOTE_AGENT_DIR"
		;;
        exec)
		shift
		shift
		HOST_EXECUTE_COMMAND=""
		REMOTE_EXECUTE_COMMAND="$@"
		;;
        *)
		echo "Usage: rinnegan <host_regex> [agent|deploy|list|stop|wipe|exec]"
		echo ""
		echo "       <host_regex>  grep regex that will be applied to filter hosts"
		echo ""
		echo "          agent    Interact with agents deployed on targets"
		echo "          deploy   Deploy agents on to targets"
		echo "          list     List all active agents"
		echo "          stop     Stop all active agents"
		echo "          wipe     Remove any file leftovers on targets, run after stopping"
		echo "          exec     Run commands on targets directly, nothing fancy"
		echo ""
		exit 1

esac


debug "Filtering hosts file $RINNEGAN_HOSTS_FILE with regex $HOST_REGEX for targets"
TARGETS=$(grep -E "$HOST_REGEX" "$RINNEGAN_HOSTS_FILE")
info ""
info "Targets: "
info ""
info "$TARGETS"
info ""

for h in $TARGETS; do
	debug "======================================================================================================="
	debug "Running on $h"
	if [ -n "$HOST_EXECUTE_COMMAND" ]; then
		debug "-------------------------------------------------------------------------------------------------------"
		debug "Executing: $HOST_EXECUTE_COMMAND"
		debug "-------------------------------------------------------------------------------------------------------"
		$(eval echo -e "$HOST_EXECUTE_COMMAND")
		debug ""
	fi
	if [ -n "$REMOTE_EXECUTE_COMMAND" ]; then
		debug "-------------------------------------------------------------------------------------------------------"
		debug "Executing: $REMOTE_EXECUTE_COMMAND"
		debug "-------------------------------------------------------------------------------------------------------"
		if [ -n "$RINNEGAN_DOCKER" ]; then
			docker exec -it $DOCKER_ARGS -u root $h $REMOTE_EXECUTE_COMMAND
		else
			ssh $h "/bin/sh -c '$REMOTE_EXECUTE_COMMAND'"
		fi
		debug ""
	fi
	debug "======================================================================================================="
done
