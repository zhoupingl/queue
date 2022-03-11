# 二叉最小堆

> 该项目当前只支持单个队列（支持多队列，略微调整代码）
> brandy_head.go 是二叉最小堆
> queue.go 是队列

## 接口


## 添加任务
##### EndPointbe
/queue/add
#### Type
GET
### Query Params
id int  数据ID

class int 等级

---


## 从队列拿取一个任务
> 从队列拿取一个任务（30s 服务器循环,在该时间范围获取到了数据。立即返回。否则超时返回）
##### EndPointbe
/queue/pull
#### Type
GET

---


## 标记任务完成
##### EndPointbe
/queue/success
#### Type
GET
### Query Params
id int

---

## 标记任务失败
> 标记任务失败。任务从新加入队列
##### EndPointbe
/queue/rejoin
#### Type
GET
### Query Params
id int

---


## 队列服务version
##### EndPointbe
/queue/version
#### Type
GET

