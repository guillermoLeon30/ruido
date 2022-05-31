package ft

import (
	"math"
	"math/cmplx"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/datos"
	"gonum.org/v1/gonum/dsp/fourier"
)

type Fourier struct {
	Fft     *fourier.FFT
	Coeff   []complex128
	Seq     []float64
	Periodo float64
	Freq_Hz []float64
}

func NewFourier(datos datos.Datos, periodo float64) Fourier {
	samples := make([]float64, len(datos))
	for i := 0; i < len(datos); i++ {
		samples[i] = datos[i].Value
	}

	fft := fourier.NewFFT(len(samples))
	coeff := fft.Coefficients(nil, samples)

	frecuencias := make([]float64, 0)
	for i := range coeff {
		frecuencia := fft.Freq(i) * periodo
		frecuencia = math.Round(frecuencia*100) / 100

		frecuencias = append(frecuencias, frecuencia)
	}

	return Fourier{
		Fft:     fft,
		Coeff:   coeff,
		Periodo: periodo,
		Freq_Hz: frecuencias,
	}
}

// Devuelve f1, f2 y fc
func (fs Fourier) Calcular(freqIni float64, freqFin float64) (float64, float64, float64) {
	var f1 float64
	var f2 float64
	var fc float64

	fCoeff := make([]complex128, 0)
	freq_indices := make([]int, 0)
	for i, f := range fs.Freq_Hz {
		if f >= freqIni && f <= freqFin {
			freq_indices = append(freq_indices, i)
		}
	}

	for _, i := range freq_indices {
		mag := fs.Coeff[i]
		fCoeff = append(fCoeff, mag)
	}

	sort.SliceStable(fCoeff, func(i, j int) bool {
		return cmplx.Abs(fCoeff[i]) > cmplx.Abs(fCoeff[j])
	})

	maximoVal := cmplx.Abs(fCoeff[0])
	limiteVal := maximoVal * 0.7

	for _, i := range freq_indices {
		if cmplx.Abs(fs.Coeff[i]) >= limiteVal {
			f1 = fs.Freq_Hz[i]
			break
		}
	}

	for i := len(freq_indices) - 1; i >= 0; i-- {
		if cmplx.Abs(fs.Coeff[freq_indices[i]]) >= limiteVal {
			f2 = fs.Freq_Hz[freq_indices[i]]
			break
		}
	}

	fc = (f1 + f2) / 2

	return f1, f2, fc
}

func (fs *Fourier) Inversa(coeff []complex128) {
	var originCoeff []complex128

	if coeff == nil {
		originCoeff = append(originCoeff, fs.Coeff...)
	} else {
		originCoeff = append(originCoeff, coeff...)
	}

	seq := fs.Fft.Sequence(nil, originCoeff)

	fs.Seq = seq
}

func (fs Fourier) GraficaParcialToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	x = append(x, fs.Freq_Hz...)
	for _, dato := range fs.Coeff {
		items = append(items, opts.LineData{
			Value: cmplx.Abs(dato),
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Fourier Magnitud",
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

func (fs Fourier) GraficaTotalToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for i := len(fs.Freq_Hz) - 1; i >= 0; i-- {
		x = append(x, -fs.Freq_Hz[i])
		items = append(items, opts.LineData{
			Value: cmplx.Abs(fs.Coeff[i]),
		})
	}

	x = append(x, fs.Freq_Hz...)
	for _, dato := range fs.Coeff {
		items = append(items, opts.LineData{
			Value: cmplx.Abs(dato),
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Fourier Magnitud",
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

func (fs Fourier) GraficaInversaToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	// x = append(x, fs.Freq_Hz...)
	for i, dato := range fs.Seq {
		tiempo := (1 / fs.Periodo) * float64(i+1) * 1000000
		tiempo = math.Round(tiempo*100) / 100
		x = append(x, tiempo)
		items = append(items, opts.LineData{
			Value: dato / float64(len(fs.Seq)),
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Inversa Fourier",
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
