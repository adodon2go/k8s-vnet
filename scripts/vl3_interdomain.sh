#!/usr/bin/env bash

usage() {
  echo "$(basename "$0")
Usage: $(basename "$0") [options...]
Options:
  --nse-hub=STRING          Hub for vL3 NSE images
                            (default=\"tiswanso\", environment variable: NSE_HUB)
  --nse-tag=STRING          Tag for vL3 NSE images
                            (default=\"kind_ci\", environment variable: NSE_TAG)
  --kconf=STRING            kube config file
  --ipam-octet=STRING       ipam octet
  --ipam-pool=STRING        ipam pool
  --remote-ip=STRING        remote ip
  --cnns-nsr-addr=STRING
  --cnns-nsr-port=STRING
" >&2
}

NSE_HUB=${NSE_HUB:-"tiswanso"}
NSE_TAG=${NSE_TAG:-"kind_ci"}
PULLPOLICY=${PULLPOLICY:-IfNotPresent}
INSTALL_OP=${INSTALL_OP:-apply}
for i in "$@"; do
    case $i in
        --nse-hub=*)
	          NSE_HUB=${i#*=}
	          ;;
        --nse-tag=*)
            NSE_TAG=${i#*=}
	          ;;
        --kconf=?*)
            KCONF=${i#*=}
            ;;
        --ipam-pool=?*)
            IPAMPOOL=${i#*=}
            ;;
        --ipam-octet=?*)
            IPAMOCTET=${i#*=}
            ;;
        --remote-ip=?*)
            REMOTE_IP=${i#*=}
            ;;
        --cnns-nsr-addr=?*)
            CNNS_NSR_ADDR=${i#*=}
            ;;
        --cnns-nsr-port=?*)
            CNNS_NSR_PORT=${i#*=}
            ;;
        --delete)
            INSTALL_OP=delete
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            usage
            exit 1
            ;;
    esac
done

sdir=$(dirname ${0})

if [[ -n ${CNNS_NSR_ADDR} ]]; then
    REMOTE_IP=${CNNS_NSR_ADDR}
fi

VL3HELMDIR=${VL3HELMDIR:-${sdir}/../k8s/helm}

CFGMAP="configmap nsm-vl3"
if [[ "${INSTALL_OP}" == "delete" ]]; then
    echo "delete configmap"
    kubectl delete ${KCONF:+--kubeconfig $KCONF} ${CFGMAP}
else
    if [[ -n ${REMOTE_IP} ]]; then
        kubectl create ${KCONF:+--kubeconfig $KCONF} ${CFGMAP} --from-literal=remote.ip_list=${REMOTE_IP}
    fi
fi

echo "---------------Install NSE-------------"
helm template ${VL3HELMDIR}/vl3 \
  --set "org=${NSE_HUB}" \
  --set "tag=${NSE_TAG}" \
  --set "pullPolicy=${PULLPOLICY}" \
  ${IPAMPOOL:+ --set "ipam.prefixPool=$IPAMPOOL"} \
  ${IPAMOCTET:+ --set "ipam.uniqueOctet=$IPAMOCTET"} \
  ${CNNS_NSR_ADDR:+ --set "cnns.addr=$CNNS_NSR_ADDR"} \
  ${CNNS_NSR_PORT:+ --set "cnns.port=$CNNS_NSR_PORT"} | kubectl ${INSTALL_OP} ${KCONF:+--kubeconfig $KCONF} -f -

if [[ "$INSTALL_OP" != "delete" ]]; then
  sleep 20
  kubectl wait ${KCONF:+--kubeconfig $KCONF} --timeout=150s --for condition=Ready -l networkservicemesh.io/app=vl3-nse pod
fi
