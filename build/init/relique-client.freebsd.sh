#!/bin/sh

# PROVIDE: relique_client
# REQUIRE: DAEMON NETWORKING
# KEYWORD: shutdown

# Add the following lines to /etc/rc.conf to enable relique_client:
#
# relique_client_enable : set to "YES" to enable the daemon, default is "NO"
#
# relique_client_config : set to "/usr/local/etc/relique/client.toml" by default
#
# relique_client_options : set to empty ("") by default

. /etc/rc.subr

name=relique_client
rcvar=relique_client_enable
start_precmd="relique_prestart"

relique_prestart() {
	install -d -o relique -g relique -m 755 /var/log/relique
}

load_rc_config $name

: ${relique_client_enable="NO"}
: ${relique_client_config="/usr/local/etc/relique/client.toml"}
: ${relique_client_options=""}

log_file="/var/log/relique/relique_client.log"

command="/usr/sbin/daemon"
command_args="-u relique -t relique -o ${log_file} /usr/local/bin/relique-client start --config ${relique_client_config} ${relique_client_options}"

run_rc_command "$1"
