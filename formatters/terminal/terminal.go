package terminal

import (
	"dimitrios_liappis/instancepeek/model"
	"io"
	"sort"

	"github.com/olekukonko/tablewriter"
)

func sortedMetaKeys(meta map[string]string) []string {
	var keys []string
	for k, _ := range meta {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys
}

func Format(instanceInfos []model.InstanceInfo, dst io.Writer) error {
	table := tablewriter.NewWriter(dst)

	table.SetAutoFormatHeaders(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	instance0 := instanceInfos[0]
	header := []string{
		instance0.InstanceType.Label,
		instance0.CPU.VCPUCount.Label,
		instance0.Disk.Typ.Label,
		instance0.Disk.Count.Label,
		instance0.Disk.SizeGiB.Label,
		instance0.Memory.SizeMiB.Label,
		instance0.Network.Performance.Label,
		instance0.Hardware.Hypervisor.Label,
		instance0.Hardware.Baremetal.Label,
	}
	metaKeys := sortedMetaKeys(instance0.Meta)
	header = append(header, metaKeys...)

	headerColors := []tablewriter.Colors{
		{tablewriter.Bold, tablewriter.BgGreenColor},
		{tablewriter.BgRedColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		{tablewriter.BgYellowColor, tablewriter.Bold, tablewriter.FgBlackColor},
		{tablewriter.BgCyanColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		{tablewriter.BgCyanColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		{tablewriter.BgCyanColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		{tablewriter.BgHiMagentaColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		{tablewriter.BgHiMagentaColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		{tablewriter.BgHiMagentaColor, tablewriter.Bold, tablewriter.FgWhiteColor},
	}

	// use the same color scheme for all meta values
	metaColor := tablewriter.Colors{tablewriter.BgBlueColor, tablewriter.Bold, tablewriter.FgWhiteColor}

	for range metaKeys {
		headerColors = append(headerColors, metaColor)
	}

	table.SetHeader(header)
	table.SetHeaderColor(headerColors...)

	var tableData [][]string

	for i, instanceInfo := range instanceInfos {
		tableData = append(tableData, []string{
			instanceInfo.InstanceType.Value,
			instanceInfo.CPU.VCPUCount.Value,
			instanceInfo.Disk.Typ.Value,
			instanceInfo.Disk.Count.Value,
			instanceInfo.Disk.SizeGiB.Value,
			instanceInfo.Memory.SizeMiB.Value,
			instanceInfo.Network.Performance.Value,
			instanceInfo.Hardware.Hypervisor.Value,
			instanceInfo.Hardware.Baremetal.Value,
		})
		for _, v := range instanceInfo.Meta {
			tableData[i] = append(tableData[i], v)
		}
	}

	table.AppendBulk(tableData)
	table.Render()
	return nil
}
