package entity

type FilePathReplacer struct {
	Source string
	Dest   string
}

func NewFilePathReplacer(source, dest string) *FilePathReplacer {
	return &FilePathReplacer{
		Source: source,
		Dest: dest,
	}
}