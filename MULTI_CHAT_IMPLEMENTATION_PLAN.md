# Multi-Chat Support Implementation Plan

## Overview
This document outlines the complete implementation plan for adding multi-chat support to the bot. The goal is to enable the bot to respond to multiple Telegram chats, with separate configurations, matchers, data storage, and file sets for each chat.

## Current State Analysis

### Existing Infrastructure
- **bot-matcher v0.1.6**: Already has per-matcher enable/disable support via `matcher.Config` with `enabled *bool` field
- **Single chat configuration**: Currently using `telegram.chatID` (string) in config
- **Database tables**:
  - `plusplus` - stores ++/-- values by name (no chat separation)
  - `stats` - stores user post counts and last post time (no chat separation)
  - `message_stats` - stores per-message word counts (no chat separation)
- **File-based matchers**: Fortune files stored in `files/fortune/` (global, no chat separation)
- **Configuration-based matchers**: Buzzwords loaded from `config/buzzwords.yaml` (global config)

### Key Dependencies
- `github.com/br0-space/bot-matcher@v0.1.6` - Provides matcher registry and base matcher functionality
- `github.com/br0-space/bot-telegramclient@v0.1.4` - Telegram client and webhook handler
- GORM for database operations (SQLite or PostgreSQL)

---

## Implementation Plan

## Phase 1: Core Infrastructure Changes

### 1.1 Configuration Changes

**File: `interfaces/config.go`**
- Change `telegram.chatID` from single string to slice of chat configurations
- Add new struct for per-chat configuration:
```go
type ChatConfigStruct struct {
    ChatID      string
    Name        string  // Human-readable name for this chat
    Description string  // Optional description
}
```
- Update `ConfigStruct` to use:
```go
Telegram TelegramConfigStruct
Chats    []ChatConfigStruct
```

**File: `pkg/config/config.go`**
- Update environment variable mapping to support multiple chats
- Consider supporting both:
  - Simple mode: Single `TELEGRAM_CHAT_ID` for backwards compatibility
  - Multi mode: Array in config.yaml with chat configurations
- Add validation to ensure at least one chat is configured

**File: `config.yaml`**
- Add new `chats` section:
```yaml
chats:
  - chatID: "-1001234567890"
    name: "main-chat"
    description: "Main team chat"
  - chatID: "-1009876543210"
    name: "test-chat"
    description: "Testing environment"
```

**Migration Strategy for Config:**
- If `telegram.chatID` is set and `chats` is empty, auto-create single chat config
- Log warning about deprecated single chatID configuration
- Update `.env.dist` and documentation with new format

---

### 1.2 Database Schema Changes

All database tables need a `chat_id` field to separate data per chat.

#### 1.2.1 Plusplus Table
**File: `interfaces/repoplusplus.go`**
```go
type Plusplus struct {
    gorm.Model `exhaustruct:"optional"`

    ChatID string `gorm:"<-:create;index:idx_plusplus_chat_name,priority:1"`
    Name   string `gorm:"<-:create;index:idx_plusplus_chat_name,priority:2"`
    Value  int    `gorm:"<-;index"`
}
```
- Change unique constraint from `name` to composite `(chat_id, name)`
- Add index on `chat_id` for efficient filtering

**File: `pkg/repo/plusplus.go`**
- Add `chatID string` parameter to all methods:
  - `Increment(chatID, name string, increment int) (int, error)`
  - `FindTops(chatID string, limit int) ([]interfaces.Plusplus, error)`
  - `FindFlops(chatID string, limit int) ([]interfaces.Plusplus, error)`
- Update queries to filter by `chat_id`

#### 1.2.2 Stats Table
**File: `interfaces/repostats.go`**
```go
type Stats struct {
    gorm.Model `exhaustruct:"optional"`

    ChatID   string    `gorm:"<-:create;index:idx_stats_chat_user,priority:1"`
    UserID   int64     `gorm:"<-:create;index:idx_stats_chat_user,priority:2"`
    Username string    `gorm:"<-"`
    Posts    uint32    `gorm:"<-"`
    LastPost time.Time
}
```
- Change unique constraint from `user_id` to composite `(chat_id, user_id)`
- Same user can have different stats in different chats

