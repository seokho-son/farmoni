// Proof of Concepts for the Cloud-Barista Multi-Cloud Project.
//      * Cloud-Barista: https://github.com/cloud-barista
//
// Short reports of local resource monitoring info.
//
// by powerkim@powerkim.co.kr, 2019.02.
 package main


 import (
         "os"
         "farmoni/localmoni/cpu_usage"
         "farmoni/localmoni/mem_usage"
         "farmoni/localmoni/disk_stat"
         "strconv"
         "time"
	 "github.com/dustin/go-humanize"
 )

 func cpu() {
        // utilization for each logical CPU
        strCPUUtilizationArr := cpuusage.GetAllUtilPercentages()
	print("  [CPU USG]")
        for i, cpupercent := range strCPUUtilizationArr {
		if(i!=0) { print(", ") }
                print(" C" + strconv.Itoa(i) +":" + cpupercent + "%")
        }
 }

 func mem() {
	// total memory in this machine
	totalMem := memusage.GetTotalMem()
	// mega byte
	//strTotalMemM := strconv.FormatUint(totalMem/1024/1024, 10)
	strTotalMemM := humanize.Comma(int64(totalMem/1024/1024))

	// used memory in this machine
	usedMem := memusage.GetUsedMem()
	// mega byte
	//strUsedMemM := strconv.FormatUint(usedMem/1024/1024, 10)
	strUsedMemM := humanize.Comma(int64(usedMem/1024/1024))

	// free memory in this machine
	freeMem := memusage.GetFreeMem()
	// mega byte
	//strFreeMemM := strconv.FormatUint(freeMem/1024/1024, 10)
	strFreeMemM := humanize.Comma(int64(freeMem/1024/1024))

	println("  [MEM USG] TOTAL: " + strTotalMemM + "MB, USED: " + strUsedMemM + "MB, FREE: " + strFreeMemM + "MB") 
 }


 func main() {
	// get Host Name
	hostname, _ := os.Hostname()

	var readBytes [] uint64 = make([]uint64, 1)
	var writeBytes [] uint64 = make([]uint64, 1)
	var beforeReadBytes [] uint64 = make([]uint64, 1)
	var beforeWriteBytes [] uint64 = make([]uint64, 1)

	for{
		println("[" + hostname + "]")
		cpu()
		println("")

		mem()

		print("  [DSK RAT]")
		// get effective partion list
		partitionList := diskstat.GetPartitionList()
		for i, partition := range partitionList {
			print(partition + ": ")
			if(len(readBytes)<(i+1)) {
				rBytes, wBytes := diskstat.GetRWBytes(partition)
				readBytes = append(readBytes, rBytes)
				writeBytes = append(writeBytes, wBytes)
			}else{
				readBytes[i], writeBytes[i] = diskstat.GetRWBytes(partition)
			}
			print("R/s:   " + strconv.FormatUint(readBytes[i]-beforeReadBytes[i], 10))
			println(", W/s:   " + strconv.FormatUint(writeBytes[i]-beforeWriteBytes[i], 10))

			if(len(readBytes)<(i+1)) {
				beforeReadBytes = append(beforeReadBytes, readBytes[i])
				beforeWriteBytes = append(beforeWriteBytes, writeBytes[i])
			}else {
				beforeReadBytes[i] = readBytes[i]
				beforeWriteBytes[i] = writeBytes[i]
			}
		}
		println("-----------")
		time.Sleep(time.Second)	
	}

 }
