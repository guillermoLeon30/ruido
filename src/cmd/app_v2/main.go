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
	"github.com/guillermoLeon30/ruido/src/domain/filtros"
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

	_, _, fc := fourier.Calcular(8000, 13000)
	bw := 2500 // BW mediante inspeccion de grafica (Hz)

	fmt.Println("fc: ", fc)
	fmt.Println("BW2: ", bw)

	fRec2 := filtros.NewFiltroRectangular2(fourier, fc, float64(bw))
	fourier.Inversa(fRec2.Coef)

	// Grafica Rectangular
	grafica := filepath.Join(dir, "out", "res_filtro_rect_v2.html")
	err = crearGraficas(
		grafica,
		datos.GraficaToRender(),
		fourier.GraficaTotalToRender(),
		fRec2.GraficaFiltroToRender(),
		fRec2.GraficaEspectroFrecuencia(),
		fourier.GraficaInversaToRender(),
	)
	if err != nil {
		panic(err)
	}

	// Grafica Gauss
	fGauss2 := filtros.NewFiltroGauss2(fourier, fc, float64(bw))
	fourier.Inversa(fGauss2.Coef)

	grafica = filepath.Join(dir, "out", "res_filtro_gauss_v2.html")
	err = crearGraficas(
		grafica,
		datos.GraficaToRender(),
		fourier.GraficaTotalToRender(),
		fGauss2.GraficaFiltroToRender(),
		fGauss2.GraficaEspectroFrecuencia(),
		fourier.GraficaInversaToRender(),
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
