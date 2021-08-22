package main

import (
	"bytes"
	"deluge-selector/entity"
	"fmt"
	human "github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/disk"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const SOURCE_SAMPLE = "SAMPLES/DRUMS/Kick/808 Kick.wav"
const DEST_SAMPLE = "SAMPLES/DRUMS/Kick/808 Kick.wav"

func main() {

	allPars, err := getPartitionInfo()
	if err != nil {
		log.Fatal(err)
	}

	err = findDelugeCards(allPars)
	if err != nil {
		log.Fatal(err)
	}

	displayPartitionInfo(allPars)

	for _, pt := range allPars {
		if pt.DelugeCard {
			err := walkDelugeCard(pt.Mountpoint)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func getPartitionInfo() ([]*entity.PartitionInfo, error) {

	var pia []*entity.PartitionInfo
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

		pi := entity.NewPartInfo(s.Fstype, s.Total, s.Used, s.Free, percent, vN, p.Mountpoint, sysDisk, delugeCard, empty)
		pia = append(pia, pi)
	}

	return pia, nil
}

func moveFile(source, dest string, repSlice []*entity.FilePathReplacer) error {
	err := os.Rename(source, dest)
	if err != nil {
		return err
	}
	repSlice = append(repSlice, entity.NewFilePathReplacer(source, dest))
	return nil
}

func displayPartitionInfo(allP []*entity.PartitionInfo) {
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

func traverseCB (path string, info os.FileInfo, err error) error {
	if err != nil {
		// ignore permissions related errors
		if !strings.Contains(err.Error(), "operation not permitted") {
			return err
		}
	}
	if !info.IsDir() {
		err = stringFindAndReplace(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func stringFindAndReplace(filepath string) error {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	matched, err := regexp.Match(SOURCE_SAMPLE, b)
	if err != nil {
		return err
	}
	if matched {
		println(fmt.Sprintf("File %s contains the string", filepath))
		//output := bytes.Replace(b, []byte(SAMPLE), []byte(SAMPLE+"XXCZ"), -1)
		output := bytes.Replace(b, []byte(SOURCE_SAMPLE), []byte(DEST_SAMPLE), -1)
		if err = ioutil.WriteFile(filepath, output, 0666); err != nil {
			return err
		}
		println("...Replaced")
	}
	return nil
}

//func findString(path string) error {
//	f, err := os.Open(path)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	scanner := bufio.NewScanner(f)
//
//	// https://golang.org/pkg/bufio/#Scanner.Scan
//	for scanner.Scan() {
//		if strings.Contains(scanner.Text(), SAMPLE) {
//			println(fmt.Sprintf("File %s contains the string", path))
//		}
//	}
//
//	return nil
//}


func walkDelugeCard(path string) error {
	//var wg sync.WaitGroup
	err := filepath.Walk(path, traverseCB)
	if err != nil {
		println("BAD JUJU")
		return err
	}

	return nil
}

func findDelugeCards(pia []*entity.PartitionInfo) error {

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