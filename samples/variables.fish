# Path to hosts file, can be relative or absolute
set -gx RINNEGAN_HOSTS_FILE ./samples/hosts

# Ensure appropriate hostname so that agents on targets can reach this to push data
# Eg: If using docker containers, add an entry in host's /etc/hosts that redirects
# this domain to localhost as containers pick up dns resolution on host generally.
set -gx RINNEGAN_INFLUX_HOST http://influxdb:8086

# Set to `true` when playing with docker containers
set -gx RINNEGAN_DOCKER false

set -gx RINNEGAN_DEBUG false

# Set to `true` when you are sshing as a user and use sudo for privesc
set -gx RINNEGAN_SUDO false