**File: `pkg/repo/userstats.go`**
- Add `chatID string` parameter to all methods:
  - `UpdateStats(chatID string, userID int64, username string) error`
  - `GetKnownUsers(chatID string) ([]interfaces.StatsUserStruct, error)`
  - `GetTopUsers(chatID string) ([]interfaces.StatsUserStruct, error)`
- Update queries to filter by `chat_id`

#### 1.2.3 MessageStats Table
**File: `interfaces/repomessagestats.go`**
```go
type MessageStats struct {
    gorm.Model `exhaustruct:"optional"`

    ChatID string    `gorm:"<-:create;index:idx_msgstats_chat"`
    UserID int64     `gorm:"<-:create;index:idx_msgstats_user"`
    Time   time.Time `gorm:"<-:create;index"`
    Words  int       `gorm:"<-:create"`
}
```
- Add `chat_id` field with index
- Composite index on `(chat_id, user_id)` for efficient queries

**File: `pkg/repo/messagestats.go`**
- Add `chatID string` parameter to all methods:
  - `InsertMessageStats(chatID string, userID int64, words int) error`
  - `GetKnownUserIDs(chatID string) ([]int64, error)`
  - `GetWordCounts(chatID string) ([]MessageStatsWordCountStruct, error)`
- Update queries to filter by `chat_id`

---

### 1.3 Database Migration Strategy

**File: `pkg/db/migration.go`**

Create a comprehensive migration system that handles:

1. **Schema migrations** (add `chat_id` columns)
2. **Data migrations** (populate `chat_id` for existing records)
3. **Index migrations** (update unique constraints and indexes)

**Migration Steps:**

```go
type Migration struct {
    version     int
    description string
    up          func(*gorm.DB, string) error
    down        func(*gorm.DB) error
}
```

**Migration 001: Add ChatID to all tables**
- Add `chat_id VARCHAR(255)` to `plusplus`, `stats`, `message_stats`
- Create temporary indexes

**Migration 002: Populate ChatID for existing data**
- Get default/first configured chat ID from config
- Update all existing records with this default chat ID:
  ```sql
  UPDATE plusplus SET chat_id = ? WHERE chat_id IS NULL OR chat_id = ''
  UPDATE stats SET chat_id = ? WHERE chat_id IS NULL OR chat_id = ''
  UPDATE message_stats SET chat_id = ? WHERE chat_id IS NULL OR chat_id = ''
  ```
- Log warning about data migration

**Migration 003: Update constraints and indexes**
- Drop old unique index on `plusplus.name`
- Create composite unique index on `(chat_id, name)`
- Drop old unique index on `stats.user_id`
- Create composite unique index on `(chat_id, user_id)`
- Create index on `message_stats.chat_id`

**Implementation:**
- Add `RunMigrations(db *gorm.DB, cfg *interfaces.ConfigStruct) error` function
- Track migration version in new table `schema_migrations`
- Only run migrations that haven't been applied yet
- Make migrations idempotent where possible

**Rollback Strategy:**
- Keep old column temporarily during migration
- Add config flag `database.dangerousAllowRollback`
- Implement down migrations for each step
- Document rollback procedure

---

## Phase 2: Matcher Registry and Per-Chat Configuration

### 2.1 Update bot-matcher Dependency

**Option A: Use existing functionality**
- Current `matcher.Config` already supports `enabled *bool`
- Extend to support per-chat configuration

**Option B: Request enhancement to bot-matcher**
- Add `ChatID string` field to matcher interface methods
- Update `Registry.Process()` to accept and pass chat ID
- This would require updating bot-matcher to v0.2.0+

**Recommended: Option A (extend current functionality)**
- Less invasive, works with current bot-matcher
- Handle per-chat logic in the bot layer

---

### 2.2 Matcher Configuration System

**New Structure:**

**File: `interfaces/config.go`**
```go
type MatcherConfigStruct struct {
    Enabled         *bool
    EnabledForChats []string  // Empty = all chats, otherwise whitelist
    DisabledForChats []string  // Blacklist, takes precedence over whitelist
}

type ChatMatcherConfigStruct struct {
    ChatID   string
    Matchers map[string]MatcherConfigStruct  // matcher identifier -> config
}
```

**Configuration Format:**

**File: `config.yaml`**
```yaml
chats:
  - chatID: "-1001234567890"
    name: "main-chat"
    matchers:
      plusplus:
        enabled: true
      fortune:
        enabled: true
      buzzwords:
        enabled: true
      goodmorning:
        enabled: false
  - chatID: "-1009876543210"
    name: "test-chat"
    matchers:
      plusplus:
        enabled: true
      fortune:
        enabled: true
      buzzwords:
        enabled: false
```

