package filtrorectangular

import (
	"math/cmplx"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/ft"
	"github.com/mjibson/go-dsp/fft"
)

type ItemFiltroRectangular struct {
	I               int
	TimeUs          float64
	Complx          complex128
	Hz              float64
	RectWindow      int
	ReqFreqSpectrum complex128
}

type FiltroRectangular []ItemFiltroRectangular

func NewFiltroRectangular(GraficaTotal ft.Fts, GraficasParciales ...ft.Fts) FiltroRectangular {
	filtro := make(FiltroRectangular, 0)

	for _, DataGTotal := range GraficaTotal {
		var reactWindow int
		var reqFreqSpectrum complex128

		for _, Graficas := range GraficasParciales {
			for _, item := range Graficas {
				if DataGTotal.Hz == item.Hz {
					reactWindow = 1
					reqFreqSpectrum = DataGTotal.Complx * complex(1, 0)
					break
				}
			}
		}

		item := ItemFiltroRectangular{
			I:               DataGTotal.I,
			TimeUs:          DataGTotal.TimeUs,
			Complx:          DataGTotal.Complx,
			Hz:              DataGTotal.Hz,
			RectWindow:      reactWindow,
			ReqFreqSpectrum: reqFreqSpectrum,
		}

		filtro = append(filtro, item)
	}

	return filtro
}

func (fr FiltroRectangular) DataFiltrada() []complex128 {
	x := make([]complex128, 0)

	for _, data := range fr {
		x = append(x, data.ReqFreqSpectrum)
	}

	filtro := fft.IFFT(x)

	return filtro
}

func (fr FiltroRectangular) GraficaFiltroToRender() *charts.Line {
	filtro := fr.DataFiltrada()
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for i := 0; i < len(filtro); i++ {
		x = append(x, fr[i].TimeUs)
		dataFiltro := filtro[i]
		items = append(items, opts.LineData{
			Value: real(dataFiltro),
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
			Name: "us",
		}),
	)

	line.SetXAxis(x).
		AddSeries("sensor", items)

	return line
}

func (fr FiltroRectangular) GraficaRectWindowToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range fr {
		x = append(x, float64(dato.Hz))
		items = append(items, opts.LineData{
			Value: dato.RectWindow,
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "RectWindow",
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

func (fr FiltroRectangular) GraficaReqFreqSpectrumToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range fr {
		r, _ := cmplx.Polar(dato.ReqFreqSpectrum)

		x = append(x, float64(dato.Hz))
		items = append(items, opts.LineData{
			Value: r,
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "ReqFreqSpectrum",
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
