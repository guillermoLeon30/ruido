package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
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

	// Fourier
	periodo := 1 / (12.5 * math.Pow10(-6))
	fourier := ft.NewFourier(datos, periodo)

	f1, f2, fc := fourier.Calcular(8000, 13000)
	bw := 2500 // BW mediante inspeccion de grafica (Hz)

	fmt.Println("f1: ", f1)
	fmt.Println("f2: ", f2)
	fmt.Println("fc: ", fc)
	fmt.Println("BW: ", f2-f1)
	fmt.Println("BW2: ", bw)

	// Graficar
	grafica := filepath.Join(dir, "out", "res_filtro_rect_v2.html")
	err = crearGraficas(
		grafica,
		datos.GraficaToRender(),
		fourier.GraficaTotalToRender(),
	)
	if err != nil {
		panic(err)
	}
}

func crearGraficas(path string, graficas ...*charts.Line) error {
	page := components.NewPage()

	for _, grafica := range graficas {
		page.AddCharts(grafica)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	page.Render(io.MultiWriter(f))

	return nil
}
