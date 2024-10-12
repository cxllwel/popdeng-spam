import httpx,time

Cli=httpx.Client(headers={
    "cache-control":"max-age=0",
    "user-agent":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36"
},timeout=60)
buildID=None
gatewayL=None
while True:
    try:
        resp=Cli.get("https://popdeng.click/")
        BIDnow=resp.text.split("static/chunks/app/page-")[1].split(".js")[0]
        gateway=resp.text.split("host\\\":\\\"")[1].split("\\\",")[0]
        if buildID != BIDnow or gateway != gatewayL:
            if gateway != gatewayL:
                print("Found New Update of gateway mqtt new",gateway,"old",gatewayL)
            else:
                print("Found New Update of page.js ID :",BIDnow)
            gatewayL=gateway
            buildID=BIDnow
            resp=Cli.get(f"https://popdeng.click/_next/static/chunks/app/page-{buildID}.js")
            # gateway=resp.text.split('connect("wss://')[1].split('",')[0]
            clickid=resp.text.split('publish("')[1].split('",')[0]
            print(f"mqtt gateway {gateway} | topic click id {clickid}")
    except Exception as e:
        print(e)
    time.sleep(5)
