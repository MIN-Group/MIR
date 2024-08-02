// Copyright [2022] [MIN-Group -- Peking University Shenzhen Graduate School Multi-Identifier Network Development Group]
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package table
// @Author: yzy
// @Description:
// @Version: 1.0.0
// @Date: 2021/12/21 4:32 PM
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//

package table

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"minlib/component"
	utils2 "minlib/utils"
	"mir-go/daemon/lf"
	"mir-go/daemon/utils"
	"os"
	"runtime"
	"testing"
)

//const datasetPath = "/home/gdcni22/linlh/prefix_dataset/prefix_dataset/dmoz_subdomain_min_name.txt"
const datasetPath = "../../../dmoz_subdomain_min_name.txt"

var insertNum = flag.Int("insertNum", 2599593, "number of insert name")          // 默认为dmoz_subdomain_min_name的数据集大小2599593
var parallelNum = flag.Int("parallelNum", 10, "number of goroutine to test lpm") // 默认为开启10个协程进行并发测试
var minNameIdentifierCache []*component.Identifier                               // 缓存创建的所有identifier
var searchIdentifiers []*component.Identifier                                    // 缓存用来查询的 Identifier

// 预加载测试数据集，生成identifier
// 测试使用的数据集请到：https://gitee.com/michael_llh/prefix_dataset 进行下载
func LoadDatasetToIdentifier() {
	//func init() {
	file, err := os.Open(datasetPath)
	if err != nil {
		log.Fatal("fail to open benchmarking dataset file", err.Error())
	}

	buf := bufio.NewReader(file)
	// 将文件指针重置到开头
	file.Seek(0, io.SeekStart)
	buf.Reset(file)
	minNameIdentifierCache = make([]*component.Identifier, *insertNum)
	searchIdentifiers = make([]*component.Identifier, *insertNum)
	for i := 0; i < *insertNum; i++ {
		minNameBytes, err := buf.ReadBytes('\n')
		if err != nil {
			log.Fatal("read line from dataset failed", err.Error())
		}
		identifier, err := component.CreateIdentifierByString(string(minNameBytes[0 : len(minNameBytes)-1]))
		if err != nil {
			log.Fatalf(err.Error())
		}
		// 缓存 ToUri 和 GetPrefix
		identifier.ToUri()
		minNameIdentifierCache[i] = identifier
		// 由于是最长前缀匹配，所以这里需要将名字随机的进行增加，为了测试的公平性，随机种子统一设置为行号
		appendName := utils.RandomMINName(1, 5, 10, int64(i+1)) // 随机生成commponet长度为[1,10]，component数量为[1,5]的随机名字
		originIdentifierStr, err := utils2.Unescape(identifier.ToUri())
		if err != nil {
			log.Fatalf(err.Error())
		}
		minName := originIdentifierStr + appendName
		queryIdentifier, err := component.CreateIdentifierByString(minName)
		if err != nil {
			log.Fatalf(err.Error())
		}
		// 缓存 ToUri 和 GetPrefix
		queryIdentifier.ToUri()
		searchIdentifiers[i] = queryIdentifier
	}
	file.Close()
	//    fmt.Println("loading dataset done, total cache:", *insertNum, " identifiers")
}

// 单元测试
func TestFIBFindLongestPrefixMatch(t *testing.T) {
	// 测试精确匹配
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	// 打印结果 &{LogicFaceId:0xc000056140 map[1:{LogicFaceId:1 1}] 0xc00001a258}
	fmt.Println(fib.FindLongestPrefixMatch(identifier))

	// 测试最长前缀匹配 存在
	identifier, err = component.CreateIdentifierByString("/min/pku/edu/cn")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 &{LogicFaceId:0xc000056140 map[1:{LogicFaceId:1 1}] 0xc00001a258}
	fmt.Println(fib.FindLongestPrefixMatch(identifier))

	// 测试最长前缀匹配 不存在
	identifier, err = component.CreateIdentifierByString("/min/pku")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 <nil>
	fmt.Println(fib.FindLongestPrefixMatch(identifier))

	// 测试异常情况 如标识没有初始化
	// 打印结果 <nil>
	fmt.Println(fib.FindLongestPrefixMatch(&component.Identifier{}))

	// 测试异常情况 加入的标识未初始化
	fib.AddOrUpdate(&component.Identifier{}, &lf.LogicFace{LogicFaceId: 1}, 1)
}

