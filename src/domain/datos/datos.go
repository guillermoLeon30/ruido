package datos

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type Dato struct {
	TimeUs float64
	Value  float64
}

type Datos []Dato

func NewDatos(path string) (Datos, error) {
	archivo, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	datos := make(Datos, 0)
	datosSensor := strings.Split(string(archivo), ",")

	for i := 0; i < len(datosSensor); i++ {
		var valor float64
		if i == len(datosSensor)-1 {
			svalor := strings.Split(datosSensor[i], "\n")[0]
			valor, err = strconv.ParseFloat(svalor, 64)
			if err != nil {
				return nil, err
			}
		} else {
			valor, err = strconv.ParseFloat(datosSensor[i], 64)
			if err != nil {
				return nil, err
			}
		}

		sensor := Dato{
			TimeUs: 12.5 * float64(i+1),
			Value:  valor,
		}

		datos = append(datos, sensor)
	}

	return datos, nil
}

func (ds Datos) Graficar(path string) error {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range ds {
		x = append(x, dato.TimeUs)
		items = append(items, opts.LineData{Value: dato.Value})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Data Sensor",
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

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	line.Render(f)

	return nil
}

func (ds Datos) GraficaToRender() *charts.Line {
	x := make([]float64, 0)
	items := make([]opts.LineData, 0)

	for _, dato := range ds {
		x = append(x, dato.TimeUs)
		items = append(items, opts.LineData{Value: dato.Value})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Data Sensor",
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

func (ds Datos) GetValues() []float64 {
	values := make([]float64, 0)

	for _, dato := range ds {
		values = append(values, dato.Value)
	}

	return values
}
