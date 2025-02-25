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

// Package common
// @Author: Jianming Que
// @Description:
// @Version: 1.0.0
// @Date: 2021/4/2 10:20 上午
// @Copyright: MIN-Group；国家重大科技基础设施——未来网络北大实验室；深圳市信息论与未来网络重点实验室
//
package common

import (
	"minlib/security"

	"gopkg.in/ini.v1"
)

// MIRConfig
// 表示 MIR 配置文件的配置，与 mirconf.ini 中的配置一一对应
//
// @Description:
//
type MIRConfig struct {
	GeneralConfig    `ini:"General"`
	LogConfig        `ini:"Log"`
	TableConfig      `ini:"Table"`
	LogicFaceConfig  `ini:"LogicFace"`
	SecurityConfig   `ini:"Security"`
	ForwarderConfig  `ini:"Forwarder"`
	StrategyConfig   `ini:"StrategyConfig"`
	ManagementConfig `ini:"Management"`
	PcapConfig       `ini:"Pcap"`

	configPath string // 存储配置文件路径
}

// Init
// 初始化配置，给所有的配置设置默认值
//
// @Description:
// @receiver mirConfig
//
func (mirConfig *MIRConfig) Init() {
	// General
	mirConfig.GeneralConfig.DefaultId = "/localhost/mir"
	mirConfig.GeneralConfig.EncryptedPasswdSavePath = "/usr/local/etc/mir/passwd"
	mirConfig.GeneralConfig.IdentifierType = []int{102, 103, 104}
	mirConfig.GeneralConfig.DefaultRouteConfigPath = "/usr/local/etc/mir/defaultRoute.xml"
	mirConfig.GeneralConfig.DefaultRouteRetryCount = 3

	// Log
	mirConfig.LogConfig.LogLevel = "INFO"
	mirConfig.LogConfig.ReportCaller = true
	mirConfig.LogConfig.LogFormat = "text"
	mirConfig.LogConfig.LogFilePath = ""

	// table
	mirConfig.TableConfig.CSSize = 500
	mirConfig.TableConfig.CSReplaceStrategy = "LRU"
	mirConfig.TableConfig.CacheUnsolicitedData = false

	// LogicFace
	mirConfig.LogicFaceConfig.SupportTCP = true
	mirConfig.LogicFaceConfig.TCPPort = 13899
	mirConfig.LogicFaceConfig.SupportUDP = true
	mirConfig.LogicFaceConfig.UDPPort = 13899
	mirConfig.LogicFaceConfig.SupportUnix = true
	mirConfig.LogicFaceConfig.UnixPath = "/tmp/mir.sock"

	mirConfig.LFRecvQueSize = 10000
	mirConfig.LFSendQueSize = 10000

	// Security
	mirConfig.SecurityConfig.VerifyPacket = false
	mirConfig.SecurityConfig.Log2BlockChain = false
	mirConfig.SecurityConfig.MiddleRouterSignature = false
	mirConfig.MaxRouterSignatureNum = 4
	mirConfig.SecurityConfig.ParallelVerifyNum = 10
	mirConfig.SecurityConfig.IdentityDBPath = security.DefaultIdentityDBPath

	// Forwarder
	mirConfig.ForwarderConfig.PacketQueueSize = 100

	// Strategy
	mirConfig.StrategyConfig.RoundRobinStrategyPrefix = "/rrs"
	mirConfig.StrategyConfig.RoundRobinStrategyRoundTime = 600
	mirConfig.StrategyConfig.EnableRoundRobinStrategy = false
}

// Save 保存当前配置状态到配置文件当中
//
// @Description:
// @receiver mirConfig
// @return error
//
func (mirConfig *MIRConfig) Save() error {
	cfg := ini.Empty()
	if err := ini.ReflectFrom(cfg, mirConfig); err != nil {
		return err
	}
	return cfg.SaveTo(mirConfig.configPath)
}

type GeneralConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// General
	////////////////////////////////////////////////////////////////////////////////////////////////
	DefaultId               string `ini:"DefaultId"`               // 默认网络身份
	EncryptedPasswdSavePath string `ini:"EncryptedPasswdSavePath"` // 加密秘钥保存位置
	IdentifierType          []int  `ini:"IdentifierType"`          // 当前路由器支持的标识类型，102 => GPPkt | 103 => 内容兴趣标识（Interest）| 104 => 内容兴趣标识（Interest）
	DefaultRouteConfigPath  string `ini:"DefaultRouteConfigPath"`  // 静态路由配置文件路径
	DefaultRouteRetryCount  int    `ini:"DefaultRouteRetryCount"`  // 静态路由创建重试次数
}

type LogConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Log
	////////////////////////////////////////////////////////////////////////////////////////////////
	LogLevel     string `ini:"LogLevel"`     // 日志等级
	ReportCaller bool   `ini:"ReportCaller"` // 日志输出时是否添加文件名和函数名
	LogFormat    string `ini:"LogFormat"`    // 输出日志的格式 "json" | "text"
	LogFilePath  string `ini:"LogFilePath"`  // 日志输出文件路径，为空则输出至控制台
}

type TableConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Table
	////////////////////////////////////////////////////////////////////////////////////////////////
	CSSize               int    `ini:"CSSize"`               // CS缓存大小，包为单位
	CSReplaceStrategy    string `ini:"CSReplaceStrategy"`    // 缓存替换策略
	CacheUnsolicitedData bool   `ini:"CacheUnsolicitedData"` // 是否缓存未请求的数据（Unsolicited Data）
}

type LogicFaceConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// LogicFace
	////////////////////////////////////////////////////////////////////////////////////////////////
	SupportTCP                 bool   `ini:"SupportTCP"`                 // 是否开启TCP
	TCPPort                    int    `ini:"TCPPort"`                    // TCP 端口号
	SupportUDP                 bool   `ini:"SupportUDP"`                 // 是否开启UDP
	UDPPort                    int    `ini:"UDPPort"`                    // UDP 端口号
	SupportUnix                bool   `ini:"SupportUnix"`                // 是否开启Unix
	UnixPath                   string `ini:"UnixPath"`                   // Unix 套接字路径设置
	LogicFaceIdleTime          int    `ini:"LogicFaceIdleTime"`          // LogicFace最大闲置时间
	CleanLogicFaceTableTimeVal int    `ini:"CleanLogicFaceTableTimeVal"` // LogicFaceSystem 清理逻辑接口的时间周期
	EtherRoutineNumber         int    `ini:"EtherRoutineNumber"`         // 以一个网卡对应的收包协程数
	UDPReceiveRoutineNumber    int    `ini:"UDPReceiveRoutineNumber"`    //UDP收包协程数
	LFRecvQueSize              int    `ini:"LFRecvQueSize"`              //	接收队列大小
	LFSendQueSize              int    `ini:"LFSendQueSize"`              // 发送队列大小
}

type SecurityConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Security
	////////////////////////////////////////////////////////////////////////////////////////////////
	VerifyPacket          bool   `ini:"VerifyPacket"`          // 是否开启包签名验证
	Log2BlockChain        bool   `ini:"Log2BlockChain"`        // 是否发送日志到区块链
	MiddleRouterSignature bool   `ini:"MiddleRouterSignature"` //是否开启中间路由器签名
	MaxRouterSignatureNum int    `ini:"MaxRouterSignatureNum"` // 最大中间路由器签名数量
	ParallelVerifyNum     int    `ini:"ParallelVerifyNum"`     // 并行包验证协程数量
	IdentityDBPath        string `ini:"IdentityDBPath"`        // 身份持久化sqlite数据库存储位置
}

type ForwarderConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Forwarder
	////////////////////////////////////////////////////////////////////////////////////////////////
	PacketQueueSize int `ini:"PacketQueueSize"` // 包缓冲队列大小
}

type StrategyConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Strategy
	////////////////////////////////////////////////////////////////////////////////////////////////
	RoundRobinStrategyPrefix    string `ini:"RoundRobinStrategyPrefix"`    // 轮询策略生效的前缀（例如：/rrs开头的包全部都会走轮询策略）=> 默认rrs
	RoundRobinStrategyRoundTime int    `ini:"RoundRobinStrategyRoundTime"` //轮询策略轮换的时间（单位为秒）=> 默认10分钟
	EnableRoundRobinStrategy    bool   `ini:"EnableRoundRobinStrategy"`    //是否开启轮询策略
}

type ManagementConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// Management
	////////////////////////////////////////////////////////////////////////////////////////////////
	CacheSize int64 `ini:"CacheSize"` // 缓存大小
}

type PcapConfig struct {
	////////////////////////////////////////////////////////////////////////////////////////////////
	//// PcapMode
	////////////////////////////////////////////////////////////////////////////////////////////////
	SetImmediateMode bool  `ini:"SetImmediateMode"` // 是否开启立即模式
	Promiscuous      bool  `ini:"Promiscuous"`      // 是否开启混杂模式
	PcapReadTimeout  int64 `ini:"PcapReadTimeout"`  // 超时时间，-1表示不超时，没有数据就卡住等待
	PcapBufferSize   int   `ini:"PcapBufferSize"`   // libpcap 抓包时的缓冲区大小 4 * 1024 * 1024 => 4194304
}

// ParseConfig
// 解析配置文件
//
// @Description:
// @receiver m
// @param configPath
// @return error
//
func ParseConfig(configPath string) (*MIRConfig, error) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		return nil, err
	}
	mirConfig := new(MIRConfig)
	mirConfig.configPath = configPath
	// 初始化配置，给所有的配置项设置默认值
	mirConfig.Init()
	// 加载配置文件中的配置
	if err = cfg.MapTo(&mirConfig); err != nil {
		return nil, err
	}
	return mirConfig, nil
}
