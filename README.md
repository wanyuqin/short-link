# 短链系统

## 功能要求

* 类似于TinyURL或Bitly 的URL 缩短服务
* 客户端（用户）在系统中输入一个长 URL，系统返回一个缩短的 URL
* 访问短 URL 的客户端必须重定向到原始长 URL
* 多个用户输入相同的长网址，必须收到相同的短网址（一对一映射）
* 短 URL 应该易于阅读
* 短网址应避免冲突
* 短网址应该是不可预测的
* 客户端应该能够选择自定义短网址
* 短网址应该对网络爬虫友好（SEO）
* 短网址应支持分析（非实时），例如从缩短的网址进行的重定向次数
* 客户端可选择定义短 URL 的过期时间

## 非功能要求

* 高可用性
* 低延迟
* 高可扩展性
* 耐用的
* 容错

## 分表

* 基于short_url进行分表
* 创建user_id -> short_url 索引表
* 创建original索引表


### 短链表
```sql  
CREATE TABLE sl_link_0
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 0';

CREATE TABLE sl_link_1
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 1';

CREATE TABLE sl_link_2
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 2';

CREATE TABLE sl_link_3
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 3';

CREATE TABLE sl_link_4
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 4';

CREATE TABLE sl_link_5
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 5';

CREATE TABLE sl_link_6
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 6';

CREATE TABLE sl_link_7
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 7';

CREATE TABLE sl_link_8
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 8';

CREATE TABLE sl_link_9
(
    id         BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    origin_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    short_url  VARCHAR(50)      NOT NULL COMMENT '短链',
    expired_at BIGINT DEFAULT 0 NOT NULL COMMENT '过期时间',
    user_id    BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del     INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE     KEY idx_origin_url (origin_url),
    UNIQUE     KEY idx_short_url (short_url),
    INDEX      idx_user_id (user_id),
    INDEX      idx_created_at (created_at),
    INDEX      idx_updated_at (updated_at)
) COMMENT '短链接表 9';

```
### 短链用户映射表

```sql
CREATE TABLE sl_user_short_url (
    id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    short_url VARCHAR(50) NOT NULL COMMENT '短链',
    user_id BIGINT UNSIGNED DEFAULT 0 NOT NULL COMMENT '用户ID',
    created_at BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del INT DEFAULT 0 NOT NULL COMMENT '删除标志',
    INDEX idx_user_id (user_id),
    UNIQUE idx_short_url (short_url),
    INDEX idx_created_at (created_at),
    INDEX idx_updated_at (updated_at)
) COMMENT '用户短链接表';

```


### 短链-原始链接表

```sql

CREATE TABLE sl_original_short_url
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
    short_url    VARCHAR(50)      NOT NULL COMMENT '短链',
    original_url VARCHAR(255)     NOT NULL COMMENT '原始链接',
    created_at   BIGINT DEFAULT 0 NOT NULL COMMENT '创建时间',
    updated_at   BIGINT DEFAULT 0 NOT NULL COMMENT '更新时间',
    is_del       INT    DEFAULT 0 NOT NULL COMMENT '删除标志',
    UNIQUE       idx_short_url (short_url),
    UNIQUE       idx_original_url (original_url),
    INDEX        idx_created_at (created_at),
    INDEX        idx_updated_at (updated_at)
) COMMENT '短链与原始链接表';

```