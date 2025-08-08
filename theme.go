package main

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type MyTheme struct {
	fyne.Theme
	fontPath string
}

func (t *MyTheme) Font(style fyne.TextStyle) fyne.Resource {
	fontFile, err := os.ReadFile(t.fontPath)
	if err != nil {
		fmt.Printf("[MyTheme::Font] 读取字体文件失败, %v\n", err)
		return theme.DefaultTextFont()
	}
	return fyne.NewStaticResource(filepath.Base(t.fontPath), fontFile)
}
