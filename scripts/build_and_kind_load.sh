#!/usr/bin/env bash

usage() {
  echo "$(basename "$0")
Usage: $(basename "$0") [options...]
Options:
  --nse-hub=STRING          Hub for vL3 NSE images
                            (default=\"tiswanso\", environment variable: NSE_HUB)
  --nse-tag=STRING          Tag for vL3 NSE images
                            (default=\"kind_ci\", environment variable: NSE_TAG)
  --kind-clusters=STRING    kind space separated cluster names
                            (default=\"kind-1\", environment variable: KIND_NAME)
" >&2
}

NSE_HUB=${NSE_HUB:-"tiswanso"}
NSE_TAG=${NSE_TAG:-"latest"}
KIND_CLUSTERS=${KIND_CLUSTERS:-"kind-1 kind-2"}

for i in "$@"; do
    case $i in
        --nse-hub=*)
	          NSE_HUB=${i#*=}
	          ;;
        --nse-tag=*)
            NSE_TAG=${i#*=}
	          ;;
        --kind-clusters=?*)
            KIND_CLUSTERS=${i#*=}
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

pushd "$(dirname "$0")/.." || (echo "the script was executed from PATH and dirname not present" && exit 1)

echo "Build go binary and docker image"

ORG=${NSE_HUB} TAG=${NSE_TAG} make all

popd || (echo "unexpected popd statement" && exit 1)

for cluster in ${KIND_CLUSTERS[@]}; do
  echo "Load docker image into cluster: $cluster"
  kind load docker-image "$NSE_HUB/vl3-nse:$NSE_TAG" --name="$cluster"
done