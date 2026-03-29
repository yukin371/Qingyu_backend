# models refactor plan

## 1. 结构性重构

系统应该划分为 **“公共资源（Bookstore）”** 和 **“用户私有状态（Reader）”** 以进行解耦。

目前存在一些**边界模糊**和**归属不当**的模型，如 `Chapter`、`Comment` 和 `Like`。

### A. `Chapter`（章节）的归属
`Chapter` 同时出现在 `reader` 代码和 `bookstore` 列表里，这是一个**冲突点**。
*   **分析**：章节是书籍的一部分，是**客观存在的内容**。无论用户是否登录，章节的标题、字数、价格（元数据）都是固定的。虽然VIP章节的内容需要权限才能看，但数据本身属于“书”这个实体。
*   **结论**：`Chapter` 应该属于 **Bookstore**（或独立的 Content 域）。`Reader` 模块应该只引用 `ChapterID` 来记录进度，而不应该包含章节的定义。

### B. `Comment`（评论）与 `Like`（点赞）的归属
这两者属于**用户生成内容 (UGC)**，具有“公共可见”的属性。
*   **分析**：虽然需要登录才能评论，但评论是给所有人看的。如果放在 `reader` 包中，会让人感觉这是用户的“私有数据”（像笔记一样）。
*   **结论**：这两者通常独立为一个 **Community (社区/互动)** 模块（目前搁置了社交模块），对于当前系统，作为书籍的**附属数据**归入 **Bookstore** 更合适。


### 2. 建议的详细划分方案

建议将模型重构为三个逻辑域（可以是三个包，也可以是微服务拆分）：

#### 域 1: Bookstore (书城/内容域)
**性质**：公共读、低频写（后台发书）、无需登录即可检索。
**包含模型**：
*   `Banner` (轮播图)
*   `Category` (分类)
*   `Book` (书籍基本信息)
*   `BookDetail` (书籍详情)
*   `BookStatistics` (统计：字数、总点击等)
*   `Ranking` (榜单)
*   **`Chapter`** (迁移至此)：包含标题、内容（或内容的OSS链接）、价格、字数。

#### 域 2: Reader (阅读器/用户状态域)
**性质**：私有读写、高频写（记录进度）、强依赖用户ID。
**包含模型**：
*   `ReadingProgress` (当前阅读进度)：记录读到了哪一章、哪一页。
*   `ReadingHistory` (阅读足迹/流水)：记录什么时间读了什么书（用于“最近阅读”列表）。
*   `Annotation` (标注)：笔记、高亮、划线（这是完全私有的）。
*   `Collection/Folder` (书架/收藏)：用户的私人书库。
*   `ReadingSettings` (阅读设置)：字体、背景、翻页模式。

#### 域 3: Community (互动域 - 可选，或并入 Bookstore)
**性质**：用户生产、公共展示。
**包含模型**：
*   **`Comment`** (迁移至此)：书评、章评。
*   **`Like`** (迁移至此)：点赞记录。
*   `BookRating` (评分)。

---

### 3. 代码层面的调整建议

基于上述划分，对 `reader` 包代码进行具体的“留存”或“剔除”建议：

#### ❌ 建议移出 `reader` 包的模型：

1.  **`Chapter`**
    *   **原因**：这是静态资源，不是用户状态。
    *   **处理**：移到 `bookstore`。
    *   **注意**：`reader` 中的业务逻辑（如 `Annotation`）只需存储 `ChapterID` 即可，不需要嵌入整个 `Chapter` 结构体。

2.  **`Comment` & `Like`**
    *   **原因**：它们属于社交属性。虽然是用户产生的，但并不属于“阅读工具”的状态。
    *   **处理**：建议新建 `package community` 或放入 `bookstore`。

#### ✅ 保留在 `reader` 包并优化的模型：

1.  **`Annotation` (标注)**
    *   **状态**：保留。
    *   **建议**：设计很好。`Range` 字段如果对应前端的 DOM Range 或定位符，需要确保书籍内容变更（如勘误）时，定位不会失效。

2.  **`Collection` (书架)**
    *   **状态**：保留。
    *   **建议**：`FolderID` 的设计很好，支持分组。

3.  **`ReadingHistory` vs `ReadingProgress`**
    *   **状态**：保留。
    *   **区分建议**：
        *   **`ReadingProgress`**：**快照数据**。一张表，每个 User+Book 只有一行记录。用于下次打开书时跳转到上次位置。
        *   **`ReadingHistory`**：**日志数据**。每次阅读行为产生一条记录。用于分析用户阅读时长、活跃时间段，或者展示“阅读轨迹”。
    *   **优化**：你的 `ReadingHistory` 中有 `Progress` 字段，这与 `ReadingProgress` 有点重叠。如果不需要精确回溯“上周五下午3点读到了第几页”，History 中可以去掉 Progress，只留时长和章节。

4.  **`ReadingSettings`**
    *   **状态**：保留。
    *   **建议**：通常每个用户只有一条记录，或者分端（Web/App）存储。

### 4. 总结架构图

```mermaid
graph TD
    User((用户))

    subgraph Bookstore_Package [书城 (公共数据)]
        B_Book[Book 书籍]
        B_Chapter[Chapter 章节内容]
        B_Cat[Category 分类]
    end

    subgraph Community_Package [社区 (互动数据)]
        C_Comment[Comment 评论]
        C_Like[Like 点赞]
        C_Rate[Rating 评分]
    end

    subgraph Reader_Package [阅读器 (私有状态)]
        R_Prog[Progress 进度]
        R_Shelf[Collection 书架]
        R_Anno[Annotation 笔记/高亮]
        R_Set[Settings 设置]
    end

    User -- 浏览 --> Bookstore_Package
    User -- 发布 --> Community_Package
    User -- 记录 --> Reader_Package
    
    R_Prog -. Ref ID .-> B_Chapter
    R_Shelf -. Ref ID .-> B_Book
    C_Comment -. Ref ID .-> B_Book
```

### 结论

建议**将 `Chapter`、`Comment`、`Like` 从 `reader` 包中移出**。让 `reader` 包保持纯粹的“用户个人阅读状态管理”。这样当你的系统扩展时（比如增加听书功能），`reader` 的逻辑不会被内容数据的变化所拖累。