func TestFIBFindExactMatch(t *testing.T) {
	// 测试精确匹配 /min/pku/edu
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	// 打印结果 &{LogicFaceId:0xc0000d0120 map[1:{LogicFaceId:1 1}] 0xc0000ec060}
	fmt.Println(fib.FindExactMatch(identifier))

	// 测试精确匹配 /min/pku/edu/cn
	identifier, err = component.CreateIdentifierByString("/min/pku/edu/cn")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 nil
	fmt.Println(fib.FindExactMatch(identifier))

	// 测试精确匹配 /min/pku
	identifier, err = component.CreateIdentifierByString("/min/pku")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 nil
	fmt.Println(fib.FindExactMatch(identifier))

	// 测试精确匹配 /min/pku2
	identifier, err = component.CreateIdentifierByString("/min/pku2")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 nil
	fmt.Println(fib.FindExactMatch(identifier))

}

func TestFIBAddOrUpdate(t *testing.T) {
	fib := CreateFIB()
	// 测试异常情况 加入的标识未初始化
	fibEntry := fib.AddOrUpdate(&component.Identifier{}, &lf.LogicFace{LogicFaceId: 1}, 1)
	// 打印结果 &{{} {<nil>} []}
	fmt.Println(fibEntry.GetIdentifier())

	// 测试add
	identifier, err := component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1))

	// 测试update
	fmt.Println(fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 0))
}

func TestFIBEraseByIdentifier(t *testing.T) {
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	// 测试删除成功
	// 打印结果 <nil>
	fmt.Println(fib.EraseByIdentifier(identifier))
	// 测试删除失败
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	identifier, err = component.CreateIdentifierByString("/min/pku/edu/cn")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 NodeError: the entry is not existed
	fmt.Println(fib.EraseByIdentifier(identifier))
	identifier, err = component.CreateIdentifierByString("/min/pku")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 NodeError: the entry is not existed
	fmt.Println(fib.EraseByIdentifier(identifier))
	identifier, err = component.CreateIdentifierByString("/min/pku2")
	if err != nil {
		fmt.Println(err)
	}
	// 打印结果 NodeError: the entry is not existed
	fmt.Println(fib.EraseByIdentifier(identifier))
}

func TestFIBEraseByFIBEntry(t *testing.T) {
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fibEntry := fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	// 测试删除成功
	// 打印结果 <nil>
	fmt.Println(fib.EraseByFIBEntry(fibEntry))
}

func TestFIBRemoveNextHopByFace(t *testing.T) {
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	identifier, err = component.CreateIdentifierByString("/min/pku/cn")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)
	identifier, err = component.CreateIdentifierByString("/min/gdcni15/mir-go")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	identifier, err = component.CreateIdentifierByString("/min/gdcni15/filegator")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)
	// 打印结果 2 0 2 0
	fmt.Println(fib.RemoveNextHopByFace(&lf.LogicFace{LogicFaceId: 0}))
	fmt.Println(fib.RemoveNextHopByFace(&lf.LogicFace{LogicFaceId: 0}))
	fmt.Println(fib.RemoveNextHopByFace(&lf.LogicFace{LogicFaceId: 1}))
	fmt.Println(fib.RemoveNextHopByFace(&lf.LogicFace{LogicFaceId: 1}))
}

func TestFIBSize(t *testing.T) {
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	fmt.Println(fib.Size())
	identifier, err = component.CreateIdentifierByString("/min/pku")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	fmt.Println(fib.Size())
	identifier, err = component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	fmt.Println(fib.Size())

	identifier, err = component.CreateIdentifierByString("/min/pku/cn")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)
	fmt.Println(fib.Size())

	identifier, err = component.CreateIdentifierByString("/min/gdcni15/mir-go")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
	fmt.Println(fib.Size())

	identifier, err = component.CreateIdentifierByString("/min/gdcni15/filegator")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)
	fmt.Println(fib.Size())

}

func TestGetDepth(t *testing.T) {
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/pku")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/pku/cn")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)

	identifier, err = component.CreateIdentifierByString("/min/gdcni15/mir-go")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/gdcni15/filegator")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)

	identifier, err = component.CreateIdentifierByString("/min/gdcni15/filegator/test/test2")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)
	fmt.Println(fib.GetDepth())
}

