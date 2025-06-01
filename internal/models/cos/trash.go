package cos

type RecoverTrashFile struct {
	OldName []string `json:"TrashFiles"`
}

type AddTrashFile struct {
	OldName []string `json:"TrashFiles"`
}
