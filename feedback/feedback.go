package feedback

import "log"

func chk(obj interface{}, err error) interface{} {
	if err != nil {
		log.Fatal("panic in mail: " + err.Error())
		panic(err.Error())
	}
	return obj
}
