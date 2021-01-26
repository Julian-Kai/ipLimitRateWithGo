### Homework

---
+ 作業講解
    1. 使用 gin framework 作為 rest 接口，進而使用 middleware 來達到限流
    2. 使用令牌桶算法來達到限流的效果
---
+ 啟動方法
    ```
  # 終端輸入以下指令, 啟動伺服器, 等待 request
  make run
    ```
---
+ 測試方法
    ```
  # 終端輸入以下指令, 產生三個不同的 ip 對 api 接口各自請求 80 次, 最後印出統計的結果
  make test_ping
    ```
