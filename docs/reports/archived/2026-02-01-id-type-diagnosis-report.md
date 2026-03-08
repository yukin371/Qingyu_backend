# ID 类型诊断报告

**生成时间**: 2026-02-01T09:02:29.237Z

## 摘要

| 指标 | 数值 |
|------|------|
| 总集合数 | 45 |
| ObjectID 类型集合 | 16 |
| String 类型集合 | 19 |
| 发现问题数 | 5 |

## 集合 _id 类型清单

| 集合 | 文档数 | _id 类型 | 示例值 |
|------|--------|----------|--------|
| collection_folders | 375 | ObjectID | `68fca4a04ba680177cc4505f` |
| book_list_likes | 273 | ObjectID | `697b529f73dedbe4807bf613` |
| novel_files | 1 | String | `68fb75f0a9679303239229a2` |
| collections | 643 | ObjectID | `68fe3404425419c8fc8f7afc` |
| users | 55 | String | `admin` |
| chapter_stats | 7119 | String | `697ed319cbc5d964a9e8deee` |
| ai_user_quotas | 1 | ObjectID | `68fc54534d1544e1b81d77f0` |
| notification_templates | 14 | ObjectID | `697e0641b7cf4a68a0d806b4` |
| comments | 3 | ObjectID | `697e18888df88e315da6ef77` |
| messages | 200 | String | `697ed3174467515bd0d5a477` |
| follows | 0 | N/A | `N/A` |
| document_contents | 6 | ObjectID | `697e1cef6ac9aea960bf1216` |
| book_stats | 30194 | String | `69759c7d0c0a8f5c9486c8fe` |
| reading_progress | 2737 | String | `d0fa5dde-96cc-4ac9-8b26-ee5f59b2bee5` |
| likes | 0 | N/A | `N/A` |
| user_collections | 0 | N/A | `N/A` |
| files | 983 | String | `68ff734a926c0e207cc781b0` |
| memberships | 15 | String | `697ed3191e9d298713b3627b` |
| ai_quotas | 20 | ObjectID | `69523411d9cab999fbd6022f` |
| reading_histories | 0 | N/A | `N/A` |
| book_lists | 0 | N/A | `N/A` |
| categories | 6 | ObjectID | `68fcbdaed8ef89fcdadd27fc` |
| banners | 2 | ObjectID | `697f11cad32074cfe3aa3f6c` |
| file_patches | 1 | String | `68fb75f0a9679303239229a4` |
| rankings | 310 | ObjectID | `697f11cad32074cfe3aa3f6e` |
| chapters | 305641 | String | `697e0641b7cf4a68a0d806d3` |
| chapter_purchases | 0 | N/A | `N/A` |
| file_revisions | 0 | N/A | `N/A` |
| reading_history | 1173 | String | `e816d9c2-466a-4cb1-ac4c-e8c706357ad3` |
| notifications | 200 | String | `697ed316cf5b5e0efb19897c` |
| ranking_items | 0 | N/A | `N/A` |
| roles | 0 | N/A | `N/A` |
| projects | 7 | ObjectID | `697e18888df88e315da6ef99` |
| documents | 6 | ObjectID | `697e1cef6ac9aea960bf1214` |
| batch_operations | 1 | ObjectID | `697cdf482b2b715242d00afb` |
| annotations | 55 | String | `58c42b1d-bf25-40db-80e9-2eb869bf1cae` |
| file_access | 62 | ObjectID | `69759c7dc9ffe8c47e3c86be` |
| multipart_uploads | 136 | String | `68ff734b926c0e207cc781b9` |
| chapter_contents | 725420 | ObjectID | `695255ddae8cbb7803ad31dc` |
| author_revenue | 1201 | String | `697ed3191e9d298713b35e00` |
| books | 100 | String | `697f11cad32074cfe3aa3ea5` |
| author_follows | 0 | N/A | `N/A` |
| conversations | 16 | String | `697ed3174467515bd0d5a476` |
| bookmarks | 224 | String | `ef528b9f-0507-4d31-bc03-29646e8d9a31` |
| announcements | 3 | String | `697ed3174467515bd0d5a54e` |

## 外键字段类型清单

