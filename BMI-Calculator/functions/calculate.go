package function

func CalculateBmi(weight float32, height float32) float32 {
	bmiResult := weight / (height*height)
	return bmiResult
}
