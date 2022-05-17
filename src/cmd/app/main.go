package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/guillermoLeon30/ruido/src/domain/datos"
	filtrorectangular "github.com/guillermoLeon30/ruido/src/domain/filtro_rectangular"
	"github.com/guillermoLeon30/ruido/src/domain/ft"
	"github.com/jedib0t/go-pretty/v6/table"
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
	dataFt := ft.NewFt(datos)

	// Calculos de frecuencias
	_, freq1_1, freq2_1, freqc_1 := dataFt.Calcular(8700, 11500)
	_, freq1_2, freq2_2, freqc_2 := dataFt.Calcular(68500, 71700)
	imprimirCalculosFrecuencias(8700, 11500, freq1_1, freq2_1, freqc_1)
	fmt.Println("")
	imprimirCalculosFrecuencias(68500, 71700, freq1_2, freq2_2, freqc_2)
	fmt.Println("")

	// Separacion de graficas
	dataFreq1 := dataFt.GetPartialWithFreq(freq1_1.Hz, freq2_1.Hz)
	dataFreq2 := dataFt.GetPartialWithFreq(freq1_2.Hz, freq2_2.Hz)
	imprimirDataFrecuancias(dataFreq1)
	fmt.Println("")
	imprimirDataFrecuancias(dataFreq2)

	dataFiltroRect := filtrorectangular.NewFiltroRectangular(dataFt, dataFreq1, dataFreq2)
	grafica := filepath.Join(dir, "out", "res_filtro_rect.html")
	err = crearGraficas(
		grafica,
		datos.GraficaToRender(),
		dataFt.GraficaManitudToRender(),
		dataFiltroRect.GraficaRectWindowToRender(),
		dataFiltroRect.GraficaReqFreqSpectrumToRender(),
		dataFiltroRect.GraficaFiltroToRender(),
	)
	if err != nil {
		panic(err)
	}
}

func imprimirCalculosFrecuencias(f1 float64, f2 float64, fq1 ft.Ft, fq2 ft.Ft, fqc ft.Ft) {
	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}
	t := table.NewWriter()
	t.AppendHeader(table.Row{"", "Frecuencia (Hz)", "Frecuencia (Hz)", "Frecuencia (Hz)", "Frecuencia (Hz)", "Frecuencia (Hz)", "Frecuencia (Hz)"}, rowConfigAutoMerge)
	t.AppendHeader(table.Row{"Id", "Inicial", "Final", "f1", "f2", "fo", "BW"})
	t.AppendRow(table.Row{
		"1",
		fmt.Sprintf("%.2f", f1),
		fmt.Sprintf("%.2f", f2),
		fmt.Sprintf("%.2f", fq1.Hz),
		fmt.Sprintf("%.2f", fq2.Hz),
		fmt.Sprintf("%.2f", fqc.Hz),
		fmt.Sprintf("%.2f", fq2.Hz-fq1.Hz),
	})

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredDark)
	t.Render()
}

func imprimirDataFrecuancias(datos ft.Fts) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"Id", "Magnitud", "Angulo (ª)", "Frecuencia (HZ)"})
	for _, dato := range datos {
		t.AppendRow(table.Row{
			fmt.Sprintf("%.2d", dato.I),
			fmt.Sprintf("%.2f", dato.Magnitud),
			fmt.Sprintf("%.2f", dato.Angulo),
			fmt.Sprintf("%.2f", dato.Hz),
		})
	}

	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleColoredDark)
	t.Render()
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
