# omo-msa-favorite
Micro Service Agent - favorite

MICRO_REGISTRY=consul micro call omo.msa.favorite FavoriteService.AddOne '{"name":"John", "owner":"11111", "remark":"test1", "type":2, "cover":"hhhhhh"}'
MICRO_REGISTRY=consul micro call omo.msa.favorite FavoriteService.RemoveOne '{"uid":"5f0fffef0d57c9d90026b782", "owner":"11111"}'
