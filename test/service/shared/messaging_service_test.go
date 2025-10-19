package shared

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockMessage 模拟消息模型
type MockMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Topic     string             `bson:"topic" json:"topic"`
	Key       string             `bson:"key" json:"key"`
	Value     []byte             `bson:"value" json:"value"`
	Headers   map[string]string  `bson:"headers" json:"headers"`
	Partition int32              `bson:"partition" json:"partition"`
	Offset    int64              `bson:"offset" json:"offset"`
	Status    string             `bson:"status" json:"status"`
	Retries   int                `bson:"retries" json:"retries"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockTopic 模拟主题模型
type MockTopic struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	Partitions int32              `bson:"partitions" json:"partitions"`
	Replicas   int16              `bson:"replicas" json:"replicas"`
	Config     map[string]string  `bson:"config" json:"config"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockConsumerGroup 模拟消费者组模型
type MockConsumerGroup struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Topics    []string           `bson:"topics" json:"topics"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockConsumer 模拟消费者模型
type MockConsumer struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID       primitive.ObjectID `bson:"group_id" json:"groupId"`
	ConsumerID    string             `bson:"consumer_id" json:"consumerId"`
	Topic         string             `bson:"topic" json:"topic"`
	Partition     int32              `bson:"partition" json:"partition"`
	Offset        int64              `bson:"offset" json:"offset"`
	LastHeartbeat time.Time          `bson:"last_heartbeat" json:"lastHeartbeat"`
	Status        string             `bson:"status" json:"status"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

// MockMessageHandler 模拟消息处理器
type MockMessageHandler func(ctx context.Context, message *MockMessage) error

// MockMessageRepository 模拟消息仓储
type MockMessageRepository struct {
	mock.Mock
}

func (m *MockMessageRepository) Create(ctx context.Context, message *MockMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockMessageRepository) GetByTopic(ctx context.Context, topic string, limit, offset int) ([]*MockMessage, error) {
	args := m.Called(ctx, topic, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockMessage), args.Error(1)
}

func (m *MockMessageRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockMessageRepository) GetPendingMessages(ctx context.Context, topic string, limit int) ([]*MockMessage, error) {
	args := m.Called(ctx, topic, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockMessage), args.Error(1)
}

// MockTopicRepository 模拟主题仓储
type MockTopicRepository struct {
	mock.Mock
}

func (m *MockTopicRepository) Create(ctx context.Context, topic *MockTopic) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

func (m *MockTopicRepository) GetByName(ctx context.Context, name string) (*MockTopic, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockTopic), args.Error(1)
}

func (m *MockTopicRepository) List(ctx context.Context) ([]*MockTopic, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*MockTopic), args.Error(1)
}

func (m *MockTopicRepository) Delete(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

// MockConsumerRepository 模拟消费者仓储
type MockConsumerRepository struct {
	mock.Mock
}

func (m *MockConsumerRepository) CreateGroup(ctx context.Context, group *MockConsumerGroup) error {
	args := m.Called(ctx, group)
	return args.Error(0)
}

func (m *MockConsumerRepository) GetGroup(ctx context.Context, name string) (*MockConsumerGroup, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*MockConsumerGroup), args.Error(1)
}

func (m *MockConsumerRepository) CreateConsumer(ctx context.Context, consumer *MockConsumer) error {
	args := m.Called(ctx, consumer)
	return args.Error(0)
}

func (m *MockConsumerRepository) UpdateOffset(ctx context.Context, consumerID string, topic string, partition int32, offset int64) error {
	args := m.Called(ctx, consumerID, topic, partition, offset)
	return args.Error(0)
}

// MockMessageBroker 模拟消息代理
type MockMessageBroker struct {
	mock.Mock
}

func (m *MockMessageBroker) Publish(ctx context.Context, topic string, key string, value []byte, headers map[string]string) error {
	args := m.Called(ctx, topic, key, value, headers)
	return args.Error(0)
}

func (m *MockMessageBroker) Subscribe(ctx context.Context, topic string, groupID string, handler MockMessageHandler) error {
	args := m.Called(ctx, topic, groupID, handler)
	return args.Error(0)
}

func (m *MockMessageBroker) CreateTopic(ctx context.Context, topic string, partitions int32, replicas int16) error {
	args := m.Called(ctx, topic, partitions, replicas)
	return args.Error(0)
}

func (m *MockMessageBroker) DeleteTopic(ctx context.Context, topic string) error {
	args := m.Called(ctx, topic)
	return args.Error(0)
}

