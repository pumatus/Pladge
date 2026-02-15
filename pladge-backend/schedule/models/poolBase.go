package models

import (
	"errors"
	"pledge-backend/db"
	"pledge-backend/log"
	"pledge-backend/utils"

	"gorm.io/gorm"
)

// 质押池基础信息 将链上的智能合约状态（State）同步到本地数据库中

// PoolBase: 存储质押池的静态和动态属性。
// SettleTime / EndTime: 结算和结束时间，决定了质押周期的生命周期。
// InterestRate: 利率，通常是借款成本或存款收益。
// MartgageRate: 质押率（LTV），决定了抵押多少资产能借出多少资产。
// LendToken / BorrowToken: 记录合约地址（Address）。
// SpCoin / JpCoin: 通常代表质押池生成的凭证代币（如 LP Token）。
// BorrowToken / LendToken: 这里的结构体看起来是用于给前端展示的 视图模型（VO），包含 Logo、价格等 UI 元素。
type PoolBase struct {
	Id                     int    `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	PoolId                 int    `json:"pool_id" gorm:"column:pool_id"`
	ChainId                string `json:"chain_id" gorm:"column:chain_id"`
	SettleTime             string `json:"settle_time" gorm:"column:settle_time"`
	EndTime                string `json:"end_time" gorm:"column:end_time"`
	InterestRate           string `json:"interest_rate" gorm:"column:interest_rate"`
	MaxSupply              string `json:"max_supply" gorm:"max_supply:"`
	LendSupply             string `json:"lend_supply" gorm:"column:lend_supply"`
	BorrowSupply           string `json:"borrow_supply" gorm:"column:borrow_supply"`
	MartgageRate           string `json:"martgage_rate" gorm:"column:martgage_rate"`
	LendToken              string `json:"lend_token" gorm:"column:lend_token"`
	LendTokenInfo          string `json:"lend_token_info" gorm:"column:lend_token_info"`
	BorrowToken            string `json:"borrow_token" gorm:"column:borrow_token"`
	BorrowTokenInfo        string `json:"borrow_token_info" gorm:"column:borrow_token_info"`
	State                  string `json:"state" gorm:"column:state"`
	SpCoin                 string `json:"sp_coin" gorm:"column:sp_coin"`
	JpCoin                 string `json:"jp_coin" gorm:"column:jp_coin"`
	LendTokenSymbol        string `json:"lend_token_symbol" gorm:"column:lend_token_symbol"`
	BorrowTokenSymbol      string `json:"borrow_token_symbol" gorm:"column:borrow_token_symbol"`
	AutoLiquidateThreshold string `json:"auto_liquidate_threshold" gorm:"column:auto_liquidate_threshold"`
	CreatedAt              string `json:"created_at" gorm:"column:created_at"`
	UpdatedAt              string `json:"updated_at" gorm:"column:updated_at"`
}

type BorrowToken struct {
	BorrowFee  string `json:"borrowFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

type LendToken struct {
	LendFee    string `json:"lendFee"`
	TokenLogo  string `json:"tokenLogo"`
	TokenName  string `json:"tokenName"`
	TokenPrice string `json:"tokenPrice"`
}

func NewPoolBase() *PoolBase {
	return &PoolBase{}
}

func (p *PoolBase) TableName() string {
	return "poolbases"
}

// SavePoolBase Save poolBase information
func (p *PoolBase) SavePoolBase(chainId, poolId string, poolBase *PoolBase) error {

	nowDateTime := utils.GetCurDateTimeFormat()

	//save token info
	err, symbol := p.SaveTokenInfo(poolBase)
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}
	poolBase.BorrowTokenSymbol = symbol[0]
	poolBase.LendTokenSymbol = symbol[1]

	// save pool info
	err = db.Mysql.Table("poolbases").Where("chain_id=? and pool_id=?", chainId, poolId).First(&p).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			poolBase.CreatedAt = nowDateTime
			poolBase.UpdatedAt = nowDateTime
			err = db.Mysql.Table("poolbases").Create(poolBase).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err
			}
		} else {
			return errors.New("record select err " + err.Error())
		}
	}

	poolBase.UpdatedAt = nowDateTime
	err = db.Mysql.Table("poolbases").Where("chain_id=? and pool_id=?", chainId, poolId).Updates(poolBase).Debug().Error
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *PoolBase) PoolBaseInfo(res *PoolBase) error {
	err := db.Mysql.Table("poolbases").Order("pool_id asc").Find(&res).Debug().Error
	if err != nil {
		log.Logger.Error(err.Error())
		return err
	}
	return nil
}

func (p *PoolBase) SaveTokenInfo(base *PoolBase) (error, []string) {
	tokenInfo := TokenInfo{}
	tokenSymbol := []string{"", ""}
	nowDateTime := utils.GetCurDateTimeFormat()

	// borrowToken
	err := db.Mysql.Table("token_info").Where("chain_id=? and token=?", base.ChainId, base.BorrowToken).First(&tokenInfo).Debug().Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			tokenInfo.Token = base.BorrowToken
			err = db.Mysql.Table("token_info").Create(&TokenInfo{
				Token:     base.BorrowToken,
				ChainId:   base.ChainId,
				CreatedAt: nowDateTime,
				UpdatedAt: nowDateTime,
			}).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err, tokenSymbol
			}
		} else {
			return errors.New("token_info record select err " + err.Error()), tokenSymbol
		}
	}
	tokenSymbol[0] = tokenInfo.Symbol

	//lendToken
	tokenInfo = TokenInfo{}
	err = db.Mysql.Table("token_info").Where("chain_id=? and token=?", base.ChainId, base.LendToken).First(&tokenInfo).Debug().Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = db.Mysql.Table("token_info").Create(&TokenInfo{
				Token:     base.LendToken,
				ChainId:   base.ChainId,
				CreatedAt: nowDateTime,
				UpdatedAt: nowDateTime,
			}).Debug().Error
			if err != nil {
				log.Logger.Error(err.Error())
				return err, tokenSymbol
			}
		} else {
			return errors.New("token_info record select err " + err.Error()), tokenSymbol
		}
	}
	tokenSymbol[1] = tokenInfo.Symbol

	return nil, tokenSymbol
}
