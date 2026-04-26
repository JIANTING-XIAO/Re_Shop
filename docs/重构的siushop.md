```
Backend/
  cmd/
    api/
      main.go                 # 程序入口，启动 HTTP 服务

  internal/
    app/
      bootstrap/              # 应用启动装配
      router/                 # 路由注册
      server/                 # HTTP 服务初始化

    shared/
      config/                 # 配置读取
      middleware/             # 中间件：JWT、日志、限流、恢复等
      response/               # 统一响应封装
      errors/                 # 业务错误码、错误定义
      utils/                  # 轻量工具函数
      db/                     # 数据库连接管理
      cache/                  # Redis 等缓存封装
      lock/                   # 分布式锁封装
      mq/                     # 消息队列封装
      idempotency/            # 幂等处理
      auth/                   # 鉴权能力
      logger/                 # 日志能力

    modules/
      user/
        controller/           # 用户接口层
        service/              # 用户业务逻辑
        repository/           # 用户数据访问
        model/                # 用户实体/数据库模型
        dto/                  # 用户请求响应对象
        routes.go             # 用户模块路由注册

      product/
        controller/           # 商品接口
        service/              # 商品业务
        repository/           # 商品数据访问
        model/                # 商品/SPU/SKU模型
        dto/                  # 商品相关 DTO
        routes.go

      cart/
        controller/           # 购物车接口
        service/
        repository/
        model/
        dto/
        routes.go

      order/
        controller/           # 订单接口
        service/              # 下单、取消、查询、关闭订单
        repository/
        model/
        dto/
        routes.go

      payment/
        controller/           # 支付下单、回调
        service/
        repository/
        model/
        dto/
        routes.go

      inventory/
        controller/           # 库存接口
        service/              # 扣减、回补、锁定
        repository/
        model/
        dto/
        routes.go

      seckill/
        controller/           # 秒杀活动、秒杀下单
        service/              # 秒杀核心逻辑
        repository/
        model/
        dto/
        routes.go

      coupon/
        controller/           # 优惠券接口
        service/
        repository/
        model/
        dto/
        routes.go

      trade/
        controller/           # 交易确认页、提交订单编排
        service/              # 聚合商品、库存、优惠、运费等
        repository/           # 可选，没有就不放
        dto/
        routes.go

  pkg/
    snowflake/                # ID 生成
    timeutil/                 # 时间工具
    validator/                # 参数校验扩展

  docs/                       # 接口文档、设计文档
  scripts/                    # 部署/初始化脚本
  deployments/                # 容器、部署配置

```



### 1. controller

负责：

- 接收 HTTP 请求
- 参数绑定与基本校验
- 调用 service
- 返回统一响应

不要放：

- 核心业务逻辑
- 数据库操作
- 复杂分支判断

------

### 2. service

负责：

- 核心业务逻辑
- 多个 repository 的协作
- 事务控制
- 业务状态流转
- 调用缓存、消息队列、锁等基础设施

这是每个模块的“大脑”。

------

### 3. repository

负责：

- 数据库 CRUD
- 查询封装
- 持久化实现

不要把业务规则写在这里。
repository 应该只关心“怎么查、怎么存”。

------

### 4. model

负责：

- 数据库实体
- 表结构映射对象
- 领域核心数据对象

建议分清：

- 数据库存储模型
- 对外返回 DTO

不要直接把数据库对象原样返回给前端。

------

### 5. dto

负责：

- 请求对象 request
- 响应对象 response

这样做的好处是字段变更更稳，不会数据库一改，接口全崩。

------

### 6. shared

这里放真正跨模块复用的能力，比如：

- JWT 鉴权
- Redis
- 限流
- 幂等
- 日志
- 统一错误码
- 统一响应体

注意：
shared/utils 一定要克制。很多项目最后最乱的地方就是 utils/。

------

## 四、模块划分建议

对于“普通电商 + 秒杀”，我建议至少有这几个核心模块。

------

### 1. user 用户模块

负责：

- 注册登录
- 用户信息
- 地址管理
- 用户状态
- 权限/角色（如果有后台）

------

### 2. product 商品模块

负责：

- 商品基础信息
- 分类、品牌
- SPU / SKU
- 商品详情
- 上下架
- 普通售价

它只关注“商品本身”，不要把秒杀活动逻辑塞进来。

------

### 3. inventory 库存模块

负责：

- 普通库存
- 锁定库存
- 扣减库存
- 回补库存

为什么建议单独拆出来？
因为库存几乎会被：

- 商品
- 秒杀
- 订单
- 支付回调
- 取消订单

这些模块同时依赖，单独拆边界更清楚。

------

### 4. cart 购物车模块

负责：

- 加入购物车
- 修改数量
- 勾选商品
- 结算前查询

秒杀商品通常不完全走普通购物车逻辑，但普通电商部分还是需要它。

------

### 5. trade 交易编排模块

这个模块很重要，也很容易被忽略。

负责：

- 确认订单页信息聚合
- 校验商品状态
- 校验库存
- 计算价格
- 运费、优惠券、总金额汇总
- 提交订单前最终校验

它更像“交易流程协调者”，而不是单一领域模块。

------

### 6. order 订单模块

负责：

- 创建订单
- 订单查询
- 取消订单
- 关闭超时订单
- 订单状态流转

建议让它只关心订单本身，不直接承载过多“优惠、秒杀、支付”细节。

------

### 7. payment 支付模块

负责：

- 创建支付单
- 支付状态更新
- 支付回调处理
- 支付与订单联动

它和订单关系很紧，但最好独立出来，不要混在订单模块里。

------

### 8. coupon 优惠券模块

负责：

- 领券
- 可用券查询
- 用券校验
- 用券核销

如果你前期不做优惠券，可以先预留结构。

------

### 9. seckill 秒杀模块

这是你项目的核心差异点。

负责：

- 秒杀活动创建与管理
- 秒杀商品绑定
- 秒杀价格
- 秒杀库存
- 限购规则
- 秒杀资格校验
- 秒杀下单流程
- 防重复抢购
- 高并发保护