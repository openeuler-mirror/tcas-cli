package tees

import (
	"tcas-cli/collectors"
	"tcas-cli/tees/csv"
	nvidia "tcas-cli/tees/nvidiamock"
	"tcas-cli/tees/virtcca"
)

func GetCollectors() map[string]collectors.EvidenceCollector {
	collectorMap := make(map[string]collectors.EvidenceCollector, 0)
	csvCollector := csv.NewCollector()
	collectorMap[csvCollector.Name()] = csvCollector

	vccaCollector := virtcca.NewCollector()
	collectorMap[vccaCollector.Name()] = vccaCollector

	nvidiaCollector := nvidia.NewCollector()
	collectorMap[nvidiaCollector.Name()] = nvidiaCollector

	return collectorMap

}
