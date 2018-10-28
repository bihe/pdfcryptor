# pdfcryptor
Golang encryption wrapper around qpdf/pdftk - change password of pdf.

## WHY
I could not remember the command line args, and I wanted to create a small golang application.

The result is a golang testbed and really simple ```cmd``` to use a PDF and change the password of the given PDF to a different one. So basically the whole logic of the cmd is to call the PDF utilities qpdf/pdftk with the command line args:


### QPDF
```
qpdf --password=YOURPASSWORD-HERE --decrypt input.pdf output.pdf
qpdf --encrypt <password>  <password> 128 -- doc_without_pass.pdf doc_with_pass.pdf
```

### PDFTK
```
pdftk document.pdf input_pw <password> output insecure.pdf
pdftk insecure.pdf output secure.pdf user_pw <password>
```
