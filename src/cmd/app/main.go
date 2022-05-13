package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/guillermoLeon30/ruido/src/domain/datos"
	"github.com/guillermoLeon30/ruido/src/domain/ft"
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

	// Graficar
	graficaSvg := filepath.Join(dir, "out", "grafica.html")
	err = datos.Graficar(graficaSvg)
	if err != nil {
		panic(err)
	}

	// Fourier
	dataFt := ft.NewFt(datos)
	// for _, data := range dataFt {
	// 	fmt.Println(data.I, ": ", data.Magnitud, " - ", data.Angulo, "f: ", data.Hz)
	// }

	// Graficar Fourier Magnitud
	graficaFM := filepath.Join(dir, "out", "f_magnitud.html")
	err = dataFt.GraficarManitud(graficaFM)
	if err != nil {
		panic(err)
	}

	// Graficar Fourier Angulo
	graficaFA := filepath.Join(dir, "out", "f_angulo.html")
	err = dataFt.GraficarAngulo(graficaFA)
	if err != nil {
		panic(err)
	}

	// Maximos
	maxs1 := dataFt.AltasMagnitudes(8700, 11500)
	fmt.Println("max1 -> ", "M: ", maxs1[0].Magnitud, " f: ", maxs1[0].Hz, "Hz")

	maxs2 := dataFt.AltasMagnitudes(68500, 71700)
	fmt.Println("max2 -> ", "M: ", maxs2[0].Magnitud, " f: ", maxs2[0].Hz, "Hz")
}
