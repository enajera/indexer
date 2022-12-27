package process

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

// ProcesarArchivo mapea el contenido de un archivo de correo electrónico a una estructura de datos de tipo Email
func ProcesarArchivo(path string) (Email, error) {

	//fmt.Println(path)
	// Abre el archivo en la ruta especificada
	file, err := os.Open(path)
	if err != nil {
		return Email{}, fmt.Errorf("Error al recorrer el archivo: %s", path)
	}
	defer file.Close() // Cierra el archivo al finalizar la función

	// Mapea el contenido del archivo a una estructura de datos de tipo Email
	email := MapearCorreos(file)

	//IndexarCorreos(*email)

	return *email, nil

	
	
}


// IndexarCorreos indexa una slice de correos en Zincsearch de manera simple
func IndexarCorreos(email Email) error {

	// Convierte el correo a formato JSON
	emailJSON, err := json.Marshal(email)
	if err != nil {
		return err
	}

	data := string(emailJSON)

	req, err := http.NewRequest("POST", "http://localhost:4080/api/enronmail/_doc", strings.NewReader(data))
	if err != nil {
		return fmt.Errorf("Error:%s", err)
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
	fmt.Println(resp.StatusCode, "-", string(email.XFolder))

	return nil

}
