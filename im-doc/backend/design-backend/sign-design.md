# Sign

## 签名验证算法

遍历json或者是query参数中的所有除`sign`, `EIO`, `transport`等框架自带参数, 按照key排序后`keyvalue`组合, 最后加上唯一的appkey进行sha1加密返回大写字符串  
example:

```javascript
{
  "name": "Jack",
  "age": 20,
  "sign": "5347011AB54914FE916A84924CDD517F702888D4",
}
```

假设`appkey`为`88888888`,混合后的字符串应为`age20nameJack88888888`,经过sha1加密后字符串:

```bash
5347011AB54914FE916A84924CDD517F702888D4
```

