# Vercel Proxy

Simple http proxy for Vercel.

### Usage

```javascript
fetch("https://project-name.vercel.app/https://example.com?param1=value1&param2=value2")
.then((res) => res.text())
.then(console.log.bind(console))
.catch(console.error.bind(console));

```

```bash
curl https://project-name.vercel.app/https:/example.com?param1=value1&param2=value2
```
> 注意 `curl` 这里需要把 `https://` 换成 `https:/`，否则会报错。