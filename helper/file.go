package helper

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	// Quelldatei öffnen
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer srcFile.Close()

	// Sicherstellen, dass Zielverzeichnis existiert
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("cannot create destination dir: %w", err)
	}

	// Zieldatei erstellen (ggf. überschreiben)
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer dstFile.Close()

	// Dateiinhalt kopieren
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("copy failed: %w", err)
	}

	return nil
}

func MoveFile(src, dst string) error {
	if err := CopyFile(src, dst); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		return err
	}
	return nil
}

func ReplaceExtension(fileName, newExt string) string {
	ext := filepath.Ext(fileName)
	return fileName[:len(fileName)-len(ext)] + newExt
}
