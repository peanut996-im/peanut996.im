# Token

通过验证用户的账号密码，获得连接整个系统的token，系统的其他应用均需要带有token以及sign双重验证，token验证失败的用户将自动跳转至SSO登录界面。

## 鉴权

通过用户登录带有的参数与后台数据库密码比较，成功则返回对应的token, 验证失败则返回对应的错误信息。

## token设计

采用sha1+base64双重加密，采用最简单的`uid_timestamp`的加密方法，由于uid采用雪花算法一定唯一，所以产生的token也最大限度避免了重复，数据量小的情况可默认为不重复：

```bash
# 全大写字母
upper(sha1(uid+timestamp))
```



token自动过期时间为`24`小时

## 存储

鉴权通过后，将会存储token至redis服务器，存在需求:

```text
uid => token
token => uid
```

有需求可知，共需要两组`string`的键值对。uid=&gt;token的键命名为`${uid}_to_token` , 相应的token=&gt;uid的键命名为`${token}_to_uid`, 示例如下：

```bash
# uid 1387109626431934464
# token 5347011AB54914FE916A84924CDD517F702888D4

1387109626431934464_to_token => 5347011AB54914FE916A84924CDD517F702888D4
5347011AB54914FE916A84924CDD517F702888D4_to_uid => 1387109626431934464
```

