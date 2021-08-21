package entity

type PartitionInfo struct {
	FileSystem string `json:"file_system"`
	Size       uint64 `json:"size"`
	Used       uint64 `json:"used"`
	Available  uint64 `json:"available"`
	Percent    string `json:"percent"`
	VolumeName string `json:"volume_name"`
	Mountpoint string `json:"mountpoint"`
	SysDisk    bool   `json:"sysdisk"`
	DelugeCard bool   `json:"deluge_card"`
	Empty      bool   `json:"empty"`
}

func NewPartInfo(fileSystem string, size, used, available uint64, percent, volumeName, mountpoint string, sysDisk, delugeCard, empty bool) *PartitionInfo {
	return &PartitionInfo{
		FileSystem: fileSystem,
		Size:       size,
		Used:       used,
		Available:  available,
		Percent:    percent,
		VolumeName: volumeName,
		Mountpoint: mountpoint,
		SysDisk:    sysDisk,
		DelugeCard: delugeCard,
		Empty:      empty,
	}
}