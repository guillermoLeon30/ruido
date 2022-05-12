package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guillermoLeon30/ruido/src/domain/datos"
	"github.com/guillermoLeon30/ruido/src/domain/ft"
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

	// Fourier
	dataFt := ft.NewFt(datos)
	for _, data := range dataFt {
		fmt.Println(data.I, ": ", data.Magnitud, " - ", data.Angulo)
	}

	// Graficar Fourier
	graficaF := filepath.Join(dir, "out", "fourier.html")
	err = dataFt.GraficarManitud(graficaF)
	if err != nil {
		panic(err)
	}
}
