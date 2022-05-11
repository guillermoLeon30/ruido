package main

import (
	"os"
	"path/filepath"

	"github.com/guillermoLeon30/ruido/src/domain/datos"
)

func main() {
	// Get path
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Get Datos
	archivoCsv := filepath.Join(dir, "data", "TOF_onereflect.csv")
	datos, err := datos.NewDatos(archivoCsv)
	if err != nil {
		panic(err)
	}

	// Graficar
	graficaSvg := filepath.Join(dir, "out", "grafica.html")
	err = datos.Graficar(graficaSvg)
	if err != nil {
		panic(err)
	}
}
