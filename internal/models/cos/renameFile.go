package cos

type RenameFile struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

type MoveFile struct {
	OldName []string `json:"oldName"`
	NewName string   `json:"newName"`
}

type DeleteFile struct {
	DeleteName string `json:"deleteName"`
}
