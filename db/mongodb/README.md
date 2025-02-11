# 1. 简介
## 教程
https://www.runoob.com/mongodb/mongodb-tutorial.html

## 文档导向
MongoDB 存储 BSON（二进制 JSON）文档，这些文档可以包含复杂的数据结构，如数组和嵌套对象。
对于存储大于 BSON 文档大小限制（16MB）的文件，MongoDB 提供了 GridFS，一种用于存储和检索大文件的规范。

## 存储引擎
MongoDB支持多种存储引擎, 本文所有涉及mongo存储引擎的只谈默认的WiredTiger引擎, 其实还有某些方面更优秀的其他引擎,例如: MongoRocks等。

## 鸡肋的Collection 和 Type

早期为了跟传统rdbms数据库保持概念一致 ，mongodb和elasticsearch都设计了跟传统数据库里面的库->表->记录行对应的概念。

Elasticsearch从6.x版本开始强制只允许一个索引使用一个type, 其实就是意识到这个这个设计的失误, 不想让你用这个type类型, 因为type和传统数据库里面的表概念其实是不一样的，这种概念类比给人造成了误解，

es的7.x版本会默认取消type类型, 就说明这个type字段真的是鸡肋的不行。

```
RDBMS	MongoDB	Elasticsearch
库	库	索引
表	集合	类型
记录	文档	文档
```

## 弱事务
MongoDB以前只是支持同一文档内的原子更新, 以此来实现伪事务功能, 不过Mongo4.0支持Replica Set事务, 大大加强了事务方面的能力。

## 无join支持


## 使用场景
* 对服务可用性和一致性有高要求
* 无schema的数据存储 + 需要索引数据
* 高读写性能要求, 数据使用场景简单的海量数据场景
* 有热点数据, 有数据分片需求的数据存储
* 日志, html, 爬虫数据等半结构化或图片，视频等非结构化数据的存储
* 有js使用经验的人员(MongoDB内置操作语言为js)*

# 2. 工具
MongoDB提供了网络和系统监控工具Munin，它作为一个插件应用于MongoDB中。

Gangila是MongoDB高性能的系统监视的工具，它作为一个插件应用于MongoDB中。

基于图形界面的开源工具 Cacti, 用于查看CPU负载, 网络带宽利用率,它也提供了一个应用于监控 MongoDB 的插件。

GUI
* Fang of Mongo – 网页式,由Django和jQuery所构成。
* Futon4Mongo – 一个CouchDB Futon web的mongodb山寨版。
* Mongo3 – Ruby写成。
* MongoHub – 适用于OSX的应用程序。
* Opricot – 一个基于浏览器的MongoDB控制台, 由PHP撰写而成。
* Database Master — Windows的mongodb管理工具
* RockMongo — 最好的PHP语言的MongoDB管理工具，轻量级, 支持多国语言.

# 3. 安装
```
brew tap mongodb/brew
brew install mongodb-community
brew install mongodb-community-shell

brew services start mongodb-community
brew services stop mongodb-community
```

mongod 命令后台进程方式：
```
mongod --config /usr/local/etc/mongod.conf --fork
这种方式启动要关闭可以进入 mongo shell 控制台来实现：
> db.adminCommand({ "shutdown" : 1 })
```

配置文件：/usr/local/etc/mongod.conf
日志文件路径：/usr/local/var/log/mongodb
数据存放路径：/usr/local/var/mongodb


# 4. shell

## 4.1  常用命令
```
mongosh
mongosh --version
mongosh --host <hostname>:<port>
```

```
查看当前数据库：db
显示数据库列表：show dbs
切换到指定数据库：use <database_name>
执行查询操作：db.<collection_name>.find()
插入文档：db.<collection_name>.insertOne({ ... })
更新文档：db.<collection_name>.updateOne({ ... })
删除文档：db.<collection_name>.deleteOne({ ... })
退出 MongoDB Shell：quit() 或者 exit
```
## 4.2 用户管理
https://www.runoob.com/mongodb/mongodb-user.html

# 5. 数据库
* 数据库（Database）：存储数据的容器，类似于关系型数据库中的数据库。
* 集合（Collection）：数据库中的一个集合，类似于关系型数据库中的表。集合是一组文档的容器。在 MongoDB 中，一个集合中的文档不需要有一个固定的模式。
* 文档（Document）：集合中的一个数据记录，类似于关系型数据库中的行（row），以 BSON 格式存储。通常是一个 JSON-like 的结构，可以包含多种数据类型。
* 主键,MongoDB自动将_id字段设置为主键

数据库名可以是满足以下条件的任意UTF-8字符串:
* 不能是空字符串（"")。
* 不得含有' '（空格)、.、$、/、\和\0 (空字符)。
* 应全部小写。
* 最多64字节。

