package error

func Check(err error){
	if err != nil {
		panic(err.Error())
	}
}