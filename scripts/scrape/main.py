import os
from playwright.sync_api import sync_playwright

LOGIN_URL = "https://accounts.poolplayers.com/login"
GRAPHQL_DOMAIN = "gql.poolplayers.com"

def run():
    email = os.getenv("APA_EMAIL")
    password = os.getenv("APA_PASSWORD")

    if not email or not password:
        print("Missing APA_EMAIL or APA_PASSWORD environment variables")
        return

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context()
        page = context.new_page()

        token_box = {"token": None}

        # capture the token
        def on_request(req):
            if GRAPHQL_DOMAIN in req.url:
                auth = req.headers.get("authorization")
                if auth and token_box["token"] is None:
                    token_box["token"] = auth

        page.on("request", on_request)

        # perform login
        page.goto(LOGIN_URL)
        page.fill("input[name='email']", email)
        page.fill("input[name='password']", password)
        page.keyboard.press("Enter")

        page.wait_for_selector("button:has-text('Continue')", timeout=20000)
        page.click("button:has-text('Continue')")
        page.wait_for_url("**/dashboard", timeout=30000)

        # Poll until token arrives or timeout passes
        for _ in range(300):   # 300 × 100 ms = 30 seconds max
            if token_box["token"]:
                print(token_box["token"])
                browser.close()
                return
            page.wait_for_timeout(100)

        # If no token found:
        browser.close()

if __name__ == "__main__":
    run()