**Alternative simpler approach:**
```yaml
chats:
  - chatID: "-1001234567890"
    name: "main-chat"
    enabledMatchers:
      - plusplus
      - fortune
      - buzzwords
      - atall
      - goodmorning
  - chatID: "-1009876543210"
    name: "test-chat"
    enabledMatchers:
      - plusplus
      - fortune
```

---

### 2.3 Matcher Registry Changes

**File: `container/container.go`**

Current implementation creates a single global `matcherRegistryInstance`. This needs to change to support per-chat matcher configuration.

**Option A: One Registry per Chat**
```go
var matcherRegistries = make(map[string]*matcher.Registry)  // chatID -> registry

func ProvideMatcherRegistry(chatID string) *matcher.Registry {
    if registry, exists := matcherRegistries[chatID]; exists {
        return registry
    }

    registry := createRegistryForChat(chatID)
    matcherRegistries[chatID] = registry
    return registry
}

func createRegistryForChat(chatID string) *matcher.Registry {
    registry := matcher.NewRegistry(ProvideLogger(), ProvideTelegramClient())

    cfg := GetChatConfig(chatID)

    // Register matchers based on chat-specific config
    if cfg.IsMatcherEnabled("atall") {
        registry.Register(atall.MakeMatcher(ProvideUserStatsRepo()))
    }
    // ... register other matchers

    return registry
}
```

**Option B: Single Registry, Filter at Process Time**
```go
func ProvideMatchersRegistry() *matcher.Registry {
    // Register ALL matchers
    // Filter based on chat at process time
}

// In webhook handler, check if matcher enabled for this chat before processing
```

**Recommended: Option A** - Cleaner separation, better performance

---

### 2.4 Webhook Handler Changes

**File: `container/container.go`**

Update `ProvideTelegramWebhookHandler` to:
1. Extract chat ID from incoming message
2. Get correct matcher registry for that chat
3. Pass chat ID through to matchers and state service

```go
func ProvideTelegramWebhookHandler() telegramclient.WebhookHandlerInterface {
    return telegramclient.NewHandler(
        &ProvideConfig().Telegram,
        func(messageIn telegramclient.WebhookMessageStruct) {
            chatID := strconv.FormatInt(messageIn.Chat.ID, 10)

            // Validate this is a configured chat
            if !IsConfiguredChat(chatID) {
                ProvideLogger().Warnf("Ignoring message from unconfigured chat: %s", chatID)
                return
            }

            // Get chat-specific registry
            registry := ProvideMatcherRegistry(chatID)
            registry.Process(messageIn)

            // Process state with chat ID
            stateService := ProvideState()
            stateService.ProcessMessage(chatID, messageIn)
        },
    )
}
```

---

## Phase 3: Per-Chat Matcher-Specific Configurations

### 3.1 Buzzwords Matcher

**Current:** Single global config at `config/buzzwords.yaml`

**New Structure:**
```
config/
  buzzwords/
    main-chat.yaml
    test-chat.yaml
    default.yaml
```

**File: `pkg/matchers/buzzwords/buzzwords.go`**

```go
func MakeMatcher(
    repo interfaces.PlusplusRepoInterface,
    chatID string,  // NEW parameter
) Matcher {
    var cfg Config

    // Try chat-specific config first, fall back to default
    configPath := fmt.Sprintf("config/buzzwords/%s.yaml", getChatName(chatID))
    if !fileExists(configPath) {
        configPath = "config/buzzwords/default.yaml"
    }

    matcher.LoadMatcherConfigFromPath(configPath, &cfg)

    // Rest of initialization...
}
```

**Update in container.go:**
```go
func createRegistryForChat(chatID string) *matcher.Registry {
    // ...
    if cfg.IsMatcherEnabled("buzzwords") {
        registry.Register(buzzwords.MakeMatcher(
            ProvidePlusplusRepo(),
            chatID,  // Pass chat ID
        ))
    }
}
```

**Migration:**
- Move `config/buzzwords.yaml` to `config/buzzwords/default.yaml`
- Support both old and new paths temporarily
- Log deprecation warning

---

### 3.2 Fortune Matcher

