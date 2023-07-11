package controller

import (
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"

	"github.com/zHenriqueGN/PDFSanitizer/pkg/models"
)

// GetRootFolders gets all the folders in the root
// where the code is executing
func GetRootFolders() (folders []string, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}

	dirEntrys, err := os.ReadDir(currentDir)
	if err != nil {
		return
	}

	for _, folder := range dirEntrys {
		if folder.IsDir() {
			if err != nil {
				return
			}
			folders = append(folders, folder.Name())
		}
	}
	return
}

// MapPDFFiles search for files in a folder and map them. Files that has a name
// that is already in the map will be ignored
func MapPDFFiles(folderPath string, PDFs *map[string]models.PDFFile) (err error) {
	dirEntrys, err := os.ReadDir(folderPath)
	if err != nil {
		return
	}

	for _, entry := range dirEntrys {
		if !IsPDF(entry.Name()) {
			continue
		}

		var PDF models.PDFFile
		PDF.Name = entry.Name()
		PDF.Path = path.Join(folderPath, PDF.Name)
		if _, keyExists := (*PDFs)[PDF.Name]; !keyExists {
			(*PDFs)[PDF.Name] = PDF
		}
	}

	return
}

// CreateSanitizedPDFsFolder move all the PDFs of a given map and move them
// to a folder where all the pdfs are unique
func CreateSanitizedPDFsFolder(PDFsMap map[string]models.PDFFile, logger *log.Logger) (err error) {
	dstFolder := "sanitized_pdfs"
	err = os.MkdirAll(dstFolder, 0777)
	if err != nil {
		return
	}

	var buffer int
	for range PDFsMap {
		buffer++
	}

	tasks := make(chan models.PDFFile, buffer)
	results := make(chan string, buffer)

	for i := 0; i < 100; i++ {
		go worker(tasks, results)
	}

	for _, PDF := range PDFsMap {
		tasks <- PDF
	}
	close(tasks)

	for i := 0; i < buffer; i++ {
		result := <-results
		logger.Println(result)
	}

	return
}

func copyPDFToSanitizedFolder(PDF models.PDFFile) string {
	sourcePDF, err := os.Open(PDF.Path)
	if err != nil {
		return fmt.Sprintf("Error at file %s: %s", PDF.Path, err)
	}
	defer sourcePDF.Close()

	dstPDF, err := os.Create(path.Join("sanitized_pdfs", PDF.Name))
	if err != nil {
		return fmt.Sprintf("Error at file %s: %s", PDF.Path, err)
	}
	defer dstPDF.Close()

	_, err = io.Copy(dstPDF, sourcePDF)
	if err != nil {
		return fmt.Sprintf("Error at file %s: %s", PDF.Path, err)
	}
	return fmt.Sprintf("'%s' copied to sanitized folder", PDF.Path)
}

func worker(tasks <-chan models.PDFFile, results chan<- string) {
	for task := range tasks {
		results <- copyPDFToSanitizedFolder(task)
	}
}

// IsPDF returns true if a file is a PDF or false if not
func IsPDF(filePath string) bool {
	ext := filepath.Ext(filePath)
	mimeType := mime.TypeByExtension(ext)
	return mimeType == "application/pdf"
}

// CreateLogger generates a new logger to write logs in the standart output and
// in a log file at the same time
func CreateLogger() (logFile *os.File, logger *log.Logger, err error) {
	logFile, err = os.OpenFile("sanitizerLogs.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return
	}
	logger = log.New(io.MultiWriter(logFile, os.Stdout), "", log.Ldate|log.Ltime)
	return
}
