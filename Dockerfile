FROM alpine:3.11
ADD omo.msa.favorite /usr/bin/omo.msa.favorite
ENV MSA_REGISTRY_PLUGIN
ENV MSA_REGISTRY_ADDRESS
ENTRYPOINT [ "omo.msa.favorite" ]