**Current:** Single global fortune files directory `files/fortune/`

**New Structure:**
```
files/
  fortune/
    main-chat/
      wisdom.txt
      jokes.txt
    test-chat/
      testing.txt
    default/
      wisdom.txt
```

**File: `pkg/fortune/service.go`**

```go
type Service struct {
    basePath string
}

func MakeService(chatID string) Service {
    chatName := getChatName(chatID)
    basePath := fmt.Sprintf("files/fortune/%s", chatName)

    // Fall back to default if chat-specific directory doesn't exist
    if !dirExists(basePath) {
        basePath = "files/fortune/default"
    }

    return Service{basePath: basePath}
}

func (f Service) GetList() []string {
    // Use f.basePath instead of hardcoded "files/fortune"
}

// Update all methods to use f.basePath
```

**File: `interfaces/fortune.go`**
```go
type FortuneServiceInterface interface {
    GetList() []string
    Exists(file string) bool
    GetRandomFortune() (FortuneInterface, error)
    GetFortune(file string) (FortuneInterface, error)
}

// Add factory function
type FortuneServiceFactory interface {
    CreateForChat(chatID string) FortuneServiceInterface
}
```

**File: `pkg/fortune/factory.go`** (new)
```go
type ServiceFactory struct{}

func MakeServiceFactory() ServiceFactory {
    return ServiceFactory{}
}

func (f ServiceFactory) CreateForChat(chatID string) interfaces.FortuneServiceInterface {
    return MakeService(chatID)
}
```

**Update matchers:**

**File: `pkg/matchers/fortune/fortune.go`**
```go
func MakeMatcher(
    fortuneService interfaces.FortuneServiceInterface,
) Matcher {
    // No changes needed - service is already injected
}
```

**File: `pkg/matchers/goodmorning/goodmorning.go`**
```go
func MakeMatcher(
    state interfaces.StateServiceInterface,
    fortuneService interfaces.FortuneServiceInterface,
) Matcher {
    // No changes needed - service is already injected
}
```

**Update container:**
```go
var fortuneServices = make(map[string]fortune.Service)

func ProvideFortuneService(chatID string) fortune.Service {
    if service, exists := fortuneServices[chatID]; exists {
        return service
    }

    service := fortune.MakeService(chatID)
    fortuneServices[chatID] = service
    return service
}

func createRegistryForChat(chatID string) *matcher.Registry {
    // ...
    fortuneService := ProvideFortuneService(chatID)

    if cfg.IsMatcherEnabled("fortune") {
        registry.Register(fortune2.MakeMatcher(fortuneService))
    }
    if cfg.IsMatcherEnabled("goodmorning") {
        registry.Register(goodmorning.MakeMatcher(ProvideState(), fortuneService))
    }
}
```

**Migration:**
- Move existing `files/fortune/*.txt` to `files/fortune/default/`
- Create chat-specific directories as needed
- Document file organization

---

## Phase 4: State Service and Repository Updates

### 4.1 State Service Changes

**File: `pkg/state/service.go`**

The state service tracks last post times and processes message stats. It needs to be chat-aware.

```go
type Service struct {
    log              logger.Interface
    userStatsRepo    interfaces.UserStatsRepoInterface
    messageStatsRepo interfaces.MessageStatsRepoInterface
    lastPost         map[string]map[int64]time.Time  // chatID -> userID -> time
}

func NewService(
    userStatsRepo interfaces.UserStatsRepoInterface,
    messageStatsRepo interfaces.MessageStatsRepoInterface,
) *Service {
    state := &Service{
        log:              logger.New(),
        userStatsRepo:    userStatsRepo,
        messageStatsRepo: messageStatsRepo,
        lastPost:         make(map[string]map[int64]time.Time),
    }
    return state
}

func (s *Service) Init(chatID string) {
    users, err := s.userStatsRepo.GetKnownUsers(chatID)
    if err != nil {
        s.log.Error("Error while getting known users from DB:", err)
        return
    }

    if s.lastPost[chatID] == nil {
        s.lastPost[chatID] = make(map[int64]time.Time)
    }

    for _, user := range users {
        s.lastPost[chatID][user.ID] = user.LastPost
    }
}

func (s *Service) ProcessMessage(chatID string, messageIn telegramclient.WebhookMessageStruct) {
    s.updateUserStats(chatID, messageIn)
    s.updateMessageStats(chatID, messageIn)
}

func (s *Service) GetLastPost(chatID string, userID int64) *time.Time {
    getLastPostLock.Lock()
    defer getLastPostLock.Unlock()

    if chatPosts, ok := s.lastPost[chatID]; ok {
        if lastPost, ok := chatPosts[userID]; ok {
            return &lastPost
        }
    }

    return nil
}

func (s *Service) updateUserStats(chatID string, messageIn telegramclient.WebhookMessageStruct) {
    if s.lastPost[chatID] == nil {
        s.lastPost[chatID] = make(map[int64]time.Time)
    }
    s.lastPost[chatID][messageIn.From.ID] = time.Now()

    if err := s.userStatsRepo.UpdateStats(
        chatID,
        messageIn.From.ID,
        messageIn.From.UsernameOrName(),
    ); err != nil {
        s.log.Error("Error while updating user stats in DB:", err)
    }
}

func (s *Service) updateMessageStats(chatID string, messageIn telegramclient.WebhookMessageStruct) {
    if err := s.messageStatsRepo.InsertMessageStats(
        chatID,
        messageIn.From.ID,
        messageIn.WordCount(),
    ); err != nil {
        s.log.Error("Error while inserting message stats in DB:", err)
    }
}
```

