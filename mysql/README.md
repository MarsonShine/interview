# 索引相关知识点

## 索引命中问题

1. 索引字段应用了**函数**则无法命中索引
2. 非前缀模糊搜索无法命中索引
3. 数据表转码格式不同则无法命中索引
4. 联合索引没有符合**最左前缀匹配原则**则无法命中索引（查询时的条件列不是联合索引中的第一个列，索引失效）

## 聚族索引与非聚族索引

1. 聚族索引的顺序就是物理存储顺序
2. 聚族索引唯一
3. 非聚族索引的顺序与数据的物理存储顺序不同
4. 可以有多个非聚族所以

# 分析 MySQL 死锁

1. 查询死锁日志 `show engine innodb status`
2. 找出具体那条语句发生死锁
3. 分析sql加锁情况
4. 模拟死锁案发
5. 分析死锁日志
6. 分析死锁结果

# 优化 MySQL 性能

1. 加索引
2. 避免返回不必要的数据
3. 适当分批量进行
4. 优化 sql 结构
5. 分库分表
6. 读写分离

# 分库分表相关

## 分库分表原则

- 水平分库/表，以字段为依据，按照一定的规则（hash，range 等），将一个库/表的数据拆分到多个库/表中
- 垂直分库，按照业务属性的不同，将不同的表拆分到不同的库中
- 垂直分表，以字段为依据，按照字段的业务性，将表中字段拆分到不同的表中（如大数据量的电商系统，用户信息就可以按照用户 id 拆分）

# MySQL 底层数据结构相关

## 索引为什么选择 B+ 树而不是红黑树

**首先树作为索引结构，我们每查一次数据就需要从磁盘中读取一个节点。**

首先红黑树是一种平衡二叉树，虽然查询速度很快速且稳定，但是它一个节点下只有两个子节点，那么查询的 “性价比” 就不是很高。而 B 树可以存储更多节点数据，并且树的高度也会降低，自然与磁盘 I/O 的次数也就变少了。

那么为什么 MySQL **还是选择 B+ 树而不是 B 树**呢？

- B+ 树非叶子节点上是不存储数据的，仅存储键值。而 B 树不仅存储索引键值，还存储数据。而 MySQL 每个数据页都是尺寸限制的（默认 16 kb），如果不存储数据，那么就能存储更多的键值，那么自然树的高度就会更矮，从而能减少 I/O 的次数进而提高查询效率。
- 在查询方面要比 B 树更加简单方便。因为 B 树非叶子节点不仅存储键还有数据，那么在查询数据自然还要判断非叶子节点的数据。而 B+ 树就只需要根据父节点判断其数据的范围即可。

# 分布式 Id 的选择

- 数据库自增长 Id
- UUID
- Redis 生成 ID
- 分布式算法生成 ID（雪花算法）
- zookeeper 生成 ID
- MongoDB 的 objectid
- 可以通过分布式算法带有业务属性的 ID（如淘宝订单号？）

# 并发下安全更新数据

- 悲观锁，在更新这行数据的时候直接上锁，其它请求要等待（如 select for update）
- 乐观锁，有并发请求来了允许修改。如果没有其它请求修改这条信息，则修改成功。否则修改失败；这是通过添加标识符来实现的，如通过版本号机制或者是 CAS 算法实现
- 条件控制法，其实这种方法有点类似 MVCC 思想，用越严格的条件更新（参与条件判断的值已经快照），并发安全的可能性就会更高，因为当其它线程修改条件中的值时，执行更新时，条件就会不成立。

# exists 和 in 的性能对比

```mysql
select * from A where bid in (select id from B);
select * from A where exists (select 1 from B where A.bid = B.id);
```

这两个语句有啥区别？

核心原则：**小表驱动大表原则**

exists：先执行主查询，然后根据查询的结果再放到子查询中匹配查询条件，再返回结果。

根据核心原则，我们知道如果 A 是大表，那么 in 是合适的写法。而 A 是小表，那么 exists 就是合适的选择。

# 自增长 ID 可能会遇到什么问题

使用自增主键对数据库做分库分表，可能出现诸如主键重复等的问题。解决方案的话，简单点的话可以考虑使用 UUID，**自增主键会产生表锁**，从而引发问题。自增主键可能用完问题。

