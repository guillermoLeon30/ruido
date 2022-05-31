package filtros

import (
	"math/cmplx"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/ft"
)

type FiltroRectangular2 struct {
	Freq_Hz       []float64
	FreqCorte_Hz  float64
	AnchoBanda_Hz float64
	Filtro        []float64
	Coef          []complex128
}

func NewFiltroRectangular2(fourier ft.Fourier, freqCorte_Hz float64, anchoBanda_Hz float64) FiltroRectangular2 {
	filtro := make([]float64, 0)
	coef := make([]complex128, 0)
	freq_minima := freqCorte_Hz - anchoBanda_Hz/2
	freq_maxima := freqCorte_Hz + anchoBanda_Hz/2

	for i, f := range fourier.Freq_Hz {
		if f >= freq_minima && f <= freq_maxima {
			filtro = append(filtro, 1)
			coef = append(coef, fourier.Coeff[i])
		} else {
			filtro = append(filtro, 0)
			coef = append(coef, complex(0, 0))
		}
	}

	fRec := FiltroRectangular2{
		Freq_Hz:       fourier.Freq_Hz,
		FreqCorte_Hz:  freqCorte_Hz,
		AnchoBanda_Hz: anchoBanda_Hz,
		Filtro:        filtro,
		Coef:          coef,
	}

	return fRec
}

func (fs FiltroRectangular2) GraficaFiltroToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for i := len(fs.Freq_Hz) - 1; i >= 0; i-- {
		x = append(x, -fs.Freq_Hz[i])
		items = append(items, opts.LineData{
			Value: fs.Filtro[i],
		})
	}

	for i, f := range fs.Freq_Hz {
		x = append(x, f)
		items = append(items, opts.LineData{
			Value: fs.Filtro[i],
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Filtro Rectangular",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type: "inside",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type: "slider",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Hz",
		}),
	)

	line.SetXAxis(x).
		AddSeries("sensor", items)

	return line
}

func (fs FiltroRectangular2) GraficaEspectroFrecuencia() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for i := len(fs.Freq_Hz) - 1; i >= 0; i-- {
		x = append(x, -fs.Freq_Hz[i])
		items = append(items, opts.LineData{
			Value: cmplx.Abs(fs.Coef[i]),
		})
	}

	for i, f := range fs.Freq_Hz {
		x = append(x, f)
		items = append(items, opts.LineData{
			Value: cmplx.Abs(fs.Coef[i]),
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Espectro de Frecuencia",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type: "inside",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type: "slider",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Hz",
		}),
	)

	line.SetXAxis(x).
		AddSeries("sensor", items)

	return line
}
