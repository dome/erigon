package estimate

import (
	"runtime"

	"github.com/c2h5oh/datasize"

	"github.com/ledgerwatch/erigon-lib/common/cmp"
	"github.com/ledgerwatch/erigon-lib/mmap"
)

type EstimatedRamPerWorker datasize.ByteSize

// Workers - return max workers amount based on total Memory/CPU's and estimated RAM per worker
func (r EstimatedRamPerWorker) Workers() int {
	// 50% of TotalMemory. Better don't count on 100% because OOM Killer may have aggressive defaults and other software may need RAM
	maxWorkersForGivenMemory := (mmap.TotalMemory() / 2) / uint64(r)
	res := cmp.Min(AlmostAllCPUs(), int(maxWorkersForGivenMemory))
	return cmp.Max(1, res) // must have at-least 1 worker
}

func (r EstimatedRamPerWorker) WorkersHalf() int    { return cmp.Max(1, r.Workers()/2) }
func (r EstimatedRamPerWorker) WorkersQuarter() int { return cmp.Max(1, r.Workers()/4) }

const (
	//elias-fano index building is single-threaded
	// when set it to 3GB - observed OOM-kil at server with 128Gb ram and 32CPU
	IndexSnapshot = EstimatedRamPerWorker(4 * datasize.GB)

	//1-file-compression is multi-threaded
	CompressSnapshot = EstimatedRamPerWorker(1 * datasize.GB)

	//state-reconstitution is multi-threaded
	ReconstituteState = EstimatedRamPerWorker(512 * datasize.MB)
)

// AlmostAllCPUs - return all-but-one cpus. Leaving 1 cpu for "work producer", also cloud-providers do recommend leave 1 CPU for their IO software
// user can reduce GOMAXPROCS env variable
func AlmostAllCPUs() int {
	return cmp.Max(1, runtime.GOMAXPROCS(-1)-1)
}
