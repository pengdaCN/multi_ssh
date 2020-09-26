package common

type FileAlloc struct {
	// 文件名与真实文件路径映射
	fileMap map[string]string
	existsFileRecord map[string]int
}
