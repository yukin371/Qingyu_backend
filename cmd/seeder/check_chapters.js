// MongoDB 查询脚本
var books = db.books.find({title: /修仙/}, {_id: 1, title: 1}).toArray();
print("书籍列表:");
books.forEach(book => {
  var chapterCount = db.chapters.countDocuments({book_id: book._id});
  print("  书籍: " + book.title + ", ID: " + book._id + ", 章节数: " + chapterCount);
  
  if (chapterCount > 0) {
    var chapters = db.chapters.find({book_id: book._id}).limit(1).toArray();
    print("    第一章节ID: " + chapters[0]._id);
  }
});
