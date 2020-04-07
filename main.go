package main

import (
	"github.com/chfenger/goNum"
	"github.com/spf13/cobra"
	"github.com/tealeg/xlsx"
	"log"
)

func excelParse(fileName string) ([][]float64, error) {
	xlFile, err := xlsx.OpenFile(fileName)

	if err != nil {
		return nil, err
	}
	//获取行数
	rowCount := len(xlFile.Sheets[0].Rows)
	log.Printf("rowCount is %d", rowCount)

	//开辟除表头外的行数的数组内存
	resourceArr := make([][]float64, rowCount-2)

	//遍历每一行
	for rowIndex, row := range xlFile.Sheets[0].Rows {

		cellCount := len(row.Cells)
		//log.Printf("cellCount is %d", cellCount)

		//跳过第前行表头信息
		if rowIndex < 2 {
			continue
		}
		//遍历每一个单元

		for cellIndex, cell := range row.Cells {
			if cellIndex < 2 {
				continue
			}
			value, _ := cell.Float()

			if resourceArr[rowIndex-2] == nil {
				resourceArr[rowIndex-2] = make([]float64, cellCount-2)
			}

			resourceArr[rowIndex-2][cellIndex-2] = value

		}
	}

	return resourceArr, nil
}

func reverse(arr []float64) {
	length := len(arr)
	for i := 0; i < length/2; i++ {
		temp := arr[length-1-i]
		arr[length-1-i] = arr[i]
		arr[i] = temp
	}
}

func main() {

	file := "./YHJ-1875-VT.xlsx"
	power := 4

	root := &cobra.Command{
		Use:   "fitter",
		Short: "fitter",
		Long:  "fitter ratio and power",

		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("main running")

			if file == "" {
				log.Printf("no input file")
				return
			}

			arr, err := excelParse(file)

			if err != nil {
				log.Printf("no input file")
				return
			}

			//log.Printf("%+v", arr)

			ratio := make([][]float64, power+1)

			log.Printf("size arr: %d, size arr[0]: %d", len(arr), len(arr[0]))
			//算全部系数
			for i := 0; i < len(arr[0]); i++ {
				var params []float64
				for j := 0; j < len(arr); j++ {
					params = append(params, float64(j))
					params = append(params, arr[j][i])
				}
				xy := goNum.NewMatrix(len(params)/2, 2, params)

				if out, _, _, ok := goNum.FittingPolynomial(xy, power); ok == true {
					reverse(out.Data)
					for i := 0; i <= power; i++ {
						ratio[i] = append(ratio[i], out.Data[i])
					}
				} else {
					log.Printf("err is %v", ok)
				}
			}

			for i := 0; i <= power; i++ {
				log.Printf("ratio[%d] is %+v", power-i, ratio[i])
				log.Printf("------------------------------------")
			}

			for i := 0; i <= power; i++ {
				var params []float64
				for j := 0; j < len(ratio[i]); j++ {
					params = append(params, float64(j))
					params = append(params, ratio[i][j])
				}

				//log.Printf("paramsA: %+v", paramsA)

				xy := goNum.NewMatrix(len(params)/2, 2, params)

				if out, _, _, ok := goNum.FittingPolynomial(xy, power); ok == true {
					reverse(out.Data)
					log.Printf("ratio[%d]: out: %+v", power-i, out.Data)
				} else {
					log.Printf("err is %v", ok)
				}

			}

		},
	}

	root.PersistentFlags().StringVar(&file, "file", file, "--file ./YHJ-1875-VT.xlsx")
	root.PersistentFlags().IntVar(&power, "power", power, "--power 4")

	root.Execute()

}
