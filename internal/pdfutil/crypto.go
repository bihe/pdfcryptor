package pdfutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/bihe/pdfcryptor/internal/config"
	"github.com/bihe/pdfcryptor/internal/utils"
	homedir "github.com/mitchellh/go-homedir"
)

// ChangeCrypt will change the password1 of pdf1 to password2 and rename it to pdf2
func ChangeCrypt(basePath, pdf1, pass1, pdf2, pass2 string, utilType config.PdfUtil) (file string, err error) {

	pdf1, err = getPath(basePath, pdf1)
	if err != nil {
		return "", err
	}
	pdf2, err = getPath(basePath, pdf2)
	if err != nil {
		return "", err
	}

	file, dir, err := decrypt(pdf1, pass1, utilType)
	if err != nil {
		cleanUpErr := os.RemoveAll(dir) // silently clean up
		if cleanUpErr != nil {
			fmt.Printf("Could not clean-up dir %s. Error: %s", dir, cleanUpErr)
		}
		return "", err
	}
	err = encrypt(dir, file, pdf2, pass2, utilType)
	if err == nil {
		// cleanup
		cleanUpErr := os.RemoveAll(dir)
		if cleanUpErr != nil {
			fmt.Printf("Could not clean-up dir %s. Error: %s", dir, cleanUpErr)
		}

		return pdf2, nil
	}
	return "", err
}

func getPath(base, p string) (string, error) {

	// 0) is this already an absolute path?
	if filepath.IsAbs(p) {
		return p, nil
	}

	// 1) check if this is a "home-path"
	exP, err := homedir.Expand(p)
	if err != nil {
		return "", err
	}

	if filepath.IsAbs(exP) {
		return exP, nil
	}

	// 2) expand "local paths"
	joinPath := path.Join(base, exP)
	joinPath, err = filepath.Abs(joinPath)

	if filepath.IsAbs(joinPath) {
		return joinPath, nil
	}

	return "", fmt.Errorf("Could not get the cleaned path of %s", p)
}

func decrypt(pdf1, pass1 string, utilType config.PdfUtil) (file, dir string, err error) {

	tmpDir, err := ioutil.TempDir("", "pdfcrypt")
	if err != nil {
		return "", "", err
	}

	// the original file is untouched, create a tempfile and copy the contents into it
	tmpFile, err := ioutil.TempFile(tmpDir, "*")
	if err != nil {
		return "", "", err
	}
	encFile := tmpFile.Name()
	tmpFile.Close()

	err = utils.CopyFile(pdf1, encFile)
	if err != nil {
		return "", "", err
	}
	// additional temp file to hold the decrypted contents
	tmpFile, err = ioutil.TempFile(tmpDir, "*")
	if err != nil {
		return "", "", err
	}
	decFile := tmpFile.Name()
	tmpFile.Close()

	switch utilType {
	case config.QPDF:
		// qpdf --password=YOURPASSWORD-HERE --decrypt input.pdf output.pdf
		err = utils.RunCmd(tmpDir, "qpdf", "--password="+pass1, "--decrypt", encFile, decFile)
	case config.PDFTK:
		// pdftk document.pdf input_pw <password> output insecure.pdf
		err = utils.RunCmd(tmpDir, "pdftk", encFile, "input_pw", pass1, "outupt", decFile)
	}

	if err != nil {
		return "", "", err
	}

	return decFile, tmpDir, nil
}

func encrypt(tmpDir, pdfDec, pdfDest, pass string, utilType config.PdfUtil) error {

	var err error

	switch utilType {
	case config.QPDF:
		// qpdf --encrypt <password>  <password> 128 -- doc_without_pass.pdf doc_with_pass.pdf
		err = utils.RunCmd(tmpDir, "qpdf", "--encrypt", pass, pass, "128", "--use-aes=y", "--", pdfDec, pdfDest)
	case config.PDFTK:
		// pdftk insecure.pdf output secure.pdf user_pw <password>
		err = utils.RunCmd(tmpDir, "pdftk", pdfDec, "output", pdfDest, "user_pw", pass)
	}

	return err
}
