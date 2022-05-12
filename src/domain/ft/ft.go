package ft

import (
	"math/cmplx"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/datos"
	"github.com/mjibson/go-dsp/fft"
)

type Ft struct {
	I        int
	Magnitud float64
	Angulo   float64
}

type Fts []Ft

func NewFt(datos datos.Datos) Fts {
	a := make([]float64, len(datos))
	for i := 0; i < len(datos); i++ {
		a[i] = datos[i].Value
	}

	x := fft.FFTReal(a)

	fts := make(Fts, 0)
	for i := 0; i < len(datos); i++ {
		r, teta := cmplx.Polar(x[i])

		ft := Ft{
			I:        i,
			Magnitud: r,
			Angulo:   teta,
		}

		fts = append(fts, ft)
	}

	return fts
}

func (fs Fts) GraficarManitud(path string) error {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range fs {
		x = append(x, float64(dato.I))
		items = append(items, opts.LineData{Value: dato.Magnitud})
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
