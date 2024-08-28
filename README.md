# Vercel Proxy

Simple http proxy for Vercel.

### Delpoy

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2FTBXark%2Fvercel-proxy)

### Usage

```javascript
fetch("https://project-name.vercel.app/https://example.com?param1=value1&param2=value2")
.then((res) => res.text())
.then(console.log.bind(console))
.catch(console.error.bind(console));

```

```bash
curl -L https://project-name.vercel.app/https:/example.com?param1=value1&param2=value2
```

Just add `https://project-name.vercel.app/` before the url you want to proxy.

### License

**vercel-proxy** is released under the MIT license. [See LICENSE](LICENSE) for details.