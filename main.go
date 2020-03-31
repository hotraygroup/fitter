package main

import (
	"github.com/chfenger/goNum"
	"github.com/tealeg/xlsx"
	"log"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func excelParse(fileName string) [][]float64 {
	xlFile, err := xlsx.OpenFile(fileName)
	checkErr(err)

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
			// for _, cell := range row.Cells {
			//  text := cell.String()
			//  fmt.Printf("%s\n", text)
			// }
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

	return resourceArr
}

func main() {

	log.Printf("main running")
	arr := excelParse("YHJ-1875-VT.xlsx")
	//log.Printf("%+v", arr)

	var A, B, C []float64

	for i := 0; i < len(arr); i++ {
		var params []float64
		for j := 0; j < len(arr[i]); j++ {
			params = append(params, float64(j))
			params = append(params, arr[j][i])
		}
		xy := goNum.NewMatrix(len(params)/2, 2, params)

		if out, _, _, ok := goNum.FittingPolynomial(xy, 2); ok == true {
			//log.Printf("out: %+v", out)
			A = append(A, out.Data[2])
			B = append(B, out.Data[1])
			C = append(C, out.Data[0])
		} else {
			log.Printf("err is %v", ok)
		}
	}

	log.Printf("A: %+v", A)
	log.Printf("-------------------")
	log.Printf("B: %+v", B)

	log.Printf("-------------------")
	log.Printf("C: %+v", C)

	log.Printf("-------------------")

	var params []float64
	for i := 0; i < len(C); i++ {
		params = append(params, float64(i))
		params = append(params, C[i])
	}

	xy := goNum.NewMatrix(len(params)/2, 2, params)

	if out, _, _, ok := goNum.FittingPolynomial(xy, 2); ok == true {
		//log.Printf("out: %+v", out)
		log.Printf("C out: %f, %f, %f", out.Data[2], out.Data[1], out.Data[0])
	} else {
		log.Printf("err is %v", ok)
	}
}
