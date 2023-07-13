#!/bin/bash

set -e

usage() {
  cat <<EOF
Generate
usage: ${0} [OPTIONS]
The following flags are required.
       --action action to take - create or remove
       --cluster the type of cluster - name of cluster
EOF
  exit 1
}

while [[ $# -gt 0 ]]; do
  case ${1} in
    --action)
        action="$2"
        shift
        ;;
    --cluster)
        cluster="$2"
        shift
        ;;
    *)
        usage
        ;;
  esac
  shift
done

[ -z "${action}" ] && action=create
[ -z "${cluster}" ] && cluster=zbi

if [[ ${action} == "create" ]]; then

  echo "creating ${cluster} cluster"
  kind create cluster --name "${cluster}" --kubeconfig=kubeconfig
  export KUBECONFIG="${PWD}/kubeconfig"

elif [[ ${action} == "remove" ]]; then
  echo "deleting ${cluster} cluster"
  kind delete cluster --name "${cluster}"

else
  echo "Unknown action ${action}"
fi
