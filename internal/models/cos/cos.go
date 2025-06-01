package cos

type Cos struct {
	SECRETID  string
	SecretKey string
	Name      string
}

type ID struct {
	ID string `json:"id"`
}

// FileName 要删除的文件名
type FileName struct {
	FileName string `json:"filename"`
}

type FileNames struct {
	FileName []string `json:"filename"`
}
