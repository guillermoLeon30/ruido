package ft

import (
	"math"
	"math/cmplx"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/datos"
	"github.com/mjibson/go-dsp/fft"
)

type Ft struct {
	I        int
	Magnitud float64
	Angulo   float64
	Hz       float64
}

type Fts []Ft

func NewFt(datos datos.Datos) Fts {
	a := make([]float64, len(datos))
	for i := 0; i < len(datos); i++ {
		a[i] = datos[i].Value
	}

	f := (1 / datos[len(datos)-1].TimeUs) * 1000000
	f = math.Round(f*100) / 100
	x := fft.FFTReal(a)

	fts := make(Fts, 0)
	for i := 0; i < len(datos); i++ {
		r, rad := cmplx.Polar(x[i])

		grados := (360 * rad) / (2 * math.Pi)
		angulo := math.Round(grados*100) / 100
		hz := float64(i) * f
		hz = math.Round(hz*100) / 100

		ft := Ft{
			I:        i,
			Magnitud: r,
			Angulo:   angulo,
			Hz:       hz,
		}

		fts = append(fts, ft)
	}

	return fts
}

func (fs Fts) GraficarManitud(path string) error {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range fs {
		x = append(x, float64(dato.Hz))
		items = append(items, opts.LineData{
			Value: dato.Magnitud,
		})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Fourier",
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

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	line.Render(f)

	return nil
}

func (fs Fts) GraficarAngulo(path string) error {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range fs {
		x = append(x, float64(dato.Hz))
		items = append(items, opts.LineData{Value: dato.Angulo})
	}

	line := charts.NewLine()
	line.SetXAxis(x).
		AddSeries("sensor", items)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	line.Render(f)

	return nil
}

func (fs Fts) AltasMagnitudes(freqIni float64, freqFin float64) Fts {
	fts := make(Fts, 0)

	for _, data := range fs {
		if data.Hz > freqIni && data.Hz < freqFin {
			fts = append(fts, data)
		}
	}

	sort.SliceStable(fts, func(i, j int) bool {
		return fts[i].Magnitud > fts[j].Magnitud
	})

	return fts
}