**File: `interfaces/state.go`**
```go
type StateServiceInterface interface {
    Init(chatID string)
    ProcessMessage(chatID string, messageIn telegramclient.WebhookMessageStruct)
    GetLastPost(chatID string, userID int64) *time.Time
}
```

**Update container:**
```go
func ProvideState() interfaces.StateServiceInterface {
    stateLock.Lock()
    defer stateLock.Unlock()

    if stateInstance == nil {
        stateInstance = state.NewService(
            ProvideUserStatsRepo(),
            ProvideMessageStatsRepo(),
        )

        // Initialize for all configured chats
        for _, chat := range ProvideConfig().Chats {
            stateInstance.Init(chat.ChatID)
        }
    }

    return stateInstance
}
```

---

### 4.2 Update All Matchers to Use ChatID

All matchers that interact with repositories or state need to be updated to pass the chat ID.

#### 4.2.1 Matchers with Database Operations

**plusplus matcher:**
```go
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
    chatID := strconv.FormatInt(messageIn.Chat.ID, 10)
    // ... existing logic
    value, err := m.repo.Increment(chatID, token.Name, token.Increment)
    // ... rest
}
```

**buzzwords matcher:**
```go
func (m Matcher) makeRepliesFromTrigger(chatID string, trigger string) ([]telegramclient.MessageStruct, error) {
    value, err := m.repo.Increment(chatID, trigger, 1)
    // ... rest
}
```

**topflop matcher:**
```go
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
    chatID := strconv.FormatInt(messageIn.Chat.ID, 10)
    // ... existing logic
    switch cmd {
    case "top":
        records, err = m.repo.FindTops(chatID, limit)
    case "flop":
        records, err = m.repo.FindFlops(chatID, limit)
    }
    // ... rest
}
```

**stats matcher:**
```go
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
    chatID := strconv.FormatInt(messageIn.Chat.ID, 10)
    users, err := m.repo.GetTopUsers(chatID)
    // ... rest
}
```

**atall matcher:**
```go
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
    chatID := strconv.FormatInt(messageIn.Chat.ID, 10)
    // ... existing logic
    users, err := m.repo.GetKnownUsers(chatID)
    // ... rest
}
```

**goodmorning matcher:**
```go
func (m Matcher) doesMatch(messageIn telegramclient.WebhookMessageStruct) bool {
    chatID := strconv.FormatInt(messageIn.Chat.ID, 10)
    now := time.Now()

    if now.Hour() < 6 || now.Hour() > 14 {
        return false
    }

    lastPost := m.state.GetLastPost(chatID, messageIn.From.ID)
    // ... rest
}
```

#### 4.2.2 Matchers Without Database (no changes needed)

These matchers don't use database or per-chat config, so they work across all chats without modification:
- **choose** - Random selection from options
- **janein** - Yes/no randomizer
- **ping** - Simple ping response
- **xkcd** - Fetches from external API

---

## Phase 5: Testing and Validation

### 5.1 Unit Tests

Update all existing unit tests to include chat ID:

