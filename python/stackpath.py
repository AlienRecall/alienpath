import requests
import js2py
import uuid
import random
import re
import string


def randString(l):
    return "".join(
        random.choice(
            "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
        for i in range(l)
    )


def solver(r, sess, m, u, h, p):
    def sS1():
        r = s.get(
            redirectUri,
            headers={
                "Host": d,
                "sec-ch-ua": '"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"',
                "sec-ch-ua-mobile": "?0",
                "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
                "accept": "image/avif,image/webp,image/apng,image/*,*/*;q=0.8",
                "sec-fetch-site": "same-origin",
                "sec-fetch-mode": "no-cors",
                "sec-fetch-dest": "image",
                "referer": f"https://{d}/",
                "accept-language": "en-GB,en;q=0.9",
            },
        )
        if r.status_code < 400:
            return r.text
        else:
            print("Error submitting first request - Invalid response")

    def sS2():
        r = s.post(
            redirectUri,
            headers={
                "Host": d,
                "cache-control": "max-age=0",
                "sec-ch-ua": '"Google Chrome";v="87", " Not;A Brand";v="99", "Chromium";v="87"',
                "sec-ch-ua-mobile": "?0",
                "upgrade-insecure-requests": "1",
                "origin": f"https://{d}/",
                "content-type": "application/x-www-form-urlencoded",
                "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
                "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
                "sec-fetch-site": "same-origin",
                "sec-fetch-mode": "navigate",
                "sec-fetch-dest": "iframe",
                "referer": redirectUri,
                "accept-language": "en-GB,en-US;q=0.9,en;q=0.8,it;q=0.7",
            },
            data=postData,
        )
        if r.status_code < 400:
            return True
        else:
            print("Error submitting second request - Invalid response")

    def solve():
        print(s.cookies)
        if m.upper() == "GET":
            r = s.get(u, headers=h)
        else:
            r = s.post(u, headers=h, data=p)
        print(r.text, r.status_code)
        ss = ""
        for c in s.cookies:
            ss = ss + c.name + "=" + c.value + "; "
        print(ss)
        if r.status_code < 400:
            return r.text
        else:
            print("Error submitting final request - Invalid response")

    s = requests.Session()
    s.cookies = sess.cookies

    d = u.split("/")[2]
    tPVal = js2py.eval_js(
        re.findall(r"sbbgs\+'(.*?)'\+(.*?)\+",
                   r)[0][1].replace("dfx", "new Date()")
    )
    gprid = js2py.eval_js(re.findall(r"genPid\(\) {return (.*?) ;", r)[0])
    sbbgs = re.findall(r'sbbsv\("(.*?)"', r)[0].split("D-")[1]
    path = re.findall(r"window\.location\.port: ''\)\+'(.*?)'", r)[0]
    sPName = re.findall(r"prid \+ '(.*?)'", r)[0]
    tPName = re.findall(r"sbbgs\+'(.*?)'\+(.*?)\+", r)[0][0]
    params = f"{path}{gprid}{sPName}{sbbgs}{tPName}{tPVal}"
    redirectUri = f"https://{d}{params}"
    s.cookies.clear(domain="www.basket4ballers.com", path="/", name="UTGv2")
    sess.cookies["sbtsck"] = re.findall(r"sbtsck=(.*?);", r)[0]
    sess.cookies["PRLST"] = gprid
    s.cookies["UTGv2"] = sbbgs
    s1B = sS1()
    funcContent = re.findall(
        r"function xrv\(keyStr,sourceStr\)(.*?)function", s1B)[0]
    funcXRV = f"function xrv(keyStr, sourceStr){funcContent}"
    trstrup = re.findall(r'sbbdep\("(.*?)"', s1B)[0]
    arr = [f"{uuid.uuid4()}.local"]
    gvl = js2py.eval_js(f'{funcXRV} xrv("{trstrup.upper()}", "{str(arr)}")')
    lext = js2py.eval_js(f'{funcXRV} xrv("{trstrup.upper()}", [0,0])')
    bhvmsg = js2py.eval_js(
        f'{funcXRV} xrv("{trstrup.upper()}", "{randString(10)}-{randString(5)}")'
    )
    cdmsg = js2py.eval_js(
        f'{funcXRV} xrv("{trstrup.upper()}", "{randString(11)}-41-${randString(9)}-${randString(11)}-${randString(11)}-noieo-90.${random.randrange(2000000000000000, 9999999999999999)}")'
    )
    postData = {
        "cdmsg": cdmsg,
        "femsg": 1,
        "bhvmsg": bhvmsg,
        "futgs": "",
        "jsdk": trstrup,
        "glv": gvl,
        "lext": lext,
        "sdrv": 0,
    }
    s2R = sS2()
    solved = solve()

    return solved
