package amdsmi

/*
#cgo CFLAGS: -I/opt/rocm-6.3.0/include
#cgo LDFLAGS: -L/opt/rocm-6.3.0/lib -lamd_smi
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "amd_smi/amdsmi.h"
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func GetGPUTempCgo() {
	// Initialize AMD SMI
	ret := C.amdsmi_init(C.AMDSMI_INIT_AMD_GPUS)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		fmt.Println("Failed to initialize AMD SMI")
		return
	}
	defer C.amdsmi_shut_down()

	// Get the socket count
	var socketCount C.uint32_t
	ret = C.amdsmi_get_socket_handles(&socketCount, nil)
	if ret != C.AMDSMI_STATUS_SUCCESS {
		fmt.Println("Failed to get socket count")
		return
	}
	fmt.Printf("Total Sockets: %d\n", socketCount)

	// Allocate memory for socket handles
	sockets := make([]C.amdsmi_socket_handle, socketCount)
	ret = C.amdsmi_get_socket_handles(&socketCount, (*C.amdsmi_socket_handle)(unsafe.Pointer(&sockets[0])))
	if ret != C.AMDSMI_STATUS_SUCCESS {
		fmt.Println("Failed to get socket handles")
		return
	}

	// Iterate over each socket
	for i := C.uint32_t(0); i < socketCount; i++ {
		// Get socket info
		var socketInfo [128]C.char
		ret = C.amdsmi_get_socket_info(sockets[i], 128, &socketInfo[0])
		if ret != C.AMDSMI_STATUS_SUCCESS {
			fmt.Println("Failed to get socket info")
			continue
		}
		fmt.Printf("Socket: %s\n", C.GoString(&socketInfo[0]))

		// Get processor count
		var deviceCount C.uint32_t
		ret = C.amdsmi_get_processor_handles(sockets[i], &deviceCount, nil)
		if ret != C.AMDSMI_STATUS_SUCCESS {
			fmt.Println("Failed to get processor count")
			continue
		}

		// Allocate memory for processor handles
		processorHandles := make([]C.amdsmi_processor_handle, deviceCount)
		ret = C.amdsmi_get_processor_handles(sockets[i], &deviceCount, (*C.amdsmi_processor_handle)(unsafe.Pointer(&processorHandles[0])))
		if ret != C.AMDSMI_STATUS_SUCCESS {
			fmt.Println("Failed to get processor handles")
			continue
		}

		// Iterate over each processor
		for j := C.uint32_t(0); j < deviceCount; j++ {
			// Get processor type
			var processorType C.processor_type_t
			ret = C.amdsmi_get_processor_type(processorHandles[j], &processorType)
			if ret != C.AMDSMI_STATUS_SUCCESS || processorType != C.AMDSMI_PROCESSOR_TYPE_AMD_GPU {
				fmt.Println("Failed to verify processor type")
				continue
			}

			// Get GPU board info
			var boardInfo C.amdsmi_board_info_t
			ret = C.amdsmi_get_gpu_board_info(processorHandles[j], &boardInfo)
			if ret != C.AMDSMI_STATUS_SUCCESS {
				fmt.Println("Failed to get board info")
				continue
			}
			fmt.Printf("\tDevice %d\n\t\tName: %s\n", j, C.GoString(&boardInfo.product_name[0]))

			// Get temperature
			var temp C.int64_t
			ret = C.amdsmi_get_temp_metric(processorHandles[j], C.AMDSMI_TEMPERATURE_TYPE_EDGE, C.AMDSMI_TEMP_CURRENT, &temp)
			if ret != C.AMDSMI_STATUS_SUCCESS {
				fmt.Println("Failed to get temperature")
				continue
			}
			fmt.Printf("\t\tTemperature: %d°C\n", temp)
		}
	}
}
