/*
 * @Author: jffan
 * @Date: 2025-11-03 16:19:06
 * @LastEditTime: 2025-11-04 10:16:22
 * @LastEditors: jffan
 * @FilePath: \tcas-cli\tees\tees.go
 * @Description:
 */
package tees

import (
	"gitee.com/nanhu-lab/tcas-cli/collectors"
	"gitee.com/nanhu-lab/tcas-cli/tees/csv"
	nvidia "gitee.com/nanhu-lab/tcas-cli/tees/nvidiamock"
	"gitee.com/nanhu-lab/tcas-cli/tees/tdx"
	"gitee.com/nanhu-lab/tcas-cli/tees/virtcca"
)

func GetCollectors() map[string]collectors.EvidenceCollector {
	collectorMap := make(map[string]collectors.EvidenceCollector, 0)
	csvCollector := csv.NewCollector()
	collectorMap[csvCollector.Name()] = csvCollector

	vccaCollector := virtcca.NewCollector()
	collectorMap[vccaCollector.Name()] = vccaCollector

	nvidiaCollector := nvidia.NewCollector()
	collectorMap[nvidiaCollector.Name()] = nvidiaCollector

	tdxCollector := tdx.NewCollector()
	collectorMap[tdxCollector.Name()] = tdxCollector

	return collectorMap

}
