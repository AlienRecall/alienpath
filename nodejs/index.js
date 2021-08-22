const requestModule = require('request-promise-native');
const crypto = require('crypto');
const uuid = require('uuid')

function randString(l) {
  let result           = '';
  let characters       = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  for (let i = 0; i < l; i++) {
     result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
}

function randInt(min, max) {
  return Math.floor(Math.random() * (Math.floor(max) - Math.ceil(min) + 1)) + Math.ceil(min);
}

function solver(schallengeSRC, scookieJar, sdomain, sopts) {

  let challengeSRC = schallengeSRC
  let domain = sdomain
  let opts = sopts
  let challengeURL;
  let newChallengeSRC;
  let challengeFORM = {
    cdmsg: "",
    femsg: 1,
    bhvmsg: "",
    futgs: "",
    jsdk: "",
    glv: "",
    lext: "",
    sdrv: 0
  }

  let cookieJar = scookieJar
  let request = requestModule.defaults({
    gzip: true,
    forever: true,
    simple: false,
    timeout: 20000,
    jar: cookieJar,
    followAllRedirects: true,
    resolveWithFullResponse: true,
    ciphers: crypto.constants.defaultCipherList + ':!ECDHE+SHA:!AES128-SHA',
    secureProtocol: "TLSv1_2_method",
    strictSSL: false
  });

  let solved = solve()

  async function solve() {

    let sbtsck = getSBTSCK()
    setCookie("sbtsck", sbtsck)
    let gprid = getGPRID()
    setCookie("PRLST", gprid)
    let sbbgs = challengeSRC.match(/sbbsv\("D-(.*?)"/)[1]
    deleteCookie("UTGv2")
    setCookie("UTGv2", sbbgs)
    let ddl = getDDL()
    challengeURL = `https://${domain}/sbbi/?sbbpg=sbbShell&gprid=${gprid}&sbbgs=${sbbgs}&ddl=${ddl}`

    newChallengeSRC = await firstRedirect()
    let adOtr = getADOTR()
    setCookie("adOtr", adOtr)
    let trstr = getTRSTR()
    challengeFORM.jsdk = trstr
    let sbbjglv = [`${uuid.v4()}.local`]
    challengeFORM.glv = xrv(trstr.toUpperCase(), String(sbbjglv))
    challengeFORM.lext = xrv(trstr.toUpperCase(), "[0,0]")
    challengeFORM.bhvmsg = xrv(trstr.toUpperCase(), `${randString(10)}-${randString(5)}`)
    challengeFORM.cdmsg = xrv(trstr.toUpperCase(), `${randString(11)}-41-${randString(9)}-${randString(11)}-${randString(11)}-noieo-90.${randInt(2000000000000000, 9999999999999999)}`) //xrv(trstr.toUpperCase(), "c3wrw6zmi89-41-btzpsbqrr-88mb7is8ih4-v2osmr4iefe-noieo-90.3095389639745667")

    console.log(challengeFORM)
    let solvedSRC = await submitChallenge()
    console.log(solvedSRC)
    let solvedRES = await getSolvedResponse()
    return [solvedRES, cookieJar]
  }

  function deleteCookie(name) {
    delete cookieJar._jar.store.idx[domain]['/'][name]
  }

  function setCookie(name, value) {
    cookieJar.setCookie(requestModule.cookie(name + '=' + value), `https://${domain}`)
  }

  function getSBTSCK() {
    return challengeSRC.match(/"sbtsck=(.*?);/)[1]
  }

  function getGPRID() {
    return eval(challengeSRC.match(/genPid\(\) {return (.*?) ;/)[1])
  }

  function getDDL() {
    return eval(challengeSRC.match(/'&ddl='\+(.*?)\+/)[1].replace('dfx', 'new Date()'))
  }

  function getADOTR() {
    return eval(newChallengeSRC.match(/parent\.otr = (.*?);/)[1]).join('')
  }

  function getTRSTR() {
    return newChallengeSRC.match(/sbbdep\("(.*?)"/)[1]
  }

  function xrv(keyStr, sourceStr) {
    var keyLength = keyStr.length;
    var targetStr = "";
    var i, rPos, a, b, c, d, targetStr;
    for (i = 0; i < sourceStr.length; i++) {
      rPos = i % keyLength;
      a = sourceStr.charCodeAt(i);
      b = keyStr.charCodeAt(rPos);
      c = a ^ b;
      d = c + "a";
      targetStr = targetStr + d;
    }
    return targetStr;
  }

  async function firstRedirect() {
    const res = await request({
      method: "GET",
      uri: challengeURL,
      headers: {
        'authority': 'www.basket4ballers.com',
        'upgrade-insecure-requests': '1',
        'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36',
        'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
        'sec-fetch-site': 'same-origin',
        'sec-fetch-mode': 'navigate',
        'sec-fetch-dest': 'iframe',
        'referer': 'https://www.basket4ballers.com/fr/authentification?back=my-account',
        'accept-language': 'en-GB,en;q=0.9',
      }
    })
    if (res.statusCode < 400) {
      return res.body
    } else {
      console.log("Error getting second challenge")
    }
  }

  async function submitChallenge() {
    const res = await request({
      method: "POST",
      uri: challengeURL,
      headers: {
        'authority': 'www.basket4ballers.com',
        'cache-control': 'max-age=0',
        'upgrade-insecure-requests': '1',
        'origin': 'https://www.basket4ballers.com',
        'content-type': 'application/x-www-form-urlencoded',
        'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36',
        'accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9',
        'sec-fetch-site': 'same-origin',
        'sec-fetch-mode': 'navigate',
        'sec-fetch-dest': 'iframe',
        'referer': challengeURL,
        'accept-language': 'en-GB,en;q=0.9',
      },
      form: challengeFORM
    })
    if (res.statusCode < 400) {
      return res.body
    } else {
      console.log("Error submitting second challenge")
    }
  }

  async function getSolvedResponse() {
    const res = await request(opts)
    if (res.statusCode < 400) {
      return res.body
    } else {
      console.log("Error getting solved response")
    }
  }

  return solved
}

module.exports = solver