**Example: `pkg/matchers/plusplus/plusplus_test.go`**
```go
func TestPlusplus(t *testing.T) {
    chatID := "test-chat"
    mockRepo := &mocks.PlusplusRepoInterface{}
    mockRepo.On("Increment", chatID, "foo", 1).Return(5, nil)
    // ... rest of test
}
```

### 5.2 Integration Tests

Create integration tests for:
1. Multiple chats with same user
2. Data isolation between chats
3. Per-chat matcher configuration
4. Migration from single-chat to multi-chat
5. Config validation

### 5.3 Manual Testing Checklist

- [ ] Configure two different chats in config
- [ ] Send messages to both chats
- [ ] Verify plusplus values are separate
- [ ] Verify user stats are separate
- [ ] Verify fortune files load correctly per chat
- [ ] Verify buzzwords configs load correctly per chat
- [ ] Test matcher enable/disable per chat
- [ ] Test migration with existing database
- [ ] Test rollback scenario
- [ ] Test config validation (missing chatID, etc.)

---

## Phase 6: Documentation and Deployment

### 6.1 Documentation Updates

**README.md:**
- Add multi-chat configuration section
- Explain chat naming conventions
- Show examples of per-chat configurations

**Migration Guide:**
- Step-by-step upgrade from single chat to multi-chat
- Database backup recommendations
- Rollback procedure if needed

**Configuration Reference:**
- Document all new config options
- Show all supported matcher configurations
- Explain file organization for per-chat files

### 6.2 Configuration Templates

Update `.env.dist`:
```bash
# Comma-separated list of chat IDs (or use config.yaml for more control)
TELEGRAM_CHAT_IDS="-1001234567890,-1009876543210"
```

Create `config.yaml.dist`:
```yaml
chats:
  - chatID: "-1001234567890"
    name: "main-chat"
    description: "Main team chat"
    enabledMatchers:
      - plusplus
      - fortune
      - buzzwords
      - atall
      - goodmorning
      - stats
      - topflop
      - xkcd
      - choose
      - janein
      - ping

  - chatID: "-1009876543210"
    name: "test-chat"
    description: "Testing environment"
    enabledMatchers:
      - plusplus
      - fortune
      - xkcd
      - choose
      - janein
      - ping
```

Create example configs:
- `config/buzzwords/default.yaml`
- `config/buzzwords/example-chat.yaml`

### 6.3 Deployment Strategy

**Step 1: Preparation**
1. Backup existing database
2. Review and update configuration files
3. Create per-chat configuration files (buzzwords, fortune)
4. Test in staging environment

**Step 2: Deployment**
1. Stop bot service
2. Deploy new version
3. Run database migrations
4. Validate migration success
5. Start bot service
6. Monitor logs for errors

**Step 3: Validation**
1. Send test messages to each configured chat
2. Verify data separation
3. Check matcher functionality
4. Monitor error logs

**Rollback Plan:**
1. Stop bot service
2. Restore database from backup
3. Deploy previous version
4. Start bot service

---

## Summary of Changes by File

### New Files
- `pkg/db/migrations.go` - Database migration system
- `pkg/fortune/factory.go` - Fortune service factory
- `config/buzzwords/default.yaml` - Default buzzwords config
- `files/fortune/default/` - Default fortune files directory
- `MULTI_CHAT_IMPLEMENTATION_PLAN.md` - This document

### Modified Files

#### Configuration
- `interfaces/config.go` - Add ChatConfigStruct, update ConfigStruct
- `pkg/config/config.go` - Support multiple chats
- `config.yaml` - Add chats array
- `.env.dist` - Document multi-chat configuration

#### Database
- `interfaces/repoplusplus.go` - Add ChatID to Plusplus model
- `interfaces/repostats.go` - Add ChatID to Stats model
- `interfaces/repomessagestats.go` - Add ChatID to MessageStats model
- `pkg/repo/plusplus.go` - Add chatID parameter to all methods
- `pkg/repo/userstats.go` - Add chatID parameter to all methods
- `pkg/repo/messagestats.go` - Add chatID parameter to all methods
- `pkg/db/migration.go` - Implement migration system

#### Core Bot Infrastructure
- `container/container.go` - Per-chat registries, chat-aware services
- `cmd/bot.go` - May need chat initialization loop

#### State Management
- `interfaces/state.go` - Add chatID parameters
- `pkg/state/service.go` - Support per-chat state tracking