还有一个问题要注意，就是自增长 ID 会引起回溯（在 MySQL 重启的时候）

> 这个问题在 MySQL8.0 被修复了。

# MySQL 主从时间差（延时）问题

产生延时的原因：客户端连接 MySQL 连接数暴增，而在数据同步时只有一个线程在读取 binlog，当某个长时间 SQL 执行时，特别时在阻塞的情况下，大量的 SQL 查询积压未被线程读取，这样就存在与从库的数据不一致。也就是主从延时。

主备流程：

1. 主库的更新事件(update、insert、delete)被写到 binlog
2. 从库发起连接，连接到主库。
3. 此时主库创建一个 binlog dump thread，把 binlog 的内容发送到从库。
4. 从库启动之后，创建一个 I/O 线程，读取主库传过来的 binlog 内容并写入到 relay log
5. 从库还会创建一个 SQL 线程，从 relay log 里面读取内容，从 Exec_Master_Log_Pos 位置开始执行读取到的更新事件，将更新内容写入到 slave 的 db

解决方案：

- 增加从服务器，这个目的还是分散读的压力，从而降低服务器负载
- 选择更好的硬件设备作为 slave（分析 ready log）提升读写能力

# 建立数据库连接的过程（为什么连接数据库开销很大）

1. 通过 TCP 协议的三次握手和数据库服务器建立连接
2. 发送数据库账号和密码，等待服务器进行身份验证
3. 完成验证之后提交 SQL 语句执行
4. 连接关闭，TCP 四次挥手关闭连接

频繁进行数据库交互就会重复执行上面的过程。

而连接池就是为了避免这个事情：在内部对象池中维护一定数量的数据库连接，暴露给客户端连接。

好处：连接复用；响应更快；统一管理，避免数据库连接泄露；

# 如何优化一条长耗时的 SQL

- 查看是否涉及多表和子查询，优化 Sql 结构，如去除冗余字段，是否可拆表等
- 优化索引结构，看是否可以适当添加索引
- 数量大的表，可以考虑进行分离/分表（如交易流水表）
- 数据库主从分离，读写分离
- explain 分析 sql 语句，查看执行计划，优化 sql
- 查看 mysql 执行日志，分析是否有其他方面的问题

# 百万级别或以上的数据如何删除

- 删除索引
- 然后批量删除要删除的数据
- 重建索引

# 如何分析 MySQL 中的慢查询

首先得先知道哪些是慢查询的 SQL：

1. MySQL 参数设置可以获取

   1. 设置 `slow_query_log=on` 就可以捕获超过一定数值的 SQL 语句
   2. 当 SQL 查询时间超过 `long_query_time` 设定的值，就会被记录到日志中
   3. `slow_query_log` 日志文件名

2. 启用 MySQL 慢查询

   ```cmd
   log-slow-queries=/data/mysqldata/slowquery.log
   long_query_time=2
   ```

   通过执行 `show variables like 'log_slow_queries'` 可以查看是否启动慢查询

   `show variables like 'long_query_time';` 查看有多少条慢查询

3. show processlist

4. explain sql，在已知慢查询的 SQL 的情况来分析

# MySQL 服务器 CPU 飙升的解决方案

排查：

- 使用 top 命令观察，确定是 mysqld 导致还是其他原因
- 如果是 mysqld 导致的，`show processlist`，查看 session 情况，确定是不是有消耗资源的 sql 在运行。
- 找出消耗高的 sql，看看执行计划是否准确， 索引是否缺失，数据量是否太大。

解决：

- kill 掉这些线程(同时观察 cpu 使用率是否下降)
- 进行相应的调整(比如说加索引、改 sql、改内存参数)

其它：

导致 CPU 飙升的情况还有短时间内连接数爆升，这种情况就要配合对应的应用客户端分析，看是否需要限制连接数等。

# 主备复制

- 主数据库有个 bin-log 二进制文件，纪录了所有增删改 Sql 语句。（binlog 线程）
- 从数据库把主数据库的 bin-log 文件的 sql 语句复制过来。（I/O 线程）
- 从数据库的 ready-log 重做日志文件中再执行一次这些 sql 语句。（Sql 执行线程，分析 ready-log 文件）


