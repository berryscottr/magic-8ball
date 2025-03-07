import asyncio
import aiohttp
import argparse
import json
import typing

class GraphQLClient:
    def __init__(self, base_url: str, headers: dict):
        self.base_url = base_url
        self.headers = headers
    
    async def query(self, session: aiohttp.ClientSession, query: str, variables: dict, operation_name: str) -> typing.List[typing.Dict]:
        json_data = {
            "query": query,
            "variables": variables,
            "operationName": operation_name,
        }
        async with session.post(self.base_url, json=json_data) as resp:
            return await resp.json()

class PoolPlayersAPI:
    def __init__(self, client: GraphQLClient):
        self.client = client
    
    async def get_team(self, session: aiohttp.ClientSession, team_id: int):
        query = """
        query teamPage($id: Int!) {
          team(id: $id) {
            id
            name
            number
            isTied
            standing
            division {
              id
              name
              number
              timeOfPlay
              nightOfPlay
              format
              state
              isTournament
              tournament {
                id
                name
                __typename
              }
              __typename
            }
            location {
              id
              name
              address {
                id
                name
                __typename
              }
              __typename
            }
            session {
              id
              name
              __typename
            }
            league {
              id
              slug
              __typename
            }
            __typename
          }
        }
        """
        return await self.client.query(session, query, {"id": team_id}, "teamPage")
    
    async def get_match(self, session: aiohttp.ClientSession, match_id: int):
        query = """
        query MatchPage($id: Int!) {
          match(id: $id) {
            id
            division {
              id
              electronicScoringEnabled
              __typename
            }
            league {
              id
              esEnabled
              __typename
            }
            ...matchForCart
            __typename
          }
        }
 
        fragment matchForCart on Match {
          __typename
          id
          type
          startTime
          week
          isBye
          isMine
          isScored
          scoresheet
          isPaid
          location {
            ...googleMapComponent
            __typename
          }
          home {
            id
            name
            number
            isMine
            ...rosterComponent
            __typename
          }
          away {
            id
            name
            number
            isMine
            ...rosterComponent
            __typename
          }
          division {
            id
            scheduleInEdit
            type
            __typename
          }
          session {
            id
            name
            year
            __typename
          }
          league {
            id
            name
            currentSessionId
            isElectronicPaymentsEnabled
            country {
              id
              __typename
            }
            __typename
          }
          fees {
            amount
            tax
            total
            __typename
          }
          orderItems {
            id
            order {
              id
              member {
                id
                firstName
                lastName
                __typename
              }
              __typename
            }
            __typename
          }
          results {
            homeAway
            overUnder
            forfeits
            matchesWon
            matchesPlayed
            points {
              bonus
              penalty
              won
              adjustment
              sportsmanship
              total
              skillLevelViolationAdjustment
              __typename
            }
            scores {
              id
              player {
                id
                displayName
                __typename
              }
              matchPositionNumber
              playerPosition
              skillLevel
              innings
              defensiveShots
              eightBallWins
              eightOnBreak
              eightBallBreakAndRun
              nineBallPoints
              nineOnSnap
              nineBallBreakAndRun
              nineBallMatchPointsEarned
              mastersEightBallWins
              mastersNineBallWins
              winLoss
              matchForfeited
              doublesMatch
              dateTimeStamp
              teamSlot
              eightBallMatchPointsEarned
              incompleteMatch
              __typename
            }
            __typename
          }
        }
 
        fragment googleMapComponent on HostLocation {
          id
          phone
          name
          address {
            id
            name
            address1
            address2
            city
            zip
            latitude
            longitude
            __typename
          }
          __typename
        }
 
        fragment rosterComponent on Team {
          id
          name
          number
          league {
            id
            slug
            __typename
          }
          division {
            id
            type
            __typename
          }
          roster {
            id
            memberNumber
            displayName
            matchesWon
            matchesPlayed
            ... on EightBallPlayer {
              pa
              ppm
              skillLevel
              __typename
            }
            ... on NineBallPlayer {
              pa
              ppm
              skillLevel
              __typename
            }
            member {
              id
              __typename
            }
            __typename
          }
          __typename
        }
        """
        return await self.client.query(session, query, {"id": match_id}, "MatchPage")

def create_parser():
    parser = argparse.ArgumentParser(description="Fetch match or team details.")
    parser.add_argument("-t", "--team-id", type=int, required=False, help="Team ID to query.")
    parser.add_argument("-m", "--match-id", type=int, required=False, help="Match ID to query.")
    return parser

async def main(args):
    headers = {
        "user-agent": "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.104 Safari/537.36",
        "referer": "https://league.poolplayers.com",
        "Authorization": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Im9seWR1dTRod0xkdHBJTEFsRGVtZi1YYjZNWjhueWg2bC1UdzN6bTRZQnMifQ.eyJhcHBsaWNhdGlvblJlZnJlc2hUb2tlbklkIjoiNTUzMTI3NSIsImlhdCI6MTc0MTM4MTk1NSwiZXhwIjoxNzQxMzgyODU1LCJpc3MiOiJBUEEiLCJzdWIiOiIyMDI2MDQzIn0.WW1kQqIGAYx10n1KzhW-xYiTCZzW2c3dkbZaTWNO-VfsrTNuFVax1fBugtg9KAewoJg4i5BseWHsvGa6yMFQSz-0kvvEu4d2I3ClpB3xb8xLk9VVr09gGqN9yneSgM9eRb1c6Kl5WEiLfzG102aQ_WFgKpa1LpQdCsYmDYnMvqmOdpKfQngkbt5p1SaVDUA7KrSofVb0da2VIpSligFQA4qFEruiyIT4SOErNtGeye8bPJ-MLR_eKlMck4RoXhdUA4Cm-0-PZWtXwag1BpWPSHK_QXNSFLUtP_B6mgSe2XQpkfl8RrG2VJMwWQDKbkeAXq-XRrG_H9OihPeGyDtWMQ",
        "authority": "gql.poolplayers.com",
    }
    
    graphql_url = "https://gql.poolplayers.com/graphql"
    client = GraphQLClient(graphql_url, headers)
    api = PoolPlayersAPI(client)
    
    async with aiohttp.ClientSession(headers=headers) as session:
        if args.team_id:
            result = await api.get_team(session, args.team_id)
        elif args.match_id:
            result = await api.get_match(session, args.match_id)
        print(json.dumps(result, indent=2))

if __name__ == "__main__":
    asyncio.run(main(create_parser().parse_args()))
