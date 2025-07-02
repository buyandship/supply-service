package db

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/buyandship/supply-svr/biz/common/config"
	"github.com/buyandship/supply-svr/biz/common/trace"
	model "github.com/buyandship/supply-svr/biz/model/mercari"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	once    sync.Once
	Handler *H
)

type H struct {
	cli *gorm.DB
}

func GetHandler() *H {
	once.Do(func() {
		c, err := gorm.Open(mysql.New(
			mysql.Config{
				DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
					config.GlobalServerConfig.Mysql.Username,
					config.GlobalServerConfig.Mysql.Password,
					config.GlobalServerConfig.Mysql.Address,
					config.GlobalServerConfig.Mysql.DBName,
				),
			},
		), &gorm.Config{})
		if err != nil {
			hlog.Fatal("mysql init err:", err)
		}
		Handler = &H{cli: c}
	})
	return Handler
}

func (h *H) HealthCheck() error {
	sqlDB, err := h.cli.DB()
	if err != nil {
		return err
	}
	if err := sqlDB.Ping(); err != nil {
		return err
	}
	return nil
}

func (h *H) UpsertAccount(ctx context.Context, account *model.Account) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "UpsertAccount")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	err = sql.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&account).Error

	return err
}

func (h *H) BanAccount(ctx context.Context, accountId int32) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "BanAccount")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	updates := map[string]any{
		"banned_at": time.Now(),
	}

	return sql.Model(&model.Account{}).Where("id = ?", accountId).Updates(updates).Error
}

func (h *H) InsertMessage(ctx context.Context, message *model.Message) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertMessage")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	return sql.Create(&message).Error
}

func (h *H) InsertReview(ctx context.Context, review *model.Review) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertReview")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	return sql.Create(&review).Error
}

func (h *H) GetTransaction(ctx context.Context, where *model.Transaction) (trx *model.Transaction, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetTransaction")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	err = sql.Where(where).First(&trx).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (h *H) InsertTransaction(ctx context.Context, transaction *model.Transaction) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertTransaction")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	return sql.Create(&transaction).Error
}

func (h *H) InsertTokenLog(ctx context.Context, token *model.Token) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "InsertTokenLog")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	return sql.Create(&token).Error
}

func (h *H) UpdateTransaction(ctx context.Context, cond *model.Transaction) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "UpdateTransaction")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	return sql.Where("ref_id = ?", cond.RefID).Updates(cond).Error
}

func (h *H) GetAccount(ctx context.Context, id int32) (account *model.Account, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetAccount")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	err = sql.Where("id = ?", id).First(&account).Error

	if err != nil {
		hlog.CtxErrorf(ctx, "get account error: %s", err.Error())
		return nil, err
	}
	return
}

func (h *H) GetToken(ctx context.Context, accountId int32) (token *model.Token, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetToken")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	err = sql.Where("account_id = ?", accountId).
		Order("created_at desc").
		First(&token).Error

	if err != nil {
		return nil, err
	}
	return
}

func (h *H) GetAccountList(ctx context.Context) (accounts []*model.Account, err error) {
	ctx, span := trace.StartDBOperation(ctx, "GetAccountList")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	err = sql.Order("priority asc").Find(&accounts).Error

	if err != nil {
		return nil, err
	}
	return
}

func (h *H) SwitchAccount(ctx context.Context, accountId int32) (err error) {
	ctx, span := trace.StartDBOperation(ctx, "SwitchAccount")
	defer trace.EndSpan(span, err)

	sql := h.cli.WithContext(ctx)

	if config.GlobalServerConfig.Env == "development" {
		sql = sql.Debug()
	}

	tx := sql.Begin()

	now := time.Now()

	if err := tx.
		Model(&model.Account{}).
		Where("active_at is not null").
		Update("active_at", nil).Error; err != nil {
		hlog.CtxErrorf(ctx, "failed to update account active_at: %v", err)
		tx.Rollback()
		return err
	}

	if err := tx.
		Where("id = ?", accountId).
		Updates(&model.Account{ActiveAt: &now}).Error; err != nil {
		hlog.CtxErrorf(ctx, "failed to update account active_at: %v", err)
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		hlog.CtxErrorf(ctx, "failed to commit tx: %v", err)
		tx.Rollback()
		return err
	}

	return
}
