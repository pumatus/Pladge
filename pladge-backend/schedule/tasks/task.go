package tasks

import (
	"pledge-backend/db"
	"pledge-backend/schedule/common"
	"pledge-backend/schedule/services"
	"time"

	"github.com/jasonlvhit/gocron"
)

func Task() {

	// get environment variables
	common.GetEnv()

	// flush redis db
	err := db.RedisFlushDB()
	if err != nil {
		panic("clear redis error " + err.Error())
	}

	//init task
	services.NewPool().UpdateAllPoolInfo()
	services.NewTokenPrice().UpdateContractPrice()
	services.NewTokenSymbol().UpdateContractSymbol()
	services.NewTokenLogo().UpdateTokenLogo()
	services.NewBalanceMonitor().Monitor()
	// services.NewTokenPrice().SavePlgrPrice()
	services.NewTokenPrice().SavePlgrPriceTestNet()

	//run pool task
	s := gocron.NewScheduler()
	s.ChangeLoc(time.UTC)
	_ = s.Every(2).Minutes().From(gocron.NextTick()).Do(services.NewPool().UpdateAllPoolInfo)
	_ = s.Every(1).Minute().From(gocron.NextTick()).Do(services.NewTokenPrice().UpdateContractPrice)
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenSymbol().UpdateContractSymbol)
	_ = s.Every(2).Hours().From(gocron.NextTick()).Do(services.NewTokenLogo().UpdateTokenLogo)
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewBalanceMonitor().Monitor)
	//_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPrice)
	_ = s.Every(30).Minutes().From(gocron.NextTick()).Do(services.NewTokenPrice().SavePlgrPriceTestNet)
	<-s.Start() // Start all the pending jobs

}

// 任务,				频率,				理由
// UpdateContractPrice,1 分钟,价格波动最快，需要最高频更新。
// UpdateAllPoolInfo,2 分钟,质押池的存借款状态变动较快，需紧密同步。
// Monitor (余额),30 分钟,手续费消耗没那么快，半小时检查一次足够安全。
// SavePlgrPrice,30 分钟,喂价涉及写链（耗 Gas），通常半小时更新一次价格中值即可。
// UpdateTokenLogo/Symbol,2 小时,代币的名字和图标几乎万年不变，低频更新节省 IO。

// 从下至上：MySQL/Redis -> ABI/Env -> Models -> Services -> Tasks。

// 从静态到动态：从配置读取到实时链上同步。
