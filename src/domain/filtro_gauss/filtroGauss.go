package filtrogauss

import (
	"math"
	"math/cmplx"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/ft"
	"github.com/mjibson/go-dsp/fft"
)

type ItemFiltroGauss struct {
	I                 int
	TimeUs            float64
	Complx            complex128
	Hz                float64
	GaussWindow       float64
	GaussFreqSpectrum complex128
}

type FiltroGauss []ItemFiltroGauss

func NewFiltroGauss(GraficaTotal ft.Fts,
	grafica1 ft.Fts,
	grafica2 ft.Fts,
	f1_1 ft.Ft,
	f2_1 ft.Ft,
	fc_1 ft.Ft,
	f1_2 ft.Ft,
	f2_2 ft.Ft,
	fc_2 ft.Ft) FiltroGauss {
	filtro := make(FiltroGauss, 0)

	for _, DataGTotal := range GraficaTotal {
		var gaussWindow float64
		var gaussFreqSpectrum complex128
		var tomada int

		for _, item := range grafica1 {
			if DataGTotal.Hz == item.Hz {
				gaussWindow = math.Exp(-1.0 / 2.0 * math.Pow((DataGTotal.Hz-fc_1.Hz)/((f2_1.Hz-f1_1.Hz)/2), 2))
				r, rad := cmplx.Polar(DataGTotal.Complx)
				r2 := r * gaussWindow
				gaussFreqSpectrum = cmplx.Rect(r2, rad)

				item := ItemFiltroGauss{
					I:                 DataGTotal.I,
					TimeUs:            DataGTotal.TimeUs,
					Complx:            DataGTotal.Complx,
					Hz:                DataGTotal.Hz,
					GaussWindow:       gaussWindow,
					GaussFreqSpectrum: gaussFreqSpectrum,
				}

				filtro = append(filtro, item)
				tomada = 1
				break
			}
		}

		for _, item := range grafica2 {
			if DataGTotal.Hz == item.Hz {
				gaussWindow = math.Exp(-1.0 / 2.0 * math.Pow((DataGTotal.Hz-fc_2.Hz)/((f2_2.Hz-f1_2.Hz)/2), 2))
				r, rad := cmplx.Polar(DataGTotal.Complx)
				r2 := r * gaussWindow
				gaussFreqSpectrum = cmplx.Rect(r2, rad)

				item := ItemFiltroGauss{
					I:                 DataGTotal.I,
					TimeUs:            DataGTotal.TimeUs,
					Complx:            DataGTotal.Complx,
					Hz:                DataGTotal.Hz,
					GaussWindow:       gaussWindow,
					GaussFreqSpectrum: gaussFreqSpectrum,
				}

				filtro = append(filtro, item)
				tomada = 1
				break
			}
		}

		if tomada == 1 {
			tomada = 0
		} else {
			item := ItemFiltroGauss{
				I:                 DataGTotal.I,
				TimeUs:            DataGTotal.TimeUs,
				Complx:            DataGTotal.Complx,
				Hz:                DataGTotal.Hz,
				GaussWindow:       gaussWindow,
				GaussFreqSpectrum: gaussFreqSpectrum,
			}

			filtro = append(filtro, item)
		}
	}

	return filtro
}

func (ftr FiltroGauss) DataFiltrada() []complex128 {
	x := make([]complex128, 0)

	for _, data := range ftr {
		x = append(x, data.GaussFreqSpectrum)
	}

	filtro := fft.IFFT(x)

	return filtro
}

func (ftr FiltroGauss) GraficaFiltroToRender() *charts.Line {
	filtro := ftr.DataFiltrada()
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for i := 0; i < len(filtro); i++ {
		x = append(x, ftr[i].TimeUs)
		dataFiltro := filtro[i]
		items = append(items, opts.LineData{
			Value: real(dataFiltro),
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
			Name: "us",
		}),
	)

	line.SetXAxis(x).
		AddSeries("sensor", items)

	return line
}

func (ftr FiltroGauss) GraficarGaussWindowToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range ftr {
		x = append(x, float64(dato.Hz))
		items = append(items, opts.LineData{
			Value: dato.GaussWindow,
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "GaussWindow",
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

func (ftr FiltroGauss) GraficarGaussFreqSpectrumToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range ftr {
		r, _ := cmplx.Polar(dato.GaussFreqSpectrum)

		x = append(x, float64(dato.Hz))
		items = append(items, opts.LineData{
			Value: r,
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "GaussFreqSpectrum",
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
