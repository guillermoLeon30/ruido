package filtros

import (
	"math"
	"math/cmplx"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/ft"
)

type FiltroGauss2 struct {
	Freq_Hz       []float64
	FreqCorte_Hz  float64
	AnchoBanda_Hz float64
	Filtro        []float64
	Coef          []complex128
}

func NewFiltroGauss2(fourier ft.Fourier, freqCorte_Hz float64, anchoBanda_Hz float64) FiltroGauss2 {
	filtro := make([]float64, 0)
	coef := make([]complex128, 0)

	for i, f := range fourier.Freq_Hz {
		gauss := math.Exp((-1.0 / 2.0) * math.Pow((f-freqCorte_Hz)/(anchoBanda_Hz/2), 2))
		filtro = append(filtro, gauss)
		r, rad := cmplx.Polar(fourier.Coeff[i])
		r = r * gauss
		nuevoCoeff := cmplx.Rect(r, rad)
		coef = append(coef, nuevoCoeff)
	}

	fgauss := FiltroGauss2{
		Freq_Hz:       fourier.Freq_Hz,
		FreqCorte_Hz:  freqCorte_Hz,
		AnchoBanda_Hz: anchoBanda_Hz,
		Filtro:        filtro,
		Coef:          coef,
	}

	return fgauss
}

func (fs FiltroGauss2) GraficaFiltroToRender() *charts.Line {
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
			Title: "Filtro Gauss",
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

func (fs FiltroGauss2) GraficaEspectroFrecuencia() *charts.Line {
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