保留数据库名
* admin： 从权限的角度来看，这是"root"数据库。要是将一个用户添加到这个数据库，这个用户自动继承所有数据库的权限。一些特定的服务器端命令也只能从这个数据库运行，比如列出所有的数据库或者关闭服务器。
* local: 这个数据永远不会被复制，可以用来存储限于本地单台服务器的任意集合
* config: 当Mongo用于分片设置时，config数据库在内部使用，用于保存分片的相关信息。

# 6. 集合
集合就是 MongoDB 文档组，集合存在于数据库中，集合没有固定的结构，这意味着你在对集合可以插入不同格式和类型的数据，但通常情况下我们插入集合的数据都会有一定的关联性。

当第一个文档插入时，集合就会被创建。

**合法的集合名：**
* 集合名不能是空字符串""。
* 集合名不能含有\0字符（空字符)，这个字符表示集合名的结尾。
* 集合名不能以"system."开头，这是为系统集合保留的前缀。
* 用户创建的集合名字不能含有保留字符。有些驱动程序的确支持在集合名里面包含，这是因为某些系统生成的集合中包含该字符。除非你要访问这种系统创建的集合，否则千万不要在名字里出现$。　

## 6.1 capped collections
就是固定大小的collection。它有很高的性能以及队列过期的特性(过期按照插入的顺序)
它非常适合类似记录日志的功能和标准的 collection 不同，你必须要显式的创建一个capped collection，指定一个 collection 的大小，单位是字节。collection 的数据存储空间值提前分配的。

Capped collections 可以按照文档的插入顺序保存到集合中，而且这些文档在磁盘上存放位置也是按照插入顺序来保存的，所以当我们更新Capped collections 中文档的时候，更新后的文档不可以超过之前文档的大小，这样话就可以确保所有文档在磁盘上的位置一直保持不变。

由于 Capped collection 是按照文档的插入顺序而不是使用索引确定插入位置，这样的话可以提高增添数据的效率。MongoDB 的操作日志文件 oplog.rs 就是利用 Capped Collection 来实现的。

```
db.createCollection("mycoll", {capped:true, size:100000})
```

## 6.2 元数据
```
dbname.system.*
```

```
dbname.system.namespaces	列出所有名字空间。
dbname.system.indexes	列出所有索引。
dbname.system.profile	包含数据库概要(profile)信息。
dbname.system.users	列出所有可访问数据库的用户。
dbname.local.sources	包含复制对端（slave）的服务器信息和状态。
```

在{{system.indexes}}插入数据，可以创建索引。但除此之外该表信息是不可变的(特殊的drop index命令将自动更新相关信息)。
{{system.users}}是可修改的。 {{system.profile}}是可删除的。


# 7. 文档
注意：
* 文档中的键/值对是有序的。
* 文档中的值不仅可以是在双引号里面的字符串，还可以是其他几种数据类型（甚至可以是整个嵌入的文档)。
* MongoDB区分类型和大小写。
* MongoDB的文档不能有重复的键。
* 文档的键是字符串。除了少数例外情况，键可以使用任意UTF-8字符。

文档键命名规范：
* 键不能含有\0 (空字符)。这个字符用来表示键的结尾。
* .和$有特别的意义，只有在特定环境下才能使用。
* 以下划线"_"开头的键是保留的(不是严格要求的)。

## 7.2 数据类型

https://www.runoob.com/mongodb/mongodb-databases-documents-collections.html

* String	字符串。存储数据常用的数据类型。在 MongoDB 中，UTF-8 编码的字符串才是合法的。
* Integer	整型数值。用于存储数值。根据你所采用的服务器，可分为 32 位或 64 位。
* Boolean	布尔值。用于存储布尔值（真/假）。
* Double	双精度浮点值。用于存储浮点值。
* Min/Max keys	将一个值与 BSON（二进制的 JSON）元素的最低值和最高值相对比。
* Array	用于将数组或列表或多个值存储为一个键。
* Timestamp	时间戳。记录文档修改或添加的具体时间。
* Object	用于内嵌文档。
* Null	用于创建空值。
* Symbol	符号。该数据类型基本上等同于字符串类型，但不同的是，它一般用于采用特殊符号类型的语言。
* Date	日期时间。用 UNIX 时间格式来存储当前日期或时间。你可以指定自己的日期时间：创建 Date 对象，传入年月日信息。
* Object ID	对象 ID。用于创建文档的 ID。
* Binary Data	二进制数据。用于存储二进制数据。
* Code	代码类型。用于在文档中存储 JavaScript 代码。
* Regular expression	正则表达式类型。用于存储正则表达式。

## 7.3 查询
```
db.collection.find(query, projection)
db.collection.findOne(query, projection)
```
* query：用于查找文档的查询条件。默认为 {}，即匹配所有文档。
* projection（可选）：指定返回结果中包含或排除的字段。

### 7.3.1 条件操作符
https://www.runoob.com/mongodb/mongodb-operators.html

### 7.3.2 $type 操作符
https://www.runoob.com/mongodb/mongodb-operators-type.html

### 7.3.3 limit / skip
https://www.runoob.com/mongodb/mongodb-limit-skip.html

