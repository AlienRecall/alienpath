const requestModule = require('request-promise-native');
const solver = require('./index');

(async () => {
  let jar = requestModule.jar()
  const request = requestModule.defaults({
    jar: jar,
    gzip: true,
    forever: true,
    simple: false,
    timeout: 20000,
    followAllRedirects: true,
    resolveWithFullResponse: true,
  })
  
  request({
    method: 'GET',
    uri: 'https://www.basket4ballers.com/fr/authentification?back=my-account',
    headers: {
      'cache-control': 'max-age=0',
      'sec-ch-ua': '"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"',
      'sec-ch-ua-mobile': '?0',
      'upgrade-insecure-requests': '1',
      'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36',
      'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
      'sec-fetch-site': 'same-origin',
      'sec-fetch-mode': 'navigate',
      'sec-fetch-dest': 'document',
      'referer': 'https://www.basket4ballers.com/fr/authentification?back=my-account',
      'accept-language': 'en-GB,en;q=0.9',
    }
  }).then(async (res) => {
    console.log(res.body.match(/<title>(.*?)<\//)[1], res.statusCode)
    if (res.body.match(/<title>(.*?)<\//).length === 2 && res.body.match(/<title>(.*?)<\//)[1] === "StackPath" && !res.body.includes('Please confirm that you are human, by typing the characters shown here')) {
      let solved = await solver(res.body, jar, 'www.basket4ballers.com', {
        method: 'GET',
        uri: 'https://www.basket4ballers.com/fr/authentification?back=my-account',
        headers: {
          'cache-control': 'max-age=0',
          'sec-ch-ua': '"Chromium";v="88", "Google Chrome";v="88", ";Not A Brand";v="99"',
          'sec-ch-ua-mobile': '?0',
          'upgrade-insecure-requests': '1',
          'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36',
          'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
          'sec-fetch-site': 'same-origin',
          'sec-fetch-mode': 'navigate',
          'sec-fetch-dest': 'document',
          'referer': 'https://www.basket4ballers.com/fr/authentification?back=my-account',
          'accept-language': 'en-GB,en;q=0.9',
        }
      })
      console.log(solved[0].body)
      jar = solved[1]
      console.log('solved')
    }
  })
})();