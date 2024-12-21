export interface Hwmon {
	name: string
	composite: number
	sensors: HwmonSensor[]
}

interface HwmonSensor {
	name: string
	temp: number
}

// // not being used
// export interface VmStat {
//     procs_r: number
//     procs_b: number
//     memory_swpd: number
//     memory_free: number
//     memory_buff: number
//     memory_cache: number
//     swap_si: number
//     swap_so: number
//     io_bi: number
//     io_bo: number
//     system_in: number
//     system_cs: number
//     cpu_us: number
//     cpu_sy: number
//     cpu_id: number
//     cpu_wa: number
//     cpu_st: number
//     cpu_gu: number
//     timestamp: Date | string
// }
