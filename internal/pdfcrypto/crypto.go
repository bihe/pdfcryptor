package pdfcrypto

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/bihe/pdfcryptor/internal/config"
	homedir "github.com/mitchellh/go-homedir"
)

// ChangePass will change the password1 of pdf1 to password2 and put the result into pdf2
func ChangePass(basePath, pdf1, pass1, pdf2, pass2 string, utilType config.PdfUtil) (string, error) {

	var err error

	if pdf1, err = getPath(basePath, pdf1); err != nil {
		return "", err
	}
	if pdf2, err = getPath(basePath, pdf2); err != nil {
		return "", err
	}

	var file, dir string

	if file, dir, err = decrypt(pdf1, pass1, utilType); err != nil {
		cleanUpErr := os.RemoveAll(dir) // silently clean up
		if cleanUpErr != nil {
			fmt.Printf("Could not clean-up dir %s. Error: %s", dir, cleanUpErr)
		}
		return "", err
	}
	if err := encrypt(dir, file, pdf2, pass2, utilType); err == nil {
		// cleanup
		if cleanUpErr := os.RemoveAll(dir); cleanUpErr != nil {
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

	var exP string
	var err error

	// 1) check if this is a "home-path"
	if exP, err = homedir.Expand(p); err != nil {
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

	var tmpDir string
	var tmpFile *os.File

	if tmpDir, err = ioutil.TempDir("", "pdfcrypt"); err != nil {
		return "", "", err
	}

	// the original file is untouched, create a tempfile and copy the contents into it
	if tmpFile, err = ioutil.TempFile(tmpDir, "*"); err != nil {
		return "", "", err
	}
	encFile := tmpFile.Name()
	tmpFile.Close()

	if err = copyFile(pdf1, encFile); err != nil {
		return "", "", err
	}
	// additional temp file to hold the decrypted contents
	if tmpFile, err = ioutil.TempFile(tmpDir, "*"); err != nil {
		return "", "", err
	}
	decFile := tmpFile.Name()
	tmpFile.Close()

	switch utilType {
	case config.QPDF:
		// qpdf --password=YOURPASSWORD-HERE --decrypt input.pdf output.pdf
		err = runCmd(tmpDir, "qpdf", "--password="+pass1, "--decrypt", encFile, decFile)
	case config.PDFTK:
		// pdftk document.pdf input_pw <password> output insecure.pdf
		err = runCmd(tmpDir, "pdftk", encFile, "input_pw", pass1, "outupt", decFile)
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
		err = runCmd(tmpDir, "qpdf", "--encrypt", pass, pass, "128", "--use-aes=y", "--", pdfDec, pdfDest)
	case config.PDFTK:
		// pdftk insecure.pdf output secure.pdf user_pw <password>
		err = runCmd(tmpDir, "pdftk", pdfDec, "output", pdfDest, "user_pw", pass)
	}

	return err
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func runCmd(dir, name string, args ...string) error {
	// Create the command.
	var stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stderr = &stderr
	cmd.Dir = dir

	// Start the command and wait for it to exit.
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(strings.TrimSpace(stderr.String()))
	}

	return nil
}
