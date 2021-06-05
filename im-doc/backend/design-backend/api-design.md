# API

{% hint style="danger" %}
 所有的Body Parameters均使用json格式传输，即`content-type=application/json`
{% endhint %}

{% api-method method="post" host="https://sso.im.peanut996.cn" path="/login" %}
{% api-method-summary %}
/login
{% endapi-method-summary %}

{% api-method-description %}
用户登录 
{% endapi-method-description %}

{% api-method-spec %}
{% api-method-request %}
{% api-method-body-parameters %}
{% api-method-parameter name="sign" type="string" required=true %}
签名验证 
{% endapi-method-parameter %}

{% api-method-parameter name="account" type="string" required=true %}
用户账号
{% endapi-method-parameter %}

{% api-method-parameter name="password" type="string" required=true %}
 用户密码
{% endapi-method-parameter %}
{% endapi-method-body-parameters %}
{% endapi-method-request %}

{% api-method-response %}
{% api-method-response-example httpCode=200 %}
{% api-method-response-example-description %}
login success
{% endapi-method-response-example-description %}

```javascript
{
    "code": 0,
    "message": "Success",
    "data": {
        "token": "1C42B01254471AFB12AB5D47C511E87C529FA3A3",
        "uid": "1387471245368365056"
    }
}
```
{% endapi-method-response-example %}
{% endapi-method-response %}
{% endapi-method-spec %}
{% endapi-method %}

{% api-method method="post" host="https://sso.im.peanut996.cn" path="/register" %}
{% api-method-summary %}
/register
{% endapi-method-summary %}

{% api-method-description %}
用户注册
{% endapi-method-description %}

{% api-method-spec %}
{% api-method-request %}
{% api-method-body-parameters %}
{% api-method-parameter name="sign" type="string" required=true %}
签名验证
{% endapi-method-parameter %}

{% api-method-parameter name="account" type="string" required=true %}
用户信息 
{% endapi-method-parameter %}

{% api-method-parameter name="password" type="string" required=true %}
用户密码
{% endapi-method-parameter %}
{% endapi-method-body-parameters %}
{% endapi-method-request %}

{% api-method-response %}
{% api-method-response-example httpCode=200 %}
{% api-method-response-example-description %}

{% endapi-method-response-example-description %}

```javascript
{
    "code": 0,
    "message": "Success",
    "data": {
        "token": "1C42B01254471AFB12AB5D47C511E87C529FA3A3",
        "uid": "1387471245368365056"
    }
}
```
{% endapi-method-response-example %}
{% endapi-method-response %}
{% endapi-method-spec %}
{% endapi-method %}

{% hint style="danger" %}
用户密码应该被sha1加密后存入数据库
{% endhint %}

{% api-method method="post" host="https://sso.im.peanut996.cn" path="/logout" %}
{% api-method-summary %}
/logout
{% endapi-method-summary %}

{% api-method-description %}
 用户注销
{% endapi-method-description %}

{% api-method-spec %}
{% api-method-request %}
{% api-method-body-parameters %}
{% api-method-parameter name="sign" type="string" required=true %}
 签名验证
{% endapi-method-parameter %}

{% api-method-parameter name="uid" type="string" required=true %}
 用户唯一ID
{% endapi-method-parameter %}

{% api-method-parameter name="token" type="string" required=true %}
 用户登录系统凭证
{% endapi-method-parameter %}
{% endapi-method-body-parameters %}
{% endapi-method-request %}

{% api-method-response %}
{% api-method-response-example httpCode=200 %}
{% api-method-response-example-description %}

{% endapi-method-response-example-description %}

```javascript
{
    "code": 0,
    "message": "Success",
    "data": null
}
```
{% endapi-method-response-example %}
{% endapi-method-response %}
{% endapi-method-spec %}
{% endapi-method %}

