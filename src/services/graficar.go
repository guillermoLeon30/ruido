package services

import (
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/guillermoLeon30/ruido/src/domain/datos"
)

func Graficar(datos []datos.Dato, path string) error {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range datos {
		x = append(x, dato.TimeUs)
		items = append(items, opts.LineData{Value: dato.Value})
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
