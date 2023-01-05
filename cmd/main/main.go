package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/enajera/indexer/internal/process"
)

func main() {

	banner, err := ioutil.ReadFile("banner.txt")
	if err != nil {
		fmt.Println("Error al abrir el banner.txt")
	}
	
	fmt.Println(string(banner))
	fmt.Println()

	// Si solo viene un argumento, lanza error
	if len(os.Args) < 2 {
		fmt.Println("Error: no se ha especificado el nombre del archivo de base de datos")
		return
	}

	file := os.Args[1]

	var wg sync.WaitGroup
	// Servidor de profiling
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wg.Add(1) //pprof - so we won't exit prematurely
	wg.Add(1) //run method
	go Run(file, &wg)
	wg.Wait()

}

func Run(file string, wg *sync.WaitGroup) {
	defer wg.Done()
	// Inicia tiempo de ejecución
	start := time.Now()
	// Obtiene las rutas de todos los archivos de la carpeta principal y sus subcarpetas
	filePaths, err := getFilePaths(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	//Procesamiento de archivos
	correos, err := process.ProcesarArchivo(filePaths)
	if err != nil {
		fmt.Printf("Error al procesar: %s", filePaths)

	}

	// Finaliza tiempo de procesamiento
	process := time.Since(start)
	fmt.Printf("Archivos procesados: %d\n", len(correos))
	fmt.Printf("Tiempo de procesamiento de correos: %s\n", process)

	fmt.Println("------------------------------------------------")
	fmt.Printf("Indexando correos en ZincSearch... \n")
	indexing := time.Now()
	//Enviar correos a ZincSearch
	//IndexarCorreosMasivos(correos)

	fmt.Println("------------------------------------------------")
	endexing := time.Since(indexing)
	fmt.Printf("Tiempo de indexado de correos: %s\n", endexing)
	fmt.Println("------------------------------------------------")
	end := time.Since(start)
	fmt.Printf("Tiempo Total: %s\n", end)

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
func IndexarCorreosMasivos(correos any) error {

	// Crea una estructura que se ajusta al formato Json que pide ZincSearch
	data := struct {
		Index   string `json:"index"`
		Records any    `json:"records"`
	}{
		Index:   "enronmails_bulk",
		Records: correos,
	}

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
