#!/bin/bash
set -x
source "${CSE_HELPERS_FILEPATH}"
source "${CSE_DISTRO_HELPERS_FILEPATH}"

echo "  dns-search ${CUSTOM_SEARCH_DOMAIN_NAME}" | tee -a /etc/network/interfaces.d/50-cloud-init.cfg
systemctl_restart 20 5 10 networking
wait_for_apt_locks
# this script is called from `cse_main.sh` where we return ERR_CUSTOM_SEARCH_DOMAINS_FAIL in case this script ends with non-zero exit code
retrycmd_if_failure 10 5 120 apt-get -y install realmd sssd sssd-tools samba-common samba samba-common python2.7 samba-libs packagekit || exit 1
wait_for_apt_locks
echo "${CUSTOM_SEARCH_REALM_PASSWORD}" | realm join -U ${CUSTOM_SEARCH_REALM_USER}@$(echo "${CUSTOM_SEARCH_DOMAIN_NAME}" | tr /a-z/ /A-Z/) $(echo "${CUSTOM_SEARCH_DOMAIN_NAME}" | tr /a-z/ /A-Z/)
