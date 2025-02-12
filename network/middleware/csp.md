# 简介
内容安全策略（CSP）是一个额外的安全层，用于检测并削弱某些特定类型的攻击，包括跨站脚本（XSS）和数据注入攻击等。

为使 CSP 可用，你需要配置你的网络服务器返回 Content-Security-Policy HTTP 标头

https://developer.mozilla.org/zh-CN/docs/Web/HTTP/CSP

## 缓解跨站脚本攻击
CSP 的主要目标是减少和报告 XSS 攻击。XSS 攻击利用了浏览器对于从服务器所获取的内容的信任。恶意脚本在受害者的浏览器中得以运行，因为浏览器信任其内容来源，即使有的时候这些脚本并非来自于它本该来的地方。

CSP 通过指定有效域——即浏览器认可的可执行脚本的有效来源——使服务器管理者有能力减少或消除 XSS 攻击所依赖的载体。一个 CSP 兼容的浏览器将会仅执行从白名单域获取到的脚本文件，忽略所有的其他脚本（包括内联脚本和 HTML 的事件处理属性）。

## 缓解数据包嗅探攻击
除限制可以加载内容的域，服务器还可指明哪种协议允许使用；比如（从理想化的安全角度来说），服务器可指定所有内容必须通过 HTTPS 加载。一个完整的数据安全传输策略不仅强制使用 HTTPS 进行数据传输，也为所有的 cookie 标记 secure 标识，并且提供自动的重定向使得 HTTP 页面导向 HTTPS 版本。网站也可以使用 Strict-Transport-Security HTTP 标头确保连接它的浏览器只使用加密通道。


# 示例

1. 一个网站管理者想要所有内容均来自站点的同一个源（不包括其子域名）。
```azure
Content-Security-Policy: default-src 'self'
```

2. 一个网站管理者允许内容来自信任的域名及其子域名（域名不必须与 CSP 设置所在的域名相同）。
```azure
Content-Security-Policy: default-src 'self' *.trusted.com
```

3. 一个网站管理者允许网页应用的用户在他们自己的内容中包含来自任何源的图片，但是限制音频或视频需从信任的资源提供者，所有脚本必须从特定主机服务器获取可信的代码。
```azure
Content-Security-Policy: default-src 'self'; img-src *; media-src media1.com media2.com; script-src userscripts.example.com
```
* 图片可以从任何地方加载 (注意“*”通配符)。
* 多媒体文件仅允许从 media1.com 和 media2.com 加载（不允许从这些站点的子域名）。
* 可运行脚本仅允许来自于 userscripts.example.com。

4. 一个线上银行网站的管理者想要确保网站的所有内容都要通过 SSL 方式获取，以避免攻击者窃听用户发出的请求。
```azure
Content-Security-Policy: default-src https://onlinebanking.jumbobank.com
```
该服务器仅允许通过 HTTPS 方式并仅从 onlinebanking.jumbobank.com 域名来访问文档。

5. 一个在线邮箱的管理者想要允许在邮件里包含 HTML，同样图片允许从任何地方加载，但不允许 JavaScript 或者其他潜在的危险内容（从任意位置加载）。
```azure
Content-Security-Policy: default-src 'self' *.mailsite.com; img-src *
```
注意这个示例并未指定 script-src；在此 CSP 示例中，站点通过 default-src 指令的对其进行配置，这也同样意味着脚本文件仅允许从原始服务器获取。