// 基准测试
// allocs/op表示每个op(单次迭代)发生了多少个不同的内存分配.
// B/op是每操作分配多少个字节.
// 测试命令：go test -benchmem -run=^$ -bench ^BenchmarkFIBAddOrUpdate$ -benchtime=2599593x
func BenchmarkFIBAddOrUpdate(b *testing.B) {
	fib := CreateFIB()

	lface := lf.LogicFace{LogicFaceId: 1}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fibEntry := fib.AddOrUpdate(minNameIdentifierCache[i], &lface, 0)
		if fibEntry == nil {
			b.Fatal("add identifier failed")
		}
	}
}

// go test -benchmem -run=^$ -bench ^BenchmarkFIBEraseByIdentifier$ -benchtime=259960x
func BenchmarkFIBEraseByIdentifier(b *testing.B) {
	fib := CreateFIB()
	for i := 0; i < b.N; i++ {
		identifier := minNameIdentifierCache[i]
		lface := lf.LogicFace{LogicFaceId: 1}
		fibEntry := fib.AddOrUpdate(identifier, &lface, 0)
		if fibEntry == nil {
			b.Fatal("add identifier failed")
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := fib.EraseByIdentifier(minNameIdentifierCache[i])
		if err != nil {
			b.Fatal("fib erase by identifier failed", err.Error())
		}
	}
}

// go test -benchmem -run=^$ -bench ^BenchmarkFIBEraseByFIBEntry$ -benchtime=259960x
func BenchmarkFIBEraseByFIBEntry(b *testing.B) {
	fib := CreateFIB()

	for i := 0; i < b.N; i++ {
		lface := lf.LogicFace{LogicFaceId: 1}
		fibEntry := fib.AddOrUpdate(minNameIdentifierCache[i], &lface, 0)
		if fibEntry == nil {
			b.Fatal("add identifier failed")
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		identifier := minNameIdentifierCache[i]
		fibEntry := FIBEntry{identifier: identifier}
		b.StartTimer()
		err := fib.EraseByFIBEntry(&fibEntry)
		if err != nil {
			b.Fatal("fib erase by fib entry failed", err.Error())
		}
	}
}

func BenchmarkFIBFindExactMatch(b *testing.B) {
	fib := CreateFIB()

	for i := 0; i < b.N; i++ {
		lface := lf.LogicFace{LogicFaceId: 1}
		fibEntry := fib.AddOrUpdate(minNameIdentifierCache[i], &lface, 0)
		if fibEntry == nil {
			b.Fatal("add identifier failed")
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fibEntry := fib.FindExactMatch(minNameIdentifierCache[i])
		if fibEntry == nil {
			b.Fatal("fib find exact match", minNameIdentifierCache[i].ToUri())
		}
	}
}

func BenchmarkFIBFindLongestPrefixMatch(b *testing.B) {
	fib := CreateFIB()

	for i := 0; i < b.N; i++ {
		identifier := minNameIdentifierCache[i]
		lface := lf.LogicFace{LogicFaceId: 1}
		fibEntry := fib.AddOrUpdate(identifier, &lface, 0)
		if fibEntry == nil {
			b.Fatal("add identifier failed")
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fibEntry := fib.FindLongestPrefixMatch(searchIdentifiers[i])
		if fibEntry == nil {
			b.Fatal("fib find longest prefix match", searchIdentifiers[i].ToUri())
		}
	}
}

type TestItem struct {
	fib       *FIB
	queryList []*component.Identifier
	cnt       int
	lastOne   bool
}

// b.StopTimer() 消除add函数添加的额外时间 测试时间60s左右 因为启用了定时器
func BenchmarkFIBRemoveNextHopByFace(b *testing.B) {
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)
		b.StartTimer()
		fib.RemoveNextHopByFace(&lf.LogicFace{LogicFaceId: 1})
	}
}

func BenchmarkFIBSize(b *testing.B) {
	fib := CreateFIB()
	identifier, err := component.CreateIdentifierByString("/min")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/pku")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/pku/edu")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/pku/cn")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)

	identifier, err = component.CreateIdentifierByString("/min/gdcni15/mir-go")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 1}, 1)

	identifier, err = component.CreateIdentifierByString("/min/gdcni15/filegator")
	if err != nil {
		fmt.Println(err)
	}
	fib.AddOrUpdate(identifier, &lf.LogicFace{LogicFaceId: 0}, 1)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fib.Size()
	}
}