// MockMessagingService 模拟消息服务
type MockMessagingService struct {
	messageRepo  *MockMessageRepository
	topicRepo    *MockTopicRepository
	consumerRepo *MockConsumerRepository
	broker       *MockMessageBroker
	handlers     map[string]MockMessageHandler
}

func NewMockMessagingService(
	messageRepo *MockMessageRepository,
	topicRepo *MockTopicRepository,
	consumerRepo *MockConsumerRepository,
	broker *MockMessageBroker,
) *MockMessagingService {
	return &MockMessagingService{
		messageRepo:  messageRepo,
		topicRepo:    topicRepo,
		consumerRepo: consumerRepo,
		broker:       broker,
		handlers:     make(map[string]MockMessageHandler),
	}
}

// CreateTopic 创建主题
func (s *MockMessagingService) CreateTopic(ctx context.Context, name string, partitions int32, replicas int16, config map[string]string) (*MockTopic, error) {
	if name == "" {
		return nil, errors.New("topic name is required")
	}

	if partitions <= 0 {
		return nil, errors.New("partitions must be greater than 0")
	}

	// 检查主题是否已存在
	existingTopic, _ := s.topicRepo.GetByName(ctx, name)
	if existingTopic != nil {
		return nil, errors.New("topic already exists")
	}

	// 在消息代理中创建主题
	err := s.broker.CreateTopic(ctx, name, partitions, replicas)
	if err != nil {
		return nil, err
	}

	// 创建主题记录
	topic := &MockTopic{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Partitions: partitions,
		Replicas:   replicas,
		Config:     config,
		Status:     "active",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err = s.topicRepo.Create(ctx, topic)
	if err != nil {
		return nil, err
	}

	return topic, nil
}

// PublishMessage 发布消息
func (s *MockMessagingService) PublishMessage(ctx context.Context, topic, key string, payload interface{}, headers map[string]string) error {
	if topic == "" {
		return errors.New("topic is required")
	}

	// 检查主题是否存在
	_, err := s.topicRepo.GetByName(ctx, topic)
	if err != nil {
		return errors.New("topic not found")
	}

	// 序列化消息
	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 发布到消息代理
	err = s.broker.Publish(ctx, topic, key, value, headers)
	if err != nil {
		return err
	}

	// 创建消息记录
	message := &MockMessage{
		ID:        primitive.NewObjectID(),
		Topic:     topic,
		Key:       key,
		Value:     value,
		Headers:   headers,
		Status:    "published",
		Retries:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.messageRepo.Create(ctx, message)
}

// SubscribeToTopic 订阅主题
func (s *MockMessagingService) SubscribeToTopic(ctx context.Context, topic, groupID string, handler MockMessageHandler) error {
	if topic == "" {
		return errors.New("topic is required")
	}

	if groupID == "" {
		return errors.New("group ID is required")
	}

	// 检查主题是否存在
	_, err := s.topicRepo.GetByName(ctx, topic)
	if err != nil {
		return errors.New("topic not found")
	}

	// 创建或获取消费者组
	group, err := s.consumerRepo.GetGroup(ctx, groupID)
	if err != nil {
		group = &MockConsumerGroup{
			ID:        primitive.NewObjectID(),
			Name:      groupID,
			Topics:    []string{topic},
			Status:    "active",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err = s.consumerRepo.CreateGroup(ctx, group)
		if err != nil {
			return err
		}
	}

	// 注册处理器
	s.handlers[topic+":"+groupID] = handler

	// 订阅消息代理
	return s.broker.Subscribe(ctx, topic, groupID, handler)
}

// ProcessPendingMessages 处理待处理消息
func (s *MockMessagingService) ProcessPendingMessages(ctx context.Context, topic string, limit int) error {
	messages, err := s.messageRepo.GetPendingMessages(ctx, topic, limit)
	if err != nil {
		return err
	}

	for _, message := range messages {
		// 查找处理器
		var handler MockMessageHandler
		for key, h := range s.handlers {
			if key == topic+":default" || key == topic {
				handler = h
				break
			}
		}

		if handler != nil {
			err := handler(ctx, message)
			if err != nil {
				// 更新重试次数
				message.Retries++
				s.messageRepo.UpdateStatus(ctx, message.ID, "failed")
			} else {
				s.messageRepo.UpdateStatus(ctx, message.ID, "processed")
			}
		}
	}

	return nil
}

// GetTopicMessages 获取主题消息
func (s *MockMessagingService) GetTopicMessages(ctx context.Context, topic string, limit, offset int) ([]*MockMessage, error) {
	return s.messageRepo.GetByTopic(ctx, topic, limit, offset)
}

// DeleteTopic 删除主题
func (s *MockMessagingService) DeleteTopic(ctx context.Context, name string) error {
	// 从消息代理删除
	err := s.broker.DeleteTopic(ctx, name)
	if err != nil {
		return err
	}

	// 删除主题记录
	return s.topicRepo.Delete(ctx, name)
}

// GetTopics 获取主题列表
func (s *MockMessagingService) GetTopics(ctx context.Context) ([]*MockTopic, error) {
	return s.topicRepo.List(ctx)
}

// 测试用例

func TestMessagingService_CreateTopic_Success(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	name := "test-topic"
	partitions := int32(3)
	replicas := int16(1)
	config := map[string]string{"retention.ms": "86400000"}

	// Mock 设置
	topicRepo.On("GetByName", ctx, name).Return(nil, errors.New("not found"))
	broker.On("CreateTopic", ctx, name, partitions, replicas).Return(nil)
	topicRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockTopic")).Return(nil)

	// 执行测试
	topic, err := service.CreateTopic(ctx, name, partitions, replicas, config)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, topic)
	assert.Equal(t, name, topic.Name)
	assert.Equal(t, partitions, topic.Partitions)
	assert.Equal(t, replicas, topic.Replicas)
	assert.Equal(t, config, topic.Config)
	assert.Equal(t, "active", topic.Status)

	topicRepo.AssertExpectations(t)
	broker.AssertExpectations(t)
}

func TestMessagingService_CreateTopic_AlreadyExists(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	name := "existing-topic"
	partitions := int32(3)
	replicas := int16(1)
	config := map[string]string{}

	existingTopic := &MockTopic{
		ID:         primitive.NewObjectID(),
		Name:       name,
		Partitions: 3,
		Replicas:   1,
		Status:     "active",
	}

	// Mock 设置
	topicRepo.On("GetByName", ctx, name).Return(existingTopic, nil)

	// 执行测试
	topic, err := service.CreateTopic(ctx, name, partitions, replicas, config)

	// 断言
	assert.Error(t, err)
	assert.Nil(t, topic)
	assert.Equal(t, "topic already exists", err.Error())

	topicRepo.AssertExpectations(t)
}

func TestMessagingService_PublishMessage_Success(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	topic := "test-topic"
	key := "test-key"
	payload := map[string]interface{}{
		"message":   "Hello, World!",
		"timestamp": time.Now().Unix(),
	}
	headers := map[string]string{
		"source": "test",
	}

	existingTopic := &MockTopic{
		ID:     primitive.NewObjectID(),
		Name:   topic,
		Status: "active",
	}

	expectedValue, _ := json.Marshal(payload)

	// Mock 设置
	topicRepo.On("GetByName", ctx, topic).Return(existingTopic, nil)
	broker.On("Publish", ctx, topic, key, expectedValue, headers).Return(nil)
	messageRepo.On("Create", ctx, mock.AnythingOfType("*shared.MockMessage")).Return(nil)

	// 执行测试
	err := service.PublishMessage(ctx, topic, key, payload, headers)

	// 断言
	assert.NoError(t, err)

	topicRepo.AssertExpectations(t)
	broker.AssertExpectations(t)
	messageRepo.AssertExpectations(t)
}

func TestMessagingService_PublishMessage_TopicNotFound(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	topic := "nonexistent-topic"
	key := "test-key"
	payload := map[string]interface{}{"message": "Hello"}
	headers := map[string]string{}

	// Mock 设置
	topicRepo.On("GetByName", ctx, topic).Return(nil, errors.New("not found"))

	// 执行测试
	err := service.PublishMessage(ctx, topic, key, payload, headers)

	// 断言
	assert.Error(t, err)
	assert.Equal(t, "topic not found", err.Error())

	topicRepo.AssertExpectations(t)
}

func TestMessagingService_SubscribeToTopic_Success(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	topic := "test-topic"
	groupID := "test-group"

	existingTopic := &MockTopic{
		ID:     primitive.NewObjectID(),
		Name:   topic,
		Status: "active",
	}

	handler := func(ctx context.Context, message *MockMessage) error {
		return nil
	}

	// Mock 设置
	topicRepo.On("GetByName", ctx, topic).Return(existingTopic, nil)
	consumerRepo.On("GetGroup", ctx, groupID).Return(nil, errors.New("not found"))
	consumerRepo.On("CreateGroup", ctx, mock.AnythingOfType("*shared.MockConsumerGroup")).Return(nil)
	broker.On("Subscribe", ctx, topic, groupID, mock.AnythingOfType("shared.MockMessageHandler")).Return(nil)

	// 执行测试
	err := service.SubscribeToTopic(ctx, topic, groupID, handler)

	// 断言
	assert.NoError(t, err)

	topicRepo.AssertExpectations(t)
	consumerRepo.AssertExpectations(t)
	broker.AssertExpectations(t)
}

func TestMessagingService_ProcessPendingMessages_Success(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	topic := "test-topic"
	limit := 10

	pendingMessages := []*MockMessage{
		{
			ID:      primitive.NewObjectID(),
			Topic:   topic,
			Key:     "key1",
			Value:   []byte(`{"message": "test1"}`),
			Status:  "pending",
			Retries: 0,
		},
		{
			ID:      primitive.NewObjectID(),
			Topic:   topic,
			Key:     "key2",
			Value:   []byte(`{"message": "test2"}`),
			Status:  "pending",
			Retries: 0,
		},
	}

	// 注册处理器
	handlerCalled := 0
	handler := func(ctx context.Context, message *MockMessage) error {
		handlerCalled++
		return nil
	}
	service.handlers[topic] = handler

	// Mock 设置
	messageRepo.On("GetPendingMessages", ctx, topic, limit).Return(pendingMessages, nil)
	messageRepo.On("UpdateStatus", ctx, pendingMessages[0].ID, "processed").Return(nil)
	messageRepo.On("UpdateStatus", ctx, pendingMessages[1].ID, "processed").Return(nil)

	// 执行测试
	err := service.ProcessPendingMessages(ctx, topic, limit)

	// 断言
	assert.NoError(t, err)
	assert.Equal(t, 2, handlerCalled)

	messageRepo.AssertExpectations(t)
}

func TestMessagingService_GetTopicMessages_Success(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	topic := "test-topic"
	limit := 10
	offset := 0

	expectedMessages := []*MockMessage{
		{
			ID:     primitive.NewObjectID(),
			Topic:  topic,
			Key:    "key1",
			Value:  []byte(`{"message": "test1"}`),
			Status: "published",
		},
		{
			ID:     primitive.NewObjectID(),
			Topic:  topic,
			Key:    "key2",
			Value:  []byte(`{"message": "test2"}`),
			Status: "processed",
		},
	}

	// Mock 设置
	messageRepo.On("GetByTopic", ctx, topic, limit, offset).Return(expectedMessages, nil)

	// 执行测试
	messages, err := service.GetTopicMessages(ctx, topic, limit, offset)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, messages)
	assert.Len(t, messages, 2)
	assert.Equal(t, expectedMessages, messages)

	messageRepo.AssertExpectations(t)
}

func TestMessagingService_DeleteTopic_Success(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()
	name := "test-topic"

	// Mock 设置
	broker.On("DeleteTopic", ctx, name).Return(nil)
	topicRepo.On("Delete", ctx, name).Return(nil)

	// 执行测试
	err := service.DeleteTopic(ctx, name)

	// 断言
	assert.NoError(t, err)

	broker.AssertExpectations(t)
	topicRepo.AssertExpectations(t)
}

func TestMessagingService_GetTopics_Success(t *testing.T) {
	messageRepo := new(MockMessageRepository)
	topicRepo := new(MockTopicRepository)
	consumerRepo := new(MockConsumerRepository)
	broker := new(MockMessageBroker)
	service := NewMockMessagingService(messageRepo, topicRepo, consumerRepo, broker)

	ctx := context.Background()

	expectedTopics := []*MockTopic{
		{
			ID:         primitive.NewObjectID(),
			Name:       "topic1",
			Partitions: 3,
			Replicas:   1,
			Status:     "active",
		},
		{
			ID:         primitive.NewObjectID(),
			Name:       "topic2",
			Partitions: 1,
			Replicas:   1,
			Status:     "active",
		},
	}

	// Mock 设置
	topicRepo.On("List", ctx).Return(expectedTopics, nil)

	// 执行测试
	topics, err := service.GetTopics(ctx)

	// 断言
	assert.NoError(t, err)
	assert.NotNil(t, topics)
	assert.Len(t, topics, 2)
	assert.Equal(t, expectedTopics, topics)

	topicRepo.AssertExpectations(t)
}
