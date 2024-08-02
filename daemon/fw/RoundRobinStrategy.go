// Package fw
// @Author: Jianming Que (阙建明)
// @Description:
// @Version: 1.0.0
// @Date: 2022/4/1 15:36
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package fw

import (
	common2 "minlib/common"
	"minlib/packet"
	"minlib/utils"
	"mir-go/daemon/lf"
	utils2 "mir-go/daemon/utils"
	"sync/atomic"
	"time"
)

// RoundRobinStrategy 轮询策略实现
//
// @Description:
// 1. 拉式采用最佳路由；
// 2. 推式采用轮询转发，每个可用链路用n秒
//
type RoundRobinStrategy struct {
	BestRouteStrategy
	roundTime    int64
	startTime    int64
	changeTicker *time.Ticker
	currentCount uint64
}

// NewRoundRobinStrategy 新建一个轮询策略
//
// @Description:
// @param roundTime 			轮询时间 => 每个可用链路一次可用的时长，单位为秒
// @return *RoundRobinStrategy
//
func NewRoundRobinStrategy(roundTime int64) *RoundRobinStrategy {
	rrs := new(RoundRobinStrategy)
	rrs.roundTime = roundTime
	rrs.changeTicker = time.NewTicker(time.Duration(roundTime) * time.Second)
	rrs.startTime = utils.GetTimestampMS()
	utils2.GoroutineNoPanic(func() {
		for true {
			<-rrs.changeTicker.C
			atomic.AddUint64(&rrs.currentCount, 1)
		}
	})
	return rrs
}

// AfterReceiveGPPkt
//
// @Description:
// @receiver r
// @param ingress
// @param gPPkt
//
func (r *RoundRobinStrategy) AfterReceiveGPPkt(ingress *lf.LogicFace, gPPkt *packet.GPPkt) {
	fibEntry := r.lookupFibForGPPkt(gPPkt)
	nextHops := fibEntry.GetNextHops()
	selectedHop := nextHops[atomic.LoadUint64(&r.currentCount)%uint64(len(nextHops))]
	if selectedHop == nil {
		// 没有路由无法转发
		common2.LogWarn("No Route")
		return
	}
	r.sendGPPkt(selectedHop.LogicFace, gPPkt)
}
