package project

import (
	"log"
	"sync"
	"time"

	pkgErrors "Qingyu_backend/pkg/errors"
)

// AutoSaveService 自动保存服务
// MVP实现：简单的定时保存机制
type AutoSaveService struct {
	versionService *VersionService
	interval       time.Duration
	sessions       map[string]*autoSaveSession
	mu             sync.RWMutex
}

// autoSaveSession 自动保存会话
type autoSaveSession struct {
	projectID  string
	documentID string
	nodeID     string
	userID     string
	ticker     *time.Ticker
	stopChan   chan struct{}
	lastSaved  time.Time
}

// NewAutoSaveService 创建自动保存服务
func NewAutoSaveService(versionService *VersionService) *AutoSaveService {
	return &AutoSaveService{
		versionService: versionService,
		interval:       30 * time.Second, // MVP: 固定30秒间隔
		sessions:       make(map[string]*autoSaveSession),
	}
}

// StartAutoSave 启动文档自动保存
// 参数：
//   - documentID: 文档ID（作为会话key）
//   - projectID: 项目ID
//   - nodeID: 节点ID
//   - userID: 用户ID
func (s *AutoSaveService) StartAutoSave(documentID, projectID, nodeID, userID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果已存在会话，先停止旧会话
	if session, exists := s.sessions[documentID]; exists {
		log.Printf("[AutoSave] 文档 %s 的自动保存已在运行，重启会话", documentID)
		close(session.stopChan)
		session.ticker.Stop()
	}

	// 创建新会话
	session := &autoSaveSession{
		projectID:  projectID,
		documentID: documentID,
		nodeID:     nodeID,
		userID:     userID,
		ticker:     time.NewTicker(s.interval),
		stopChan:   make(chan struct{}),
		lastSaved:  time.Now(),
	}

	s.sessions[documentID] = session

	// 启动后台定时保存
	go s.runAutoSave(session)

	log.Printf("[AutoSave] 启动文档自动保存: documentID=%s, interval=%v", documentID, s.interval)
	return nil
}

// StopAutoSave 停止文档自动保存
func (s *AutoSaveService) StopAutoSave(documentID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.sessions[documentID]
	if !exists {
		return pkgErrors.ProjectFactory.BusinessError(
			"AUTOSAVE_NOT_FOUND",
			"未找到自动保存会话",
		)
	}

	// 停止ticker和goroutine
	close(session.stopChan)
	session.ticker.Stop()
	delete(s.sessions, documentID)

	log.Printf("[AutoSave] 停止文档自动保存: documentID=%s", documentID)
	return nil
}

// runAutoSave 运行自动保存循环
func (s *AutoSaveService) runAutoSave(session *autoSaveSession) {
	for {
		select {
		case <-session.ticker.C:
			// 执行自动保存
			if err := s.performAutoSave(session); err != nil {
				log.Printf("[AutoSave] 自动保存失败: documentID=%s, error=%v",
					session.documentID, err)
				// MVP: 记录错误但不停止自动保存
			} else {
				session.lastSaved = time.Now()
				log.Printf("[AutoSave] 自动保存成功: documentID=%s, time=%v",
					session.documentID, session.lastSaved)
			}

		case <-session.stopChan:
			// 收到停止信号
			log.Printf("[AutoSave] 自动保存协程退出: documentID=%s", session.documentID)
			return
		}
	}
}

// performAutoSave 执行自动保存操作
func (s *AutoSaveService) performAutoSave(session *autoSaveSession) error {
	// MVP: 直接调用VersionService创建版本
	// 注意：这里不检查内容是否变化，每30秒必定保存（简化实现）
	_, err := s.versionService.BumpVersionAndCreateRevision(
		session.projectID,
		session.nodeID,
		session.userID,
		"自动保存", // 固定消息
	)

	if err != nil {
		return pkgErrors.ProjectFactory.InternalError(
			"AUTOSAVE_FAILED",
			"自动保存版本创建失败",
			err,
		)
	}

	return nil
}

// GetStatus 获取自动保存状态
// 返回：是否正在自动保存、最后保存时间
func (s *AutoSaveService) GetStatus(documentID string) (isRunning bool, lastSaved *time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[documentID]
	if !exists {
		return false, nil
	}

	return true, &session.lastSaved
}

// GetActiveSessions 获取所有活跃的自动保存会话数量（用于监控）
func (s *AutoSaveService) GetActiveSessions() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

// StopAll 停止所有自动保存会话（用于服务关闭）
func (s *AutoSaveService) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for documentID, session := range s.sessions {
		close(session.stopChan)
		session.ticker.Stop()
		log.Printf("[AutoSave] 停止自动保存会话: %s", documentID)
	}

	s.sessions = make(map[string]*autoSaveSession)
	log.Println("[AutoSave] 所有自动保存会话已停止")
}

// UpdateInterval 更新自动保存间隔（仅影响新启动的会话）
// MVP: 简化实现，不影响已运行的会话
func (s *AutoSaveService) UpdateInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if interval < 10*time.Second {
		log.Println("[AutoSave] 自动保存间隔过短，设置为最小值10秒")
		interval = 10 * time.Second
	}

	s.interval = interval
	log.Printf("[AutoSave] 自动保存间隔已更新为: %v（仅影响新会话）", interval)
}

// --- 以下为未来增强功能的占位符 ---

// TODO: Phase 2 增强功能
// - [ ] 内容变更检测（只在内容变化时保存）
// - [ ] 变更字符数阈值触发（100字符变更触发保存）
// - [ ] 离线保存队列（网络断开时排队）
// - [ ] 保存失败重试机制
// - [ ] 增量保存（只保存diff）
// - [ ] 保存前的内容校验
// - [ ] 冲突检测（乐观锁）
// - [ ] 保存进度通知（WebSocket推送）
// - [ ] 用户手动保存与自动保存协调
// - [ ] 自动保存历史清理策略

// SaveImmediately 立即执行一次保存（手动触发）
// MVP暂不实现，预留接口
func (s *AutoSaveService) SaveImmediately(documentID string) error {
	s.mu.RLock()
	session, exists := s.sessions[documentID]
	s.mu.RUnlock()

	if !exists {
		return pkgErrors.ProjectFactory.BusinessError(
			"AUTOSAVE_NOT_FOUND",
			"未找到自动保存会话",
		)
	}

	return s.performAutoSave(session)
}

// --- 测试辅助方法 ---

// GetSession 获取会话信息（仅用于测试）
func (s *AutoSaveService) GetSession(documentID string) *autoSaveSession {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[documentID]
}
