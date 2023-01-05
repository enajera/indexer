package process


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

var archivosProcesados int
var totalArchivos int

// ProcesarArchivo mapea el contenido de un archivo de correo electrónico a una estructura de datos de tipo Email
func ProcesarArchivo(filePaths []string) ([]Email, error) {

	totalArchivos = len(filePaths)
	archivosProcesados = 0

	//recorre todos los archivos
	correos := []Email{}
	for _, filePath := range filePaths {

		// Mapea el contenido del archivo a una estructura de datos de tipo Email
		correos = append(correos, *MapearCorreos(filePath))
       // archivosProcesados+=1
		//fmt.Printf("Progreso:[%d/%d]\n", archivosProcesados, totalArchivos)
	}

	return correos, nil

}