func BenchmarkFIBFindLongestPrefixMatch_Parallel(b *testing.B) {
	numOfGoroutine := *parallelNum // 获取命令行输入的并行数量，对应的就是创建的goroutine并行数量
	testItemChan := make(chan *TestItem, numOfGoroutine)
	testItemList := make([]TestItem, numOfGoroutine)
	for i := 0; i < numOfGoroutine; i++ {
		testItemList[i].queryList = make([]*component.Identifier, b.N/numOfGoroutine+1)
		testItemList[i].cnt = 0
		if i == numOfGoroutine-1 {
			testItemList[i].lastOne = true
		} else {
			testItemList[i].lastOne = false
		}
	}
	fib := CreateFIB()

	for i := 0; i < b.N; i++ {
		identifier := minNameIdentifierCache[i]
		lface := lf.LogicFace{LogicFaceId: 1}
		fibEntry := fib.AddOrUpdate(identifier, &lface, 0)
		if fibEntry == nil {
			b.Fatal("add identifier failed")
		}
		testItemList[i%numOfGoroutine].queryList[testItemList[i%numOfGoroutine].cnt] = searchIdentifiers[i]
		testItemList[i%numOfGoroutine].cnt++
	}

	// 分发到各个协程的channel
	for n := 0; n < numOfGoroutine; n++ {
		testItemList[n].fib = fib
		testItemChan <- &testItemList[n]
	}

	//      b.SetParallelism(numOfGoroutine) // 使用这个函数设置最大的并发数，不知道为什么channel会卡住
	runtime.GOMAXPROCS(numOfGoroutine) // 通过设置最大线程数来控制RunParallel启动的携程数量
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// 取出协程要完成的工作
		testItem := <-testItemChan
		var idx int = 0
		fib := testItem.fib
		cnt := testItem.cnt
		for idx < cnt && pb.Next() {
			fibEntry := fib.FindLongestPrefixMatch(searchIdentifiers[idx])
			if fibEntry == nil {
				b.Fatal("fib find longest prefix match", searchIdentifiers[idx].ToUri())
			}
			idx++
		}
		if testItem.lastOne {
			pb.Next()
		}
	})
}

func BenchmarkFIBFindExactMatch_Parallel(b *testing.B) {
	numOfGoroutine := *parallelNum // 获取命令行输入的并行数量，对应的就是创建的goroutine并行数量
	testItemChan := make(chan *TestItem, numOfGoroutine)
	testItemList := make([]TestItem, numOfGoroutine)
	for i := 0; i < numOfGoroutine; i++ {
		testItemList[i].queryList = make([]*component.Identifier, b.N/numOfGoroutine+1)
		testItemList[i].cnt = 0
		if i == numOfGoroutine-1 {
			testItemList[i].lastOne = true
		} else {
			testItemList[i].lastOne = false
		}
	}
	fib := CreateFIB()

	for i := 0; i < b.N; i++ {
		identifier := minNameIdentifierCache[i]
		lface := lf.LogicFace{LogicFaceId: 1}
		fibEntry := fib.AddOrUpdate(identifier, &lface, 0)
		if fibEntry == nil {
			b.Fatal("add identifier failed")
		}
		testItemList[i%numOfGoroutine].queryList[testItemList[i%numOfGoroutine].cnt] = minNameIdentifierCache[i]
		testItemList[i%numOfGoroutine].cnt++
	}

	// 分发到各个协程的channel
	for n := 0; n < numOfGoroutine; n++ {
		testItemList[n].fib = fib
		testItemChan <- &testItemList[n]
	}

	//      b.SetParallelism(numOfGoroutine) // 使用这个函数设置最大的并发数，不知道为什么channel会卡住
	runtime.GOMAXPROCS(numOfGoroutine) // 通过设置最大线程数来控制RunParallel启动的携程数量
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// 取出协程要完成的工作
		testItem := <-testItemChan
		var idx int = 0
		fib := testItem.fib
		cnt := testItem.cnt
		for idx < cnt && pb.Next() {
			fibEntry := fib.FindExactMatch(minNameIdentifierCache[idx])
			if fibEntry == nil {
				b.Fatal("fib find longest prefix match", minNameIdentifierCache[idx].ToUri())
			}
			idx++
		}
		if testItem.lastOne {
			pb.Next()
		}
	})
}

func TestMain(m *testing.M) {
	flag.Parse() // 使用到了insertNum来控制每次加载的数据数量，所以这里需要提前调用解析命令行参数
	LoadDatasetToIdentifier()
	m.Run()
}
