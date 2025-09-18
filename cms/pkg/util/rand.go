package util

import (
	"math/rand"
	"sync"
	"time"
)

//随机种子（随机数底层不是线程安全，必须手动加锁）
var rndSeed = rand.NewSource(time.Now().UnixNano())
var randomer = rand.New(rndSeed)
var randMutex = sync.Mutex{}

//生成范围内的随机整数pseudo-random number in [min,max)
func RndInt(min int, max int) int {
	randMutex.Lock()
	defer randMutex.Unlock()

	v := randomer.Intn(max - min)

	return min + v
}

//生成0,1范围内的随机小数
func RndFloat() float32 {
	randMutex.Lock()
	defer randMutex.Unlock()

	return randomer.Float32()
}

func RndNInt(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}
	//存放结果的slice
	nums := make([]int, 0)

	for len(nums) < count {
		//生成随机数
		num := RndInt(start, end)
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}
