import requests
import re
import stackpath

s = requests.Session()

# example code
r = s.get(
    "https://www.basket4ballers.com/en/authentification?back=my-account",
    headers={
        "authority": "www.basket4ballers.com",
        "upgrade-insecure-requests": "1",
        "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
        "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
        "sec-fetch-site": "same-origin",
        "sec-fetch-mode": "navigate",
        "sec-fetch-dest": "document",
        "referer": "https://www.basket4ballers.com/",
        "accept-language": "en-GB,en;q=0.9",
    },
)
print(r.status_code)
if "StackPath" in re.findall(r"<title>(.*?)<\/", r.text)[0]:
    print("StackPath Detected!")
    res = stackpath.solver(
        r.text,
        s,
        "GET",
        "https://www.basket4ballers.com/en/authentification?back=my-account",
        {
            "authority": "www.basket4ballers.com",
            "upgrade-insecure-requests": "1",
            "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36",
            "accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
            "sec-fetch-site": "same-origin",
            "sec-fetch-mode": "navigate",
            "sec-fetch-dest": "document",
            "referer": "https://www.basket4ballers.com/",
            "accept-language": "en-GB,en;q=0.9",
        },
        False,
    )
    body = res[0]
    cookies = res[1]
    print(cookies, body)
