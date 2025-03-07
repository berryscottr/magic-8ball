import requests
import json
import os
from dotenv import load_dotenv


class AuthManager:
    def __init__(self, username, password, login_url, headers):
        self.username = username
        self.password = password
        self.login_url = login_url
        self.headers = headers
        self.session = requests.Session()
        self.device_refresh_token = None
        self.access_token = None

    def login(self):
        """Login to get the device refresh token."""
        login_payload = {
            "operationName": "login",
            "query": """
            mutation login($username: String!, $password: String!) {
              login(input: {username: $username, password: $password}) {
                __typename
                ... on SuccessLoginPayload {
                  deviceRefreshToken
                }
                ... on PartialSuspendedLoginPayload {
                  deviceRefreshToken
                }
                ... on DeniedLoginPayload {
                  reason
                }
              }
            }
            """,
            "variables": {
                "username": self.username,
                "password": self.password
            }
        }
        login_response = self.session.post(self.login_url, json=login_payload, headers=self.headers)

        if login_response.status_code == 200:
            login_data = login_response.json()
            if "errors" in login_data:
                print("Login failed:", login_data["errors"])
                return None

            self.device_refresh_token = login_data["data"]["login"].get("deviceRefreshToken", None)
            if not self.device_refresh_token:
                print("No device refresh token received. Please check the login response.")
                return None

            print("Device Refresh Token:", self.device_refresh_token)
            return self.device_refresh_token
        else:
            print("Login request failed:", login_response.text)
            return None

    def refresh_access_token(self):
        """Use the device refresh token to get an access token."""
        if not self.device_refresh_token:
            print("No device refresh token available. Please login first.")
            return None

        refresh_payload = {
            "operationName": "GenerateAccessTokenMutation",
            "query": """
            mutation GenerateAccessTokenMutation($refreshToken: String!) {
            generateAccessToken(refreshToken: $refreshToken) {
                accessToken
                __typename
            }
            }
            """,
            "variables": {
                "refreshToken": self.device_refresh_token
            }
        }

        refresh_response = self.session.post(self.login_url, json=refresh_payload, headers=self.headers)

        if refresh_response.status_code == 200:
            refresh_data = refresh_response.json()
            if "errors" in refresh_data:
                print("Token refresh failed:", refresh_data["errors"])
                return None

            self.access_token = refresh_data["data"]["generateAccessToken"].get("accessToken", None)
            if not self.access_token:
                print("No access token received. Please check the refresh response.")
                return None

            print("Access Token:", self.access_token)

            # Log the refresh token to see if it was updated during this request
            if "deviceRefreshToken" in refresh_data:
                print(f"New device refresh token: {refresh_data['deviceRefreshToken']}")
                self.device_refresh_token = refresh_data["deviceRefreshToken"]

            return self.access_token
        else:
            print("Token refresh request failed:", refresh_response.text)
            return None


class APIClient:
    def __init__(self, auth_manager):
        self.auth_manager = auth_manager

    def execute_query(self, query, variables):
        """Execute a GraphQL query using the access token."""
        access_token = self.auth_manager.access_token
        if not access_token:
            access_token = self.auth_manager.refresh_access_token()
        
        if not access_token:
            print("Failed to get an access token.")
            return None
        
        headers = {**self.auth_manager.headers, "Authorization": f"Bearer {access_token}"}
        response = self.auth_manager.session.post(
            self.auth_manager.login_url,
            headers=headers,
            json={"query": query, "variables": variables}
        )
        
        try:
            response_data = response.json()
            return response_data
        except json.JSONDecodeError:
            print("Failed to parse response:", response.text)
            return None


class GraphQLQuery:
    @staticmethod
    def get_match_query():
        return """
        query MatchPage($id: Int!) {
          match(id: $id) {
            id
            type
            startTime
          }
        }
        """


# Main logic
if __name__ == "__main__":
    load_dotenv()
    # Define the necessary variables
    LOGIN_URL = "https://gql.poolplayers.com/graphql/"
    USERNAME = os.getenv("USERNAME")
    PASSWORD = os.getenv("PASSWORD")
    
    HEADERS = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:136.0) Gecko/20100101 Firefox/136.0",
        "Accept": "*/*",
        "Accept-Language": "en-US,en;q=0.5",
        "Accept-Encoding": "gzip, deflate, br, zstd",
        "Referer": "https://league.poolplayers.com/",
        "Content-Type": "application/json",
        "apollographql-client-name": "MemberServices",
        "apollographql-client-version": "3.18.33-3463",
        "Origin": "https://league.poolplayers.com",
        "DNT": "1",
        "Sec-GPC": "1",
        "Connection": "keep-alive",
        "Sec-Fetch-Dest": "empty",
        "Sec-Fetch-Mode": "cors",
        "Sec-Fetch-Site": "same-site",
        "Priority": "u=4",
        "TE": "trailers"
    }

    # Initialize the AuthManager with your credentials
    auth_manager = AuthManager(USERNAME, PASSWORD, LOGIN_URL, HEADERS)

    # Step 1: Login and get the device refresh token
    if not auth_manager.login():
        print("Login failed. Exiting.")
        exit(1)

    # Step 2: Get the access token using the device refresh token
    if not auth_manager.refresh_access_token():
        print("Failed to refresh access token. Exiting.")
        exit(1)

    # Step 3: Initialize the APIClient with the AuthManager
    api_client = APIClient(auth_manager)

    # Step 4: Use the access token to make a query
    query = GraphQLQuery.get_match_query()
    variables = {"id": 45117061}
    
    response_data = api_client.execute_query(query, variables)

    if response_data:
        print(json.dumps(response_data, indent=2))
    else:
        print("No data received.")
