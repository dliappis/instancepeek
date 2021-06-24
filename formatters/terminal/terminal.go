package terminal

import (
	"dimitrios_liappis/instancepeek/model"
	"io"

	"github.com/olekukonko/tablewriter"
)

func Format(instanceInfos []model.InstanceInfo, dst io.Writer) error {
	table := tablewriter.NewWriter(dst)
	header := []string{"Instance Type", "vCPUs", "Memory MiB", "Local Disk Type", "Local Disk Count", "Local Disk Size", "Network", "Baremetal", "Hypervisor"}
	table.SetHeader(header)

	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
		tablewriter.Colors{tablewriter.BgRedColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgYellowColor, tablewriter.Bold, tablewriter.FgBlackColor},
		tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgHiMagentaColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgHiMagentaColor, tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.BgHiMagentaColor, tablewriter.Bold, tablewriter.FgWhiteColor},
	)

	// table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
	// 	tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
	// 	tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor})

	for _, instanceInfo := range instanceInfos {
		var row []string
		row = append(
			row,
			instanceInfo.InstanceType,
			instanceInfo.CPU.VCPUCount,
			instanceInfo.Memory.SizeMiB,
			instanceInfo.Disk.Typ,
			instanceInfo.Disk.Count,
			instanceInfo.Disk.SizeGiB,
			instanceInfo.Network.Performance,
			instanceInfo.Hardware.Baremetal,
			instanceInfo.Hardware.Hypervisor,
		)
		table.Append(row)
	}

	table.Render()
	return nil
}
