package process

import (
	"bufio"
	"os"
	"strings"
)

func MapearCorreos(file *os.File) *Email {
	// Crea un nuevo escaner para leer línea a línea del archivo
	scanner := bufio.NewScanner(file)
	email := &Email{}
	//Banderas de region en el correo segun el campo
	inMessageID := false
	inDate := false
	inFrom := false
	inTo := false
	inSubject := false
	inMime := false
	inContentType := false
	inContentTransfer := false
	inXFrom := false
	inXTo := false
	inXCc := false
	inXBcc := false
	inXFolder := false
	inXOrigin := false
	inXFileName := false
	inBody := false

	//Recorre el archivo
	for scanner.Scan() {

		line := scanner.Text()

		//Evaluar cantidad de lineas para Message-ID
		if strings.HasPrefix(line, "Message-ID:") && !inBody {
			inMessageID = true
			email.MessageID = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inMessageID && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "Date:") {
			email.MessageID += " " + line
			continue
		}

		//Evaluar cantidad de lineas para Date
		if strings.HasPrefix(line, "Date:") && !inBody {
			inMessageID = false
			inDate = true
			email.Date = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inDate && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "From:") {
			email.Date += " " + line
			continue
		}

		//Evaluar cantidad de lineas para From
		if strings.HasPrefix(line, "From:") && !inBody {
			inDate = false
			inFrom = true
			email.From = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inFrom && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "To:") {
			email.From += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para To
		if strings.HasPrefix(line, "To:") && !inBody {
			inFrom = false
			inTo = true
			email.To = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inTo && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "Subject:") {
			email.To += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para Subject
		if strings.HasPrefix(line, "Subject:") && !inBody {
			inTo = false
			inSubject = true
			email.Subject = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inSubject && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "Mime-Version:") {
			email.Subject += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para Mime-Version:
		if strings.HasPrefix(line, "Mime-Version:") && !inBody {
			inSubject = false
			inMime = true
			email.MimeVersion = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inMime && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "Content-Type:") {
			email.MimeVersion += " " + line
			continue
		}

		//Evaluar cantidad de lineas para Content-Type:
		if strings.HasPrefix(line, "Content-Type:") && !inBody {
			inMime = false
			inContentType = true
			email.ContentType = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inContentType && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "Content-Transfer-Encoding:") {
			email.ContentType += " " + line
			continue
		}

		//Evaluar cantidad de lineas para Content-Transfer-Encoding
		if strings.HasPrefix(line, "Content-Transfer-Encoding:") && !inBody {
			inContentType = false
			inContentTransfer = true
			email.ContentTransferEncoding = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inContentTransfer && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "X-From:") {
			email.ContentTransferEncoding += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para X-From
		if strings.HasPrefix(line, "X-From:") && !inBody {
			inContentTransfer = false
			inXFrom = true
			email.XFrom = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inXFrom && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "X-To:") {
			email.XFrom += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para X-To
		if strings.HasPrefix(line, "X-To:") && !inBody {
			inXFrom = false
			inXTo = true
			email.XTo = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inXTo && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "X-cc:") {
			email.XTo += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para X-cc
		if strings.HasPrefix(line, "X-cc:") && !inBody {
			inXTo = false
			inXCc = true
			email.XCc = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inXCc && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "X-bcc:") {
			email.XCc += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para X-bcc
		if strings.HasPrefix(line, "X-bcc:") && !inBody {
			inXCc = false
			inXBcc = true
			email.XBcc = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inXBcc && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "X-Folder:") {
			email.XBcc += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para X-Folder
		if strings.HasPrefix(line, "X-Folder:") && !inBody {
			inXBcc = false
			inXFolder = true
			email.XFolder = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inXFolder && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "X-Origin:") {
			email.XFolder += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para X-Origin
		if strings.HasPrefix(line, "X-Origin") && !inBody {
			inXFolder = false
			inXOrigin = true
			email.XOrigin = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inXOrigin && len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "X-FileName:") {
			email.XOrigin += " " + strings.Replace(line, "\t", "", -1)
			continue
		}

		//Evaluar cantidad de lineas para X-FileName
		if strings.HasPrefix(line, "X-FileName:") && !inBody {
			inXOrigin = false
			inXFileName = true
			email.XFileName = strings.TrimSpace(strings.SplitAfter(line, ":")[1])
			continue
		}

		if inXFileName && len(strings.TrimSpace(line)) == 0 {
			inXFileName = false
			inBody = true
			continue
		}

		if inBody {

			email.Body += " " + line

			continue
		}

	} //fin del for

	return email

}
