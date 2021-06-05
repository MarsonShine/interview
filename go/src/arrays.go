package interview

func FirstMissingPositiveForce(nums []int) int {
	var i int = 1
	for i < 2<<31 {
		b := true
		for j := 0; j < len(nums); j++ {
			if i == nums[j] {
				b = true
				break
			} else {
				b = false
			}
		}
		if !b {
			return i
		}
		i += 1
	}
	return i
}

func FirstMissingPositive(nums []int) int {
	n := len(nums)
	for i := 0; i < n; i++ {
		for nums[i] > 0 && nums[i] <= n && nums[nums[i]-1] != nums[i] {
			nums[nums[i]-1], nums[i] = nums[i], nums[nums[i]-1]
		}
	}
	for i := 0; i < n; i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}
	return n + 1
}
