package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/enajera/indexer/process"
)

func main() {
	// Inicia tiempo de ejecución
	start := time.Now()

	// Si solo viene un argumento, lanza error
	if len(os.Args) < 2 {
		fmt.Println("Error: no se ha especificado el nombre del archivo de base de datos")
		return
	}

	// Servidor de profiling
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	file := os.Args[1]

	// Obtiene las rutas de todos los archivos de la carpeta principal y sus subcarpetas
	filePaths, err := getFilePaths(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Printf("Paths: %d",len(filePaths))

	correos := []process.Email{}

	// Crea una goroutine por cada archivo a procesar
	var wg sync.WaitGroup
	for i, filePath := range filePaths {
		wg.Add(i)
		go func(path string) {
			// Procesa el archivo y cierra el WaitGroup al finalizar
			email, err := process.ProcesarArchivo(path)
			if err != nil {
				fmt.Printf("Error al procesar el archivo: %s", path)
				return
			}
			correos = append(correos, email)
			wg.Done()
		}(filePath)
	}

	// // Espera a que todas las goroutines hayan terminado su trabajo
	wg.Wait()

	//Crea un archivo json llamado correos,json
	jsonFile, err := os.Create("correos.json")
	if err != nil {
		fmt.Printf("Error al crear el archivo: correos.json")
	}
	defer jsonFile.Close()

	//Crea una estructura que se ajusta al formato a subir
	data := struct {
		Index   string          `json:"index"`
		Records []process.Email `json:"records"`
	}{
		Index:   "enronmails_bulk",
		Records: correos,
	}

	if err := json.NewEncoder(jsonFile).Encode(&data); err != nil {
		fmt.Printf("Error al escribir el archivo: correos.json")
	}

	IndexarCorreosMasivos(data)

	// Finaliza tiempo de ejecución
	elapsed := time.Since(start)
	fmt.Printf("Archivos procesados: %d\n", len(correos))
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)

}

// getFilePaths obtiene las rutas de todos los archivos de la carpeta principal y sus subcarpetas
func getFilePaths(root string) ([]string, error) {

	var paths []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error al recorrer el archivo: %s", err)
	}
	return paths, nil
}

// IndexarCorreos indexa una slice de correos en Zincsearch de manera compleja
func IndexarCorreosMasivos(data any) error {

	emailJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://localhost:4080/api/_bulkv2", bytes.NewBuffer(emailJSON))
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
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(body))
	fmt.Println(resp.StatusCode, "-", string(body))

	return nil

}
