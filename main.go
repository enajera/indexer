package main

import (
	"fmt"
	"os"
	"time"
	"github.com/enajera/indexer/zincsearch"
)

func main() { 

    //inicia tiempo de ejecucion
	start := time.Now()


	//Si solo viene un argumento, lanza error
	if len(os.Args) < 2 {
		fmt.Println("Error: no se ha especificado el nombre del archivo de base de datos")
		return
	}

	archivo := os.Args[1]
	//Se llama al procesador del archivo
	if err := zincsearch.ProcesarArchivo(archivo); err != nil {
		fmt.Println(err)
	}


	//finaliza tiempo de ejecucion
	elapsed := time.Since(start)
	fmt.Printf("Tiempo de ejecución: %s\n", elapsed)
}