| 集合 | 字段 | 类型 | 示例值 |
|------|------|------|--------|
| collection_folders | user_id | string | `68fc402bcd736a40d4220ba6` |
| book_list_likes | booklist_id | string | `697b529f73dedbe4807bf611` |
| book_list_likes | user_id | string | `697b529f73dedbe4807bf612` |
| novel_files | project_id | string | `test_project` |
| novel_files | node_id | string | `node_patch_1` |
| collections | user_id | string | `68fc402bcd736a40d4220ba6` |
| collections | book_id | string | `507f1f77bcf86cd799439011` |
| chapter_stats | chapter_id | string | `697e0641b7cf4a68a0d806d3` |
| ai_user_quotas | user_id | string | `68fc402acd736a40d4220ba4` |
| comments | author_id | string | `697e18878df88e315da6ef72` |
| comments | target_id | string | `697e18888df88e315da6ef74` |
| messages | conversation_id | string | `697ed3174467515bd0d5a476` |
| messages | sender_id | string | `admin` |
| messages | receiver_id | string | `author1` |
| document_contents | document_id | ObjectID | `697e1cef6ac9aea960bf1214` |
| book_stats | book_id | string | `book_69759c7d0c0a8f5c9486c8fc` |
| book_stats | author_id | string | `author_69759c7d0c0a8f5c9486c8fd` |
| reading_progress | user_id | string | `admin` |
| reading_progress | book_id | string | `697eca76c2ff09ef0e9ad98a` |
| files | user_id | string | `user_123` |
| memberships | user_id | string | `admin` |
| ai_quotas | user_id | string | `6911b83eac0bff9fbf62604d` |
| file_patches | project_id | string | `test_project` |
| file_patches | node_id | string | `node_patch_1` |
| rankings | book_id | ObjectID | `697f11cad32074cfe3aa3f0d` |
| chapters | book_id | string | `697e0641b7cf4a68a0d806d2` |
| reading_history | user_id | string | `25e1e57e-1ffd-4522-98c5-398c25200b05` |
| reading_history | book_id | string | `697eca76c2ff09ef0e9ad9fc` |
| reading_history | chapter_id | string | `697e067034a47b8fb262dd95` |
| notifications | user_id | string | `admin` |
| projects | author_id | string | `697e18888df88e315da6ef98` |
| documents | project_id | ObjectID | `697e1cef6ac9aea960bf1213` |
| batch_operations | project_id | ObjectID | `697cdf482b2b715242d00af4` |
| batch_operations | target_ids | ObjectID | `697cdf482b2b715242d00af5,697cdf482b2b715242d00af7` |
| annotations | user_id | string | `50db10c3-5519-4461-8b27-3e533f613f15` |
| annotations | book_id | string | `697eca76c2ff09ef0e9add12` |
| annotations | chapter_id | string | `697e189e64f82fef2eebab13` |
| file_access | file_id | string | `69759c7dc9ffe8c47e3c86bd` |
| file_access | user_id | string | `user_456` |
| multipart_uploads | upload_id | string | `68ff734b926c0e207cc781ba` |
| chapter_contents | chapter_id | ObjectID | `6911b840ac0bff9fbf626062` |
| author_revenue | user_id | string | `author1` |
| author_revenue | book_id | string | `697eca76c2ff09ef0e9ad98a` |
| bookmarks | user_id | string | `cf9dcb61-a7da-4449-8d92-ca049876efeb` |
| bookmarks | book_id | string | `697eca76c2ff09ef0e9adc80` |
| bookmarks | chapter_id | string | `697e066f34a47b8fb262dd93` |

## 发现的问题

### 问题 1: ID_TYPE_MISMATCH

| 属性 | 值 |
|------|-----|
| 严重程度 | HIGH |
| 描述 | 存在混合的 ID 类型 |
| 详情 | `{"objectIDCollections":["collection_folders","book_list_likes","collections","ai_user_quotas","notification_templates","comments","document_contents","ai_quotas","categories","banners","rankings","projects","documents","batch_operations","file_access","chapter_contents"],"stringCollections":["novel_files","users","chapter_stats","messages","book_stats","reading_progress","files","memberships","file_patches","chapters","reading_history","notifications","annotations","multipart_uploads","author_revenue","books","conversations","bookmarks","announcements"]}` |
| 建议 | 统一使用 ObjectID 类型，外键字段应存储为 string (ObjectID.Hex()) |

### 问题 2: FOREIGN_KEY_TYPE_MISMATCH

| 属性 | 值 |
|------|-----|
| 严重程度 | HIGH |
| 描述 | 外键字段 project_id 类型与目标集合 projects._id 类型不匹配 |
| 详情 | `{"collection":"novel_files","field":"project_id","currentType":"string","targetCollection":"projects","targetIdType":"ObjectID"}` |
| 建议 | 外键应存储为 string 类型 (ObjectID.Hex()) |

### 问题 3: FOREIGN_KEY_TYPE_MISMATCH

| 属性 | 值 |
|------|-----|
| 严重程度 | HIGH |
| 描述 | 外键字段 project_id 类型与目标集合 projects._id 类型不匹配 |
| 详情 | `{"collection":"file_patches","field":"project_id","currentType":"string","targetCollection":"projects","targetIdType":"ObjectID"}` |
| 建议 | 外键应存储为 string 类型 (ObjectID.Hex()) |

### 问题 4: FOREIGN_KEY_TYPE_MISMATCH

| 属性 | 值 |
|------|-----|
| 严重程度 | HIGH |
| 描述 | 外键字段 book_id 类型与目标集合 books._id 类型不匹配 |
| 详情 | `{"collection":"rankings","field":"book_id","currentType":"ObjectID","targetCollection":"books","targetIdType":"String"}` |
| 建议 | 外键应存储为 string 类型 (ObjectID.Hex()) |

### 问题 5: FOREIGN_KEY_TYPE_MISMATCH

| 属性 | 值 |
|------|-----|
| 严重程度 | HIGH |
| 描述 | 外键字段 chapter_id 类型与目标集合 chapters._id 类型不匹配 |
| 详情 | `{"collection":"chapter_contents","field":"chapter_id","currentType":"ObjectID","targetCollection":"chapters","targetIdType":"String"}` |
| 建议 | 外键应存储为 string 类型 (ObjectID.Hex()) |

## 修复建议

1. **统一 ID 类型**: 确保所有集合使用一致的 ID 类型
2. **修复外键类型**: 将所有外键字段转换为 string (ObjectID.Hex())
3. **更新关联查询**: 确保关联查询时正确处理 ID 类型转换
