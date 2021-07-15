package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
*/
import "C"

func COptions(options []string) **C.char {
	length := len(options)
	opts := make([]*C.char, length+1)
	for i := 0; i < length; i++ {
		opts[i] = C.CString(options[i])
	}
	opts[length] = (*C.char)(nil)

	return (**C.char)(&opts[0])
}
