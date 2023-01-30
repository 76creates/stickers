package flexbox

import (
	"math"
)

// TODO: explain this mess using some comments

func calculateRatioWithMinimum(distribute int, matrix []int, minimumMatrix []int) (ratioDistribution []int) {
	for range matrix {
		ratioDistribution = append(ratioDistribution, 0)
	}

	dist := calculateRatio(distribute, matrix)
	for i, d := range dist {
		if minimumMatrix[i] > d {
			ratioDistribution[i] = minimumMatrix[i]

			distribute -= minimumMatrix[i]
			matrix[i] = 0
			minimumMatrix[i] = 0
			_dist := calculateRatioWithMinimum(distribute, matrix, minimumMatrix)
			for ii, _d := range _dist {
				if ii != i {
					ratioDistribution[ii] = _d
				}
			}
			break
		} else {
			ratioDistribution[i] = d
		}
	}
	// TODO: calculate reminder and if negative shrink right most column

	return ratioDistribution
}

func calculateRatio(distribute int, matrix []int) (ratioDistribution []int) {
	if distribute == 0 {
		for range matrix {
			ratioDistribution = append(ratioDistribution, 0)
		}
		return ratioDistribution
	}

	var combinedRatios int
	for _, value := range matrix {
		combinedRatios += value
	}

	if combinedRatios > 0 {
		var reminder int
		ratioDistribution, reminder = distributeToMatrix(distribute, combinedRatios, matrix)
		if reminder > 0 {
			for index, reminderAdded := range distributeReminder(reminder, matrix) {
				ratioDistribution[index] += reminderAdded
				reminder -= reminderAdded
			}
		}
		// TODO: rethink maybe, does this fn belong here
		if reminder < 0 {
			// happens when there is minimum value
		}
	}

	return ratioDistribution
}

func distributeToMatrix(distribute int, combinedRatio int, matrix []int) (distribution []int, reminder int) {
	reminder = distribute
	for _, max := range matrix {
		ratioDistributionValue := int(math.Floor((float64(max) / float64(combinedRatio)) * float64(distribute)))
		distribution = append(distribution, ratioDistributionValue)
		reminder -= ratioDistributionValue

	}
	return distribution, reminder
}

func calculateMatrixRatio(distribute int, matrix [][]int) (ratioDistribution []int) {
	// get matrix max ratio for each int in matrix slice
	var maxRatio []int
	for matrixIndex, ratios := range matrix {
		maxRatio = append(maxRatio, 0)
		for _, ratio := range ratios {
			if ratio > maxRatio[matrixIndex] {
				maxRatio[matrixIndex] = ratio
			}
		}
	}

	return calculateRatio(distribute, maxRatio)
}

// distributeReminder is simple reminder distributor, it will distribute add 1 to next highest
// matrix value till it runs out of reminder to distribute, this might be improved for some more
// complex cases
func distributeReminder(reminder int, matrixMaxRatio []int) (reminderDistribution []int) {
	for range matrixMaxRatio {
		reminderDistribution = append(reminderDistribution, 0)
	}

	distributed := 0
	for reminder > 0 {
		maxIndex := 0
		maxRatio := 0
		for index, ratio := range matrixMaxRatio {
			// skip if already expanded
			if reminderDistribution[index] > 0 {
				continue
			}
			if ratio > maxRatio {
				maxRatio = ratio
				maxIndex = index
			}
		}
		reminderDistribution[maxIndex] += 1
		distributed += 1
		reminder -= 1
	}

	return reminderDistribution
}
