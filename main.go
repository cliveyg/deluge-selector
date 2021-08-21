package main

import (
	"fmt"
	human "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/cliveyg/deluge-selector/entity"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)



func main() {

	allPars, err := getPartitionInfo()
	if err != nil {
		println(fmt.Errorf("we has a bad ting %s", err))
		os.Exit(1)
	}

	err = findDelugeCards(allPars)
	if err != nil {
		println(fmt.Errorf("we has a bad ting %s", err))
	}

	displayPartitionInfo(allPars)

}

func getPartitionInfo() ([]*PartitionInfo, error) {

	var pia []*PartitionInfo
	parts, _ := disk.Partitions(true)
	for _, p := range parts {
		device := p.Mountpoint
		s, err := disk.Usage(device)
		if err != nil {
			return nil, err
		}
		if s.Total == 0 {
			continue
		}

		percent := fmt.Sprintf("%2.f%%", s.UsedPercent)
		empty := false
		if percent == " 0%" {
			empty = true
		}

		volName := strings.Split(p.Mountpoint, string(os.PathSeparator))
		var vN string
		sysDisk := true
		delugeCard := false

		if runtime.GOOS == "darwin" && (volName[1] == "" || volName[1] == "dev" || volName[1] == "System") {
			vN = "MacOS System"
		} else {
			sysDisk = false
			vN = volName[len(volName)-1]
		}

		pi := NewPartInfo(s.Fstype, s.Total, s.Used, s.Free, percent, vN, p.Mountpoint, sysDisk, delugeCard, empty)
		pia = append(pia, pi)
	}

	return pia, nil
}

func displayPartitionInfo(allP []*PartitionInfo) {
	formatter := "%-14s %7s %7s %7s %4s %20s %8s %6s %5s %s\n"
	fmt.Printf(formatter, "Filesystem", "Size", "Used", "Avail", "Use%", "Volume Name", "SysDisk", "Deluge", "Empty", "Mounted on")

	for _, p := range allP {
		sysDisk := "Y   "
		if !p.SysDisk {
			sysDisk = "N   "
		}
		dCard := "N  "
		if p.DelugeCard {
			dCard = "Y  "
		}
		empty := "N "
		if p.Empty {
			empty = "Y "
		}
		fmt.Printf(formatter,
			p.FileSystem,
			human.Bytes(p.Size),
			human.Bytes(p.Used),
			human.Bytes(p.Available),
			p.Percent,
			p.VolumeName,
			sysDisk,
			dCard,
			empty,
			p.Mountpoint,
		)
	}

	//b, err := json.Marshal(allP)
	//if err != nil{
	//	fmt.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(string(b))

}

func findDelugeCards(pia []*PartitionInfo) error {

	for _, pi := range pia {
		if pi.SysDisk {
			continue
		}
		files, err := ioutil.ReadDir(pi.Mountpoint)
		if err != nil {
			return err
		}
		fCount := 0
		for _, item := range files {
			if item.IsDir() && (item.Name() == "KITS" || item.Name() == "SAMPLES" || item.Name() == "SONGS" || item.Name() == "SYNTHS") {
				fCount = fCount + 1
			}
		}
		if fCount == 4 {
			pi.DelugeCard = true
		}
	}
	return nil
}