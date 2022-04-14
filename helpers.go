package main

import (
	"strconv"
)

// Helpers

func getData(trialid string) []SensorData {
	data := []SensorData{}
	err := db.Select(&data, "SELECT * FROM data WHERE trialid=$1", trialid)
	if err != nil {
		return []SensorData{}
	}
	return data
}

func getTrials() []Trial {
	trials := []Trial{}
	err := db.Select(&trials, "SELECT id, start_time FROM trials")
	if err != nil {
		return []Trial{}
	}
	return trials
}

func getTrialsComplete() []Trial {
	trials := getTrials()

	for i, t := range trials {
		trials[i].SensorData = getData(strconv.Itoa(t.ID))
	}

	return trials
}

// func getRows(trialid string) int {
// 	if trialid == "-1" {
// 		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
// 		if err != nil {
// 			return -1
// 		}
// 	}
// 	count := 0
// 	if trialid != "0" {
// 		err := db.Get(&count, "SELECT COUNT(*) FROM data", trialid)

// 		if err != nil {
// 			return -1
// 		}
// 	} else {
// 		err := db.Get(&count, "SELECT COUNT(*) FROM data WHERE trialid=$1", trialid)

// 		if err != nil {
// 			return -1
// 		}
// 	}

// 	return count
// }

func getThreshold(trialid string) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}
	threshold := 0
	max := 0
	for i := 100; i < 500; i++ {
		steps := 0

		if trialid != "0" {
			err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data WHERE trialid=$2) AS data WHERE steps", i, trialid)

			if err != nil {
				return -1
			}
		} else {
			err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data) AS data WHERE steps", i)

			if err != nil {
				return -1
			}
		}

		if steps > max {
			max = steps
			threshold = i
		}
	}
	return threshold
}

func getPronation(trialid string, pThreshold int) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}

	threshold := getThreshold(trialid)

	pronation := 0

	if trialid != "0" {
		err := db.Get(&pronation, "SELECT COUNT(pronation) FROM (SELECT ((((a2 + a3) - (a1 + a5)) > $1) AND ((a1 + a2 + a3 + a5) > $2)) AS pronation FROM data WHERE trialid=$3) AS data WHERE pronation", strconv.Itoa(pThreshold), strconv.Itoa(threshold), trialid)
		if err != nil {
			return -1
		}
	} else {
		err := db.Get(&pronation, "SELECT COUNT(pronation) FROM (SELECT ((((a2 + a3) - (a1 + a5)) > $1) AND ((a1 + a2 + a3 + a5) > $2)) AS pronation FROM data) AS data WHERE pronation", strconv.Itoa(pThreshold), strconv.Itoa(threshold))
		if err != nil {
			return -1
		}
	}

	return pronation
}

func getSupination(trialid string, sThreshold int) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}

	threshold := getThreshold(trialid)

	supination := 0

	if trialid != "0" {
		err := db.Get(&supination, "SELECT COUNT(supination) FROM (SELECT ((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $2)) AS supination FROM data WHERE trialid=$3) AS data WHERE supination", strconv.Itoa(sThreshold), strconv.Itoa(threshold), trialid)
		if err != nil {
			return -1
		}
	} else {
		err := db.Get(&supination, "SELECT COUNT(supination) FROM (SELECT ((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $2)) AS supination FROM data) AS data WHERE supination", strconv.Itoa(sThreshold), strconv.Itoa(threshold))
		if err != nil {
			return -1
		}
	}

	return supination
}

