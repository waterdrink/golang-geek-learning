# redis性能测试

redis版本：6.2.6


redis brenchmark本地测试结果如下：

set
| key size | result                                                                  |
| -------- | ----------------------------------------------------------------------- |
| 10       | 100000 requests completed in 2.58 seconds, 38789.76 requests per second |
| 20       | 100000 requests completed in 2.46 seconds, 40567.95 requests per second |
| 50       | 100000 requests completed in 2.66 seconds, 37551.63 requests per second |
| 100      | 100000 requests completed in 2.46 seconds, 40683.48 requests per second |
| 200      | 100000 requests completed in 2.65 seconds, 37707.39 requests per second |
| 1k       | 100000 requests completed in 2.58 seconds, 38759.69 requests per second |
| 5k       | 100000 requests completed in 2.56 seconds, 39138.94 requests per second |


get
| key size | result                                                                  |
| -------- | ----------------------------------------------------------------------- |
| 10       | 100000 requests completed in 2.51 seconds, 39777.25 requests per second |
| 20       | 100000 requests completed in 3.03 seconds, 33057.85 requests per second |
| 50       | 100000 requests completed in 2.43 seconds, 41135.34 requests per second |
| 100      | 100000 requests completed in 2.50 seconds, 39952.06 requests per second |
| 200      | 100000 requests completed in 2.53 seconds, 39588.28 requests per second |
| 1k       | 100000 requests completed in 2.58 seconds, 38729.67 requests per second |
| 5k       | 100000 requests completed in 2.70 seconds, 37009.62 requests per second |

可见 在排除网络因素的测试条件下，key size 对性能影响不大

# redis单key空间占用测试

每次测试前重启redis server，避免内存碎片对结果的影响
测试方法：debug populate 1000 test-value-10 10

| value size | count | 平均每个key占用 |
| ---------- | ----- | --------------- |
| 10         | 1k    | 88 字节         |
| 10         | 2k    | 88 字节         |
| 10         | 5k    | 93 字节         |
| 10         | 10k   | 93 字节         |
| 10         | 20k   | 93 字节         |
| 10         | 50k   | 90 字节         |


| value size | count | 平均每个key占用 |
| ---------- | ----- | --------------- |
| 20         | 1k    | 96 字节         |
| 20         | 2k    | 96 字节         |
| 20         | 5k    | 101 字节        |
| 20         | 10k   | 102 字节        |
| 20         | 20k   | 101 字节        |
| 20         | 50k   | 98 字节         |

| value size | count | 平均每个key占用 |
| ---------- | ----- | --------------- |
| 50         | 1k    | 128 字节        |
| 50         | 2k    | 128 字节        |
| 50         | 5k    | 133 字节        |
| 50         | 10k   | 133 字节        |
| 50         | 20k   | 133 字节        |
| 50         | 50k   | 130 字节        |

| value size | count | 平均每个key占用 |
| ---------- | ----- | --------------- |
| 5k         | 1k    | 6,216 字节      |
| 5k         | 2k    | 6,216 字节      |
| 5k         | 5k    | 6,221 字节      |
| 5k         | 10k   | 6,221 字节      |
| 5k         | 20k   | 6,221 字节      |
| 5k         | 50k   | 6,218 字节      |

可见 
1. redis实际使用的内存会略大于键值对的实际大小(考虑redis存储内存模型是由结构体dictEntry和redisObject封装，且字符串对象由SDS表示，SDS类型会有free的余量，为了减少变更时发生的内存重分配)
2. 键值对的大小越大，实际使用的内存越贴近于键值对本身的大小(考虑构体dictEntry和redisObject占用的空间比较小，且SDS在实际存储的值大于1M时，free的余量最多只会占用1M的空间)
3. 在键值对大小相同的情况下，数量>20k和<2k时，实际使用的内存越小(越有效率)