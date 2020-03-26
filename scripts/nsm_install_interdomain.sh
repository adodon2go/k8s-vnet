#!/usr/bin/env bash

print_usage() {
  echo "$(basename "$0")
Usage: $(basename "$0") [options...]
Options:
  --nsm-hub=STRING          Hub for NSM images
                            (default=\"tiswanso\", environment variable: NSM_HUB) 
  --nsm-tag=STRING          Tag for NSM images
                            (default=\"vl3_api_rebase\", environment variable: NSM_TAG)
  --nsm-dir=STRING          nsm directory
  --kconf=STRING            kube config file
" >&2
}
PULLPOLICY=${PULLPOLICY:-IfNotPresent}
NSM_HUB="${NSM_HUB:-"tiswanso"}"
NSM_TAG="${NSM_TAG:-"vl3_api_rebase"}"
INSTALL_OP=${INSTALL_OP:-apply}

for i in "$@"
do
case $i in
    --nsm-hub=*)
    NSM_HUB="${i#*=}"
    ;;
    --nsm-tag=*)
    NSM_TAG="${i#*=}"
    ;;
    --nsm-dir=*)
    NSMDIR=$(realpath "${i#*=}")
    ;;
    --kconf=*)
    KCONF="${i#*=}"
    ;;
    -h|--help)
      print_usage
      exit 0
    ;;
    *)
      print_usage
      exit 1
    ;;
esac
done

sdir=$(dirname ${0})

NSMDIR=${NSMDIR:-${sdir}/../../../../networkservicemesh}
VL3DIR=${VL3DIR:-${sdir}/..}

echo "------------- Create nsm-system namespace ----------"
if [[ "${INSTALL_OP}" != "delete" ]]; then
  kubectl create ns nsm-system ${KCONF:+--kubeconfig $KCONF}
fi
echo "------------Installing NSM monitoring-----------"

helm template ${NSMDIR}/deployments/helm/crossconnect-monitor --namespace nsm-system --set insecure="true" --set global.JaegerTracing="true" | kubectl ${INSTALL_OP} ${KCONF:+--kubeconfig $KCONF} -f -
helm template ${NSMDIR}/deployments/helm/jaeger --namespace nsm-system --set insecure="true" --set global.JaegerTracing="true" --set monSvcType=NodePort | kubectl ${INSTALL_OP} ${KCONF:+--kubeconfig $KCONF} -f -
helm template ${NSMDIR}/deployments/helm/skydive --namespace nsm-system --set insecure="true" --set global.JaegerTracing="true" | kubectl ${INSTALL_OP} ${KCONF:+--kubeconfig $KCONF} -f -

echo "------------Installing NSM-----------"
helm template ${NSMDIR}/deployments/helm/nsm --namespace nsm-system --set org=${NSM_HUB},tag=${NSM_TAG} --set "pullPolicy=$PULLPOLICY" --set insecure="true" --set global.JaegerTracing="true" | kubectl ${INSTALL_OP} ${KCONF:+--kubeconfig $KCONF} -f -

echo "------------Installing NSM-addons -----------"
helm template ${VL3DIR}/k8s/helm/nsm-addons --namespace nsm-system --set global.NSRegistrySvc=true  | kubectl ${INSTALL_OP} ${KCONF:+--kubeconfig $KCONF} -f -

echo "------------Installing proxy NSM-----------"
helm template ${NSMDIR}/deployments/helm/proxy-nsmgr --namespace nsm-system --set org=${NSM_HUB},tag=${NSM_TAG} --set "pullPolicy=$PULLPOLICY" --set insecure="true" --set global.JaegerTracing="true" | kubectl ${INSTALL_OP} ${KCONF:+--kubeconfig $KCONF} -f -
