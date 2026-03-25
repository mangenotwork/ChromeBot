package host

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Zip 压缩：支持 文件 / 目录（递归）
func Zip(src string, dstZip string) error {
	// 创建目标 zip 文件
	zipFile, err := os.Create(dstZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// 创建 zip 写入器
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 获取源信息
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// 处理压缩
	var basePath string
	if srcInfo.IsDir() {
		basePath = filepath.Dir(src)
	} else {
		basePath = filepath.Dir(src)
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建 zip 中的文件头
		relPath, err := filepath.Rel(basePath, path)
		if err != nil {
			return err
		}
		relPath = filepath.ToSlash(relPath)

		// 目录
		if info.IsDir() {
			_, err = zipWriter.Create(relPath + "/")
			return err
		}

		// 文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// 写入文件信息
		zipFileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		zipFileHeader.Name = relPath

		// 创建写入
		writer, err := zipWriter.CreateHeader(zipFileHeader)
		if err != nil {
			return err
		}

		// 复制内容
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// Unzip 解压：解压 zip 到目标目录
func Unzip(zipFile string, dstDir string) error {
	// 打开 zip
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	// 创建目标目录
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// 遍历解压
	for _, file := range reader.File {
		// 路径安全处理（防止跨目录攻击）
		filePath := filepath.Join(dstDir, file.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(dstDir)+string(os.PathSeparator)) {
			continue
		}

		// 目录
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, 0755)
			continue
		}

		// 创建文件所在目录
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		// 打开文件
		srcFile, err := file.Open()
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// 创建目标文件
		dstFile, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// 复制内容
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
	}

	return nil
}
