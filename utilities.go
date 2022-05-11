package gdal

/*
#include "go_gdal.h"
#include "gdal_version.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

/* --------------------------------------------- */
/* Command line utility wrapper functions        */
/* --------------------------------------------- */

func stringArrayContains(array []string, needle string) bool {
	for _, s := range array {
		if s == needle {
			return true
		}
	}
	return false
}

func Warp(dstDS string, sourceDS []Dataset, options []string) (Dataset, error) {
	if dstDS == "" {
		dstDS = "MEM:::"
		if !stringArrayContains(options, "-of") {
			options = append([]string{"-of", "MEM"}, options...)
		}
	}

	warpopts := C.GDALWarpAppOptionsNew(COptions(options), nil)
	defer C.GDALWarpAppOptionsFree(warpopts)

	srcDS := make([]C.GDALDatasetH, len(sourceDS))
	for i, ds := range sourceDS {
		srcDS[i] = ds.cval
	}

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))

	ds := C.GDALWarp(cdstDS, nil, C.int(len(sourceDS)), &srcDS[0], warpopts, &cerr)

	if cerr != 0 {
		return Dataset{}, fmt.Errorf("warp failed with code %d", cerr)
	}

	if ds == nil {
		return Dataset{}, fmt.Errorf(C.GoString(C.CPLGetLastErrorMsg()))
	}

	return Dataset{ds}, nil
}

func Translate(dstDS string, sourceDS Dataset, options []string) (Dataset, error) {
	if dstDS == "" {
		dstDS = "MEM:::"
		if !stringArrayContains(options, "-of") {
			options = append([]string{"-of", "MEM"}, options...)
		}
	}

	translateopts := C.GDALTranslateOptionsNew(COptions(options), nil)
	defer C.GDALTranslateOptionsFree(translateopts)

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))

	ds := C.GDALTranslate(cdstDS, sourceDS.cval, translateopts, &cerr)
	if cerr != 0 {
		return Dataset{}, fmt.Errorf("translate failed with code %d", cerr)
	}

	if ds == nil {
		return Dataset{}, fmt.Errorf(C.GoString(C.CPLGetLastErrorMsg()))
	}

	return Dataset{ds}, nil
}

func VectorTranslate(dstDS string, sourceDS []Dataset, options []string) (Dataset, error) {
	if dstDS == "" {
		dstDS = "MEM:::"
		if !stringArrayContains(options, "-f") {
			options = append([]string{"-f", "MEM"}, options...)
		}
	}

	translateopts := C.GDALVectorTranslateOptionsNew(COptions(options), nil)
	defer C.GDALVectorTranslateOptionsFree(translateopts)

	srcDS := make([]C.GDALDatasetH, len(sourceDS))
	for i, ds := range sourceDS {
		srcDS[i] = ds.cval
	}

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))

	ds := C.GDALVectorTranslate(cdstDS, nil, C.int(len(sourceDS)), &srcDS[0], translateopts, &cerr)

	if cerr != 0 {
		return Dataset{}, fmt.Errorf("vector translate failed with code %d", cerr)
	}

	if ds == nil {
		return Dataset{}, fmt.Errorf(C.GoString(C.CPLGetLastErrorMsg()))
	}

	return Dataset{ds}, nil
}

func Rasterize(dstDS string, sourceDS Dataset, options []string) (Dataset, error) {
	if dstDS == "" {
		dstDS = "MEM:::"
		if !stringArrayContains(options, "-f") {
			options = append([]string{"-of", "MEM"}, options...)
		}
	}

	rasterizeopts := C.GDALRasterizeOptionsNew(COptions(options), nil)
	defer C.GDALRasterizeOptionsFree(rasterizeopts)

	var cerr C.int
	cdstDS := C.CString(dstDS)
	defer C.free(unsafe.Pointer(cdstDS))

	ds := C.GDALRasterize(cdstDS, nil, sourceDS.cval, rasterizeopts, &cerr)

	if cerr != 0 {
		return Dataset{}, fmt.Errorf("rasterize failed with code %d", cerr)
	}

	if ds == nil {
		return Dataset{}, fmt.Errorf(C.GoString(C.CPLGetLastErrorMsg()))
	}

	return Dataset{ds}, nil
}

func RasterizeOverwrite(dstDS Dataset, sourceDS Dataset, options []string) (Dataset, error) {
	rasterizeopts := C.GDALRasterizeOptionsNew(COptions(options), nil)
	defer C.GDALRasterizeOptionsFree(rasterizeopts)

	var cerr C.int

	x := C.GDALRasterize(nil, dstDS.cval, sourceDS.cval, rasterizeopts, &cerr)

	if cerr != 0 {
		return Dataset{}, fmt.Errorf("rasterize failed with code %d", cerr)
	}

	if x == nil {
		return Dataset{}, fmt.Errorf(C.GoString(C.CPLGetLastErrorMsg()))
	}

	return dstDS, nil
}
