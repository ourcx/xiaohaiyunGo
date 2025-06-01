package v1

//这是sequenceDiagram
//    客户端->>服务端: 注销请求（携带JWT）
//    服务端->>Redis: SET jwt_blacklist:<token_id> "invalidated" EX <剩余时间>
//    Redis-->>服务端: 操作结果
//    服务端->>客户端: 注销成功
//检验是否有jwt到jwt是否有效再到检查redis黑名单，不存在则进行清除，并且做一个定期清除redis黑名单的策略