func getNormal(trialid string, pThreshold int, sThreshold int) int {
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return -1
		}
	}

	threshold := getThreshold(trialid)

	normal := 0

	if trialid != "0" {
		err := db.Get(&normal, "SELECT COUNT(normal) FROM (SELECT ((NOT((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $3)) AND NOT((((a2 + a3) - (a1 + a5)) > $2) AND ((a1 + a2 + a3 + a5) > $3))) AND ((a1 + a2 + a3 + a5) > $3)) AS normal FROM data WHERE trialid=$4) AS data WHERE normal", strconv.Itoa(sThreshold), strconv.Itoa(pThreshold), strconv.Itoa(threshold), trialid)
		if err != nil {
			return -1
		}
	} else {
		err := db.Get(&normal, "SELECT COUNT(normal) FROM (SELECT ((NOT((((a2 + a3) - (a1 + a5)) < $1) AND ((a1 + a2 + a3 + a5) > $3)) AND NOT((((a2 + a3) - (a1 + a5)) > $2) AND ((a1 + a2 + a3 + a5) > $3))) AND ((a1 + a2 + a3 + a5) > $3)) AS normal FROM data) AS data WHERE normal", strconv.Itoa(sThreshold), strconv.Itoa(pThreshold), strconv.Itoa(threshold))
		if err != nil {
			return -1
		}
	}

	return normal
}

func getSteps(trialid string) string {
	steps := 0
	if trialid == "-1" {
		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
		if err != nil {
			return "-1"
		}
	}
	if trialid == "0" {
		dataC := getTrialsComplete()
		total := 0
		for _, temp := range dataC {
			data := temp.SensorData
			max := 0
			for j := 1; j < 1000; j++ {
				steps = 0

				for i, v := range data {
					if i != len(data)-1 {
						sum := v.A1 + v.A2 + v.A3 + v.A5
						sumN := data[i+1].A1 + data[i+1].A2 + data[i+1].A3 + data[i+1].A5
						if (sum > j) && (sumN < j) {
							steps += 2
						}
					}
				}

				if steps > max {
					max = steps
				}
			}
			total += max
		}
		return strconv.Itoa(total)
	}
	data := getData(trialid)
	max := 0
	for j := 100; j < 500; j++ {
		steps = 0
		if len(data) == 0 {
			return "-1"
		}

		for i, v := range data {
			if i != len(data)-1 {
				sum := v.A1 + v.A2 + v.A3 + v.A5
				sumN := data[i+1].A1 + data[i+1].A2 + data[i+1].A3 + data[i+1].A5
				if (sum > j) && (sumN < j) {
					steps += 2
				}
			}
		}

		if steps > max {
			max = steps
		}
	}

	return strconv.Itoa(max)
}

// func getSteps(trialid string, threshold int) string {
// 	if trialid == "-2" {
// 		steps := getStepsO("-1", -1)
// 		return steps
// 	}
// 	if trialid == "-1" {
// 		err := db.Get(&trialid, "SELECT max(trialid) FROM data")
// 		if err != nil {
// 			return "-1"
// 		}
// 	}
// 	if trialid == "0" && threshold != -1 {

// 		steps := 0
// 		err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY $1) < $1)) AS steps FROM data) AS data WHERE steps", threshold)

// 		if err != nil {
// 			return "-1"
// 		}
// 		return strconv.Itoa(steps)
// 	}
// 	if threshold == -1 {
// 		max := 0
// 		for i := 100; i < 500; i++ {
// 			steps := 0

// 			if trialid != "0" {
// 				err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data WHERE trialid=$2) AS data WHERE steps", i, trialid)

// 				if err != nil {
// 					return "-1"
// 				}
// 			} else {
// 				err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data) AS data WHERE steps", i)

// 				if err != nil {
// 					return "-1"
// 				}
// 			}

// 			if steps > max {
// 				max = steps
// 			}
// 		}
// 		return strconv.Itoa(max)
// 	}
// 	steps := 0
// 	err := db.Get(&steps, "SELECT COUNT(steps)*2 FROM (SELECT ((a1 + a2 + a3 + a5 > $1) AND (LEAD(a1 + a2 + a3 + a5,1) OVER (ORDER BY trialid) < $1)) AS steps FROM data WHERE trialid=$2) AS data WHERE steps", threshold, trialid)

// 	if err != nil {
// 		return "-1"
// 	}

// 	return strconv.Itoa(steps)
// }