#### Fortune Service
- `interfaces/fortune.go` - Add factory interface
- `pkg/fortune/service.go` - Support per-chat file paths

#### Matchers with Changes
- `pkg/matchers/plusplus/plusplus.go` - Pass chatID to repo
- `pkg/matchers/buzzwords/buzzwords.go` - Per-chat config, pass chatID
- `pkg/matchers/topflop/topflop.go` - Pass chatID to repo
- `pkg/matchers/stats/stats.go` - Pass chatID to repo
- `pkg/matchers/atall/atall.go` - Pass chatID to repo
- `pkg/matchers/goodmorning/goodmorning.go` - Use chat-aware state
- `pkg/matchers/fortune/fortune.go` - Use chat-specific service

#### Tests
- All `*_test.go` files - Update to include chatID

---

## Estimated Implementation Effort

### By Phase
1. **Phase 1 (Core Infrastructure)**: 8-12 hours
   - Config changes: 2 hours
   - Database schema: 3 hours
   - Migration system: 5 hours

2. **Phase 2 (Matcher Registry)**: 6-8 hours
   - Registry changes: 3 hours
   - Webhook handler: 2 hours
   - Per-chat config loading: 3 hours

3. **Phase 3 (Matcher-Specific Configs)**: 4-6 hours
   - Buzzwords: 2 hours
   - Fortune: 3 hours

4. **Phase 4 (State and Matchers)**: 6-8 hours
   - State service: 2 hours
   - Update all matchers: 4 hours

5. **Phase 5 (Testing)**: 6-8 hours
   - Unit tests: 3 hours
   - Integration tests: 3 hours
   - Manual testing: 2 hours

6. **Phase 6 (Documentation)**: 3-4 hours
   - Documentation: 2 hours
   - Templates and examples: 1 hour

**Total Estimated Time**: 33-46 hours

### Risk Factors
- **Database migration complexity**: Could take longer if existing data is large
- **Bot-matcher compatibility**: May need to update dependency or fork
- **Testing edge cases**: Multi-chat scenarios can be complex
- **Configuration validation**: Need robust error handling

### Recommended Order
1. Start with Phase 1 (infrastructure) - foundational
2. Then Phase 2 (registry) - core functionality
3. Then Phase 4 (state) - needed for most matchers
4. Then Phase 3 (configs) - nice to have
5. Then Phase 5 (testing) - validation
6. Finally Phase 6 (docs) - polish

---

## Additional Considerations

### Security
- Validate incoming chat IDs against configured list
- Log unauthorized access attempts
- Consider rate limiting per chat

### Performance
- Per-chat registries cache matchers (good)
- Database queries filtered by chat_id with indexes (good)
- State service tracks all chats in memory (acceptable for reasonable chat count)
- Consider periodic cleanup of inactive chat data

### Monitoring
- Add metrics per chat ID
- Track message processing time per chat
- Monitor database query performance with new indexes
- Alert on errors from specific chats

### Future Enhancements
- Admin commands to add/remove chats without restart
- Web UI for per-chat configuration
- Import/export chat configurations
- Chat templates (copy settings from one chat to another)
- Per-chat rate limiting
- Per-chat logging levels
- Shared global configurations with chat overrides

---

## Questions and Open Issues

1. **bot-matcher dependency**: Should we update it to natively support multi-chat, or handle at bot layer?
   - **Decision**: Handle at bot layer for now, may propose upstream changes later

2. **Configuration format**: YAML vs environment variables vs database?
   - **Decision**: YAML for per-chat config, env vars for backwards compatibility

3. **Chat identification**: Use chat ID string or integer?
   - **Decision**: Use string for flexibility (Telegram IDs can be large)

4. **Default behavior**: What happens if a chat isn't configured?
   - **Decision**: Ignore messages, log warning

5. **Matcher naming**: Should we use chat name or chat ID in configs?
   - **Decision**: Use chat name in file paths, chat ID for runtime identification

---

## Conclusion

This plan provides a comprehensive roadmap for implementing multi-chat support. The implementation is designed to:

- Maintain backwards compatibility where possible
- Provide clear migration path from single to multi-chat
- Ensure complete data isolation between chats
- Support flexible per-chat configuration
- Scale to support many chats without performance degradation

The phased approach allows for incremental implementation and testing, reducing risk of introducing bugs or breaking existing functionality.
