package processer

import (
	"fmt"

	"github.com/LinkinStars/words/internal/notebook"
)

func ImportOrExport(importFilePath, exportFilePath string) {
	if len(exportFilePath) > 0 {
		if err := notebook.ExportMistakeWords(exportFilePath); err != nil {
			panic(err)
		}
		fmt.Printf("成功导出到 %s", exportFilePath)
	} else if len(importFilePath) > 0 {
		if err := notebook.ImportMistakeWords(importFilePath); err != nil {
			panic(err)
		}
		fmt.Printf("成功导入到 %s", importFilePath)
	}
	return
}
