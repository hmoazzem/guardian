export interface GPU {
	gpu: number
	asic: ASIC
	bus: Bus
	vbios: VBIOS
	driver: Driver 
	board: Board
	ras: RAS
	process_isolation: string
	vram: VRAM
	cache_info: CacheInfo[]
}

export interface GPUMetric {
	gpu: number
	usage: object | string
	power: object
	clock: object
	temperature: Temp
	pcie: object
	ecc: ECC
	ecc_blocks: string
	mem_usage: MemUsage
}

// static
interface CacheSize {
	value: number
	unit: string
}

interface VRAM {
	type: string
	vendor: string
	size: object
	bit_width: number
}

interface ASIC {
	market_name: string
	vendor_id: string
	vendor_name: string
	subvendor_id: string
	device_id: string
	subsystem_id: string
	rev_id: string
	asic_serial: string
	oam_id: number
	num_compute_units: number
	target_graphics_version: string
}

interface Driver {
	name: string
    version: string
}

interface Board {
	model_number: string
	product_serial: string
	fru_id: string
	product_name: string
	manufacturer_name: string
}

interface RAS {
	eeprom_version: string
	parity_schema: string
	single_bit_schema: string
	double_bit_schema: string
	poison_schema: string
	ecc_block_state: string
}

interface VBIOS {
	name: string
	build_date: string
	part_number: string
	version: string
}

interface Bus {
	bdf: string
	max_pcie_width: string | number
	max_pcie_speed: object
	pcie_interface_version: string
	slot_type: string
}

interface CacheSize {
	value: number
	unit: string
}

// CacheInfo represents the details of a cache.
interface CacheInfo {
	cache: string
	cache_properties: string[]
	cache_size: CacheSize
	cache_level: number
	max_num_cu_shared: number
	num_cache_instance: number
}

// metric
interface ValueUnit {
  value: string
  unit: number
}

interface Temp {
	edge: ValueUnit
	hotspot: ValueUnit | string // returns "NA" for iGPUs
	mem: ValueUnit | string
}

interface MemUsage {
	total_vram: ValueUnit
	used_vram: ValueUnit
	free_vram: ValueUnit
	total_visible_vram: ValueUnit
	used_visible_vram: ValueUnit
	free_visible_vram: ValueUnit
	total_gtt: ValueUnit
	used_gtt: ValueUnit
	free_gtt: ValueUnit
}

interface ECC {
	total_correctable_count: number
	total_uncorrectable_count: number
	total_deferred_count: number
	cache_correctable_count: string
	cache_uncorrectable_count: string
}