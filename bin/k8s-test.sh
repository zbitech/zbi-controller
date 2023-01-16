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
    --snapshot)
        snapshot="$2"
        shift
        ;;
    *)
        usage
        ;;
  esac
  shift
done

[ -z "${action}" ] && action=create
[ -z "${cluster}" ] && cluster=zbi-test
[ -z "${snapshot}" ] && snapshot=release-5.0

if [[ ${action} == "create" ]]; then
  echo "creating ${cluster} cluster"
  kind create cluster --name "${cluster}" --kubeconfig=kubeconfig
  export KUBECONFIG="${PWD}/kubeconfig"

  # project contour
  helm upgrade --install contour bitnami/contour -n contour --create-namespace

  # csi driver
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${snapshot}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotclasses.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${snapshot}/client/config/crd/snapshot.storage.k8s.io_volumesnapshotcontents.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${snapshot}/client/config/crd/snapshot.storage.k8s.io_volumesnapshots.yaml

  # create csi snapshot controller
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${snapshot}/deploy/kubernetes/snapshot-controller/rbac-snapshot-controller.yaml
  kubectl apply -f https://raw.githubusercontent.com/kubernetes-csi/external-snapshotter/${snapshot}/deploy/kubernetes/snapshot-controller/setup-snapshot-controller.yaml

  rm -rf csi-driver-host-path
  git clone https://github.com/kubernetes-csi/csi-driver-host-path.git
  curr_dir=${PWD}
  cd csi-driver-host-path/deploy/kubernetes-latest
  ./deploy.sh
  cd "$curr_dir"
  rm -rf csi-driver-host-path

  kubectl apply -f fake-zbi/data/manifest/zbi.yaml

elif [[ ${action} == "remove" ]]; then
  echo "deleting ${cluster} cluster"
  kind delete cluster --name "${cluster}"

else
  echo "Unknown action ${action}"
fi





