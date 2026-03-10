package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// TimeoutLock 无死锁的超时锁（核心：用原子状态+非嵌套加锁）
type TimeoutLock struct {
	mu      sync.Mutex    // 仅用于保护状态，不做业务锁
	locked  bool          // 是否已加锁
	expired bool          // 是否已超时
	timeout time.Duration // 超时时间
	timer   *time.Timer   // 超时计时器
}

// NewTimeoutLock 创建超时锁实例
func NewTimeoutLock(timeout time.Duration) *TimeoutLock {
	return &TimeoutLock{
		timeout: timeout,
	}
}

// Lock 加锁（非阻塞检查+无嵌套加锁，避免死锁）
// 返回值：true=加锁成功，false=已超时/已被锁定
func (tl *TimeoutLock) Lock() error {
	// 单锁保护所有状态，避免嵌套加锁
	tl.mu.Lock()
	defer tl.mu.Unlock()

	// 1. 检查超时/已锁定状态
	if tl.expired {
		return errors.New("锁已超时，无法加锁")
	}
	if tl.locked {
		return errors.New("锁已被持有，请勿重复加锁")
	}

	// 2. 标记为已锁定
	tl.locked = true

	// 3. 启动超时计时器（先停止旧计时器）
	if tl.timer != nil {
		tl.timer.Stop()
	}
	tl.timer = time.AfterFunc(tl.timeout, func() {
		// 计时器内仅修改状态，不做复杂操作
		tl.mu.Lock()
		tl.expired = true
		tl.locked = false // 超时后自动解锁
		tl.mu.Unlock()
		fmt.Println("[超时提醒] 锁已超时，自动释放并失效")
	})

	Debug("✅ 加锁成功")
	return nil
}

// Unlock 解锁（安全释放，无死锁）
func (tl *TimeoutLock) Unlock() error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	// 检查状态
	if !tl.locked && !tl.expired {
		return errors.New("未持有锁，无需释放")
	}
	if tl.expired {
		return errors.New("锁已超时，释放操作无效")
	}

	// 停止计时器+重置状态
	if tl.timer != nil {
		tl.timer.Stop()
		tl.timer = nil
	}
	tl.locked = false
	tl.expired = false

	Debug("✅ 解锁成功")
	return nil
}

// TryLock 非阻塞尝试加锁（可选，避免长时间阻塞）
func (tl *TimeoutLock) TryLock() bool {
	err := tl.Lock()
	return err == nil
}
