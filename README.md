# DOAn (Digital Ocean Ansible Agent)

DOAn is a Golang based agent that runs on Digital Ocean droplets and syncs Ansible playbooks to the droplet from artifactory.

## Building debian package

Debian building has been updated to use fpm and was moved into the Makefile.
Install `fpm` following the steps here and run `make build-deb` to build the debian package. 
