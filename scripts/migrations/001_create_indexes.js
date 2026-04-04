// MongoDB Index Creation Script
// Version: 001
// Date: 2026-04-04
// Description: Create indexes for all collections to improve query performance

// ============== Writer Module ==============

// documents collection
db.documents.createIndex({ project_id: 1, parent_id: 1, order_key: 1 });
db.documents.createIndex({ project_id: 1, type: 1, updated_at: -1 });

// document_contents collection
db.document_contents.createIndex({ document_id: 1, version: -1 });

// versions collection
db.versions.createIndex({ document_id: 1, created_at: -1 });

// characters collection
db.characters.createIndex({ project_id: 1, name: 1 });

// character_relations collection
db.character_relations.createIndex({ project_id: 1, source_id: 1 });

// locations collection
db.locations.createIndex({ project_id: 1, name: 1 });

// outlines collection
db.outlines.createIndex({ project_id: 1, type: 1, order_key: 1 });

// timeline_events collection
db.timeline_events.createIndex({ project_id: 1, timestamp: 1 });

// ============== Bookstore Module ==============

// books collection
db.books.createIndex({ status: 1, published_at: -1 });
db.books.createIndex({ author_id: 1, status: 1 });
db.books.createIndex({ category: 1, status: 1 });
db.books.createIndex({ "stats.hot_score": -1 });
db.books.createIndex({ "stats.rating_avg": -1 });

// chapters collection
db.chapters.createIndex({ book_id: 1, order_key: 1 });
db.chapters.createIndex({ book_id: 1, status: 1 });

// book_ratings collection
db.book_ratings.createIndex({ book_id: 1, user_id: 1 }, { unique: true });

// comments collection
db.comments.createIndex({ target_id: 1, target_type: 1, created_at: -1 });

// ============== Reader Module ==============

// reading_progress collection
db.reading_progress.createIndex({ user_id: 1, book_id: 1 }, { unique: true });
db.reading_progress.createIndex({ user_id: 1, updated_at: -1 });

// bookmarks collection
db.bookmarks.createIndex({ user_id: 1, book_id: 1, chapter_id: 1 });

// ============== User Module ==============

// users collection
db.users.createIndex({ username: 1 }, { unique: true });
db.users.createIndex({ email: 1 }, { unique: true });
db.users.createIndex({ phone: 1 }, { sparse: true });

// sessions collection
db.sessions.createIndex({ user_id: 1, expires_at: 1 });
db.sessions.createIndex({ expires_at: 1 }, { expireAfterSeconds: 0 });

// ============== Notification Module ==============

// notifications collection
db.notifications.createIndex({ user_id: 1, is_read: 1, created_at: -1 });
db.notifications.createIndex({ user_id: 1, type: 1 });

// ============== Audit Module ==============

// audit_logs collection
db.audit_logs.createIndex({ user_id: 1, created_at: -1 });
db.audit_logs.createIndex({ action: 1, created_at: -1 });

// ============== Print Summary ==============
print("========================================");
print("Index creation completed successfully!");
print("========================================");
print("\nCreated indexes summary:");

// Writer
print("\n--- Writer Module ---");
print("documents: 2 indexes");
print("document_contents: 1 index");
print("versions: 1 index");
print("characters: 1 index");
print("character_relations: 1 index");
print("locations: 1 index");
print("outlines: 1 index");
print("timeline_events: 1 index");

// Bookstore
print("\n--- Bookstore Module ---");
print("books: 5 indexes");
print("chapters: 2 indexes");
print("book_ratings: 1 index (unique)");
print("comments: 1 index");

// Reader
print("\n--- Reader Module ---");
print("reading_progress: 2 indexes (1 unique)");
print("bookmarks: 1 index");

// User
print("\n--- User Module ---");
print("users: 3 indexes (1 unique, 1 sparse)");
print("sessions: 2 indexes (1 TTL)");

// Notification
print("\n--- Notification Module ---");
print("notifications: 2 indexes");

// Audit
print("\n--- Audit Module ---");
print("audit_logs: 2 indexes");

print("\nTotal: 27 indexes created");
print("========================================");
