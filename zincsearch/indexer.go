package zincsearch

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

)

type Email struct {
	MessageID               string
	Date                    string
	From                    string
	To                      string
	Subject                 string
	MimeVersion             string
	ContentType             string
	ContentTransferEncoding string
	XFrom                   string
	XTo                     string
	XCc                     string
	XBcc                    string
	XFolder                 string
	XOrigin                 string
	XFileName               string
	Body                    string
}

// ProcesarArchivo recorre los archivos de correo electrónico de la base de datos especificada
func ProcesarArchivo(database string) error {
	
	emails := []Email{}
	// Recorrer las carpetas y archivos de la base de datos
	if err := filepath.Walk(filepath.Join(database), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			//Es carpeta
		} else {
			//fmt.Println("Archivo: " + path)
			// Abre el archivo en la ruta especificada
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close() // Cierra el archivo al finalizar la función

			// Crea un nuevo escaner para leer línea a línea del archivo
			scanner := bufio.NewScanner(file)

			email := &Email{}

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

			// Añade el correo electrónico a la lista de correos
			// emailJSON, err := json.Marshal(email)
			// if err != nil {
			// 	panic(err)
			// }

			//data := string(emailJSON)
			//SubirInfo(data)

			emails = append(emails, *email)

		}
		return nil

	}); err != nil {
		return fmt.Errorf("error al indexar la base de datos: %v", err)
	}

	// // Reemplaza caracateres invalidos
	// buf := &bytes.Buffer{}
	// enc := json.NewEncoder(buf)
	// enc.SetEscapeHTML(false)
	// err := enc.Encode(emails)
	// if err != nil {
	// 	panic(err)
	// }

	// // Escribe el resultado en el archivo "emails.json"
	// err = os.WriteFile("emails.json", buf.Bytes(), 0644)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println("Archivos procesados:", len(emails))

	
	return nil
}

// Indexador a Zincsearch
func SubirInfo(data string) {

	req, err := http.NewRequest("POST", "http://localhost:4080/api/enronmail/_doc", strings.NewReader(data))
	if err != nil {
		panic(err)
	}
	req.SetBasicAuth("admin", "Complexpass#123")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//fmt.Println(resp.StatusCode)
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(body))
}
