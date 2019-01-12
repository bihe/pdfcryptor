package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/bihe/pdfcryptor/internal/config"
	"github.com/bihe/pdfcryptor/internal/pdfcrypto"
)

func usage() {
	fmt.Fprintf(os.Stderr,
		"\n"+
			"pdfcryptor - Utility to change the encryption-password of a given PDF document.\n\n"+
			"usage: %s\n"+
			"\t--pdfIn=path of origin pdf-document\n"+
			"\t[--passIn=<password of the origin-pdf>]\n"+
			"\t--pdfOut=<path of the destination pdf-document>\n"+
			"\t--passOut=<password of the dest-pdf>\n"+
			"\nInfo: The application uses \"qpdf\" or \"pdftk\". If the binary is not available in the path, the application will stop!"+
			"\n\n",
		os.Args[0])
	os.Exit(1)
}

// The application only works if a compatible pdf-utility is available.
// Check if the qpdf/pdftk utility exists - by looking-up the path.
func checkPdfUtility() config.PdfUtil {
	_, err := exec.LookPath("qpdf")
	if err != nil {
		_, err = exec.LookPath("pdftk")
		if err != nil {
			exitError("Cannot find \"qpdf\" or \"pdftk\" utility in the path!\nProgram is terminated!\n")
		}
		return config.PDFTK
	}
	return config.QPDF
}

func exitError(errMsg string) {
	fmt.Fprintf(os.Stderr, errMsg+"\n")
	os.Exit(2)
}

func main() {
	pdf1 := flag.String("pdfIn", "", "The path of the origin-document.")
	pass1 := flag.String("passIn", "", "The password of the origin-pdf.")
	pdf2 := flag.String("pdfOut", "", "The path of the destination-document.")
	pass2 := flag.String("passOut", "", "The password of the destination-pdf.")
	flag.Parse()

	if len(*pdf1) > 0 && len(*pdf2) > 0 && len(*pass2) > 0 {
		utilType := checkPdfUtility()

		basePath, err := os.Getwd()
		if err != nil {
			exitError(fmt.Sprintf("Could not get the current path. Error: %s", err))
		}

		file, err := pdfcrypto.ChangePass(basePath, *pdf1, *pass1, *pdf2, *pass2, utilType)
		if err != nil {
			exitError(fmt.Sprintf("Could not decrypt/encrypt. Error: %s", err))
		}
		fmt.Printf("Created new encrypted file at %s\n", file)

		os.Exit(0)
	}
	usage()
}
