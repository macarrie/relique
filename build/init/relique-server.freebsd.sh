#!/bin/sh

# PROVIDE: relique_server
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

# Add the following lines to /etc/rc.conf to enable relique_server:
#
# relique_server_enable : set to "YES" to enable the daemon, default is "NO"
#
# relique_server_config : set to "/usr/local/etc/relique/server.toml" by default
#
# relique_server_options : set to empty ("") by default

. /etc/rc.subr

name=relique_server
rcvar=relique_server_enable

load_rc_config $name

: ${relique_server_enable="NO"}
: ${relique_server_config="/usr/local/etc/relique/server.toml"}
: ${relique_server_options=""}

log_file="/var/log/relique/relique_server.log"

command="/usr/sbin/daemon"
command_args="-u relique -t relique -o ${log_file} /usr/local/bin/relique-server start --config ${relique_server_config} ${relique_server_options}"

run_rc_command "$1"
