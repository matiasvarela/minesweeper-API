import requests
import time

class MinesweeperAPI:
    baseURL = 'http://ec2-3-14-1-190.us-east-2.compute.amazonaws.com:8080'

    def createGame(self, userId, rows, columns, bombs):
        return requests.post(
            url=self.baseURL+'/users/'+str(userId)+'/games',
            json={
                "rows": rows,
                "columns": columns,
                "bombs_number": bombs
            }
        ).json()

    def getGame(self, userId, gameId):
        return requests.get(url=self.baseURL+'/users/'+str(userId)+'/games/'+str(gameId)).json()

    def markCell(self, userId, gameId, row, column):
        return requests.put(
            url=self.baseURL+'/users/'+str(userId)+'/games/'+str(gameId)+'/mark',
            json={
                "row": row,
                "column": column
            }).json()

    def RevealCell(self, userId, gameId, row, column):
        return requests.put(
            url=self.baseURL+'/users/+'+str(userId)+'+/games/'+str(gameId)+'/play-square',
            json={
                "row": row,
                "column": column
            }).json()