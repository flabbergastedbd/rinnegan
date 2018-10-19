Rinnegan
########

   Blessed are those who write good READMEs
   - Baba Blanka (fictitious)

Why the hell this exists?
*************************

To help visualize in reverse engineering complex distributed systems. Helps in
greatly decoding communications and events happening.

Requirements
************

Hosts
=====

* bash
* golang

Targets
=======

Non sudo privileged access is required with ptrace_scope set to 0.

Walkthrough
***********

Modules
*******

1. Network Routing
==================

Prerequisites
^^^^^^^^^^^^^

* Iptables need to be installed and kernel module need to be enabled on target.

.. code-block::

   # Disabled by default!
   $> echo "1" > /proc/sys/net/ipv4/ip_forward

   # Load iptables module, if not already
   $> modprobe ip_tables


Known Issues
************

This is quite unstable, so expect a bumpy ride

* ps on alpine is from busybox, so to get the real deal install procps
  on all alpine containers.

.. code-block::

   ./bin/rinnegan.sh "." exec apk add --no-cache procps

* Agent just smoked up too many processes, what to do? Just restart your containers/vms.
