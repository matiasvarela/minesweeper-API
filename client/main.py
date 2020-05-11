import requests
import time

class MinesweeperAPI:
    baseURL = 'http://ec2-3-14-1-190.us-east-2.compute.amazonaws.com:8080'

    def createGame(self, rows, columns, bombs):
        return requests.post(
            url=self.baseURL+'/games',
            json={
                "rows": rows,
                "columns": columns,
                "bombs_number": bombs
            }
        ).json()

    def getGame(self, gameId):
        return requests.get(url=self.baseURL+'/games/'+gameId).json()

    def markCell(self, gameId, row, column):
        return requests.put(
            url=self.baseURL+'/games/'+gameId+'/mark',
            json={
                "row": row,
                "column": column
            }).json()

    def RevealCell(self, gameId, row, column):
        return requests.put(
            url=self.baseURL+'/games/'+gameId+'/play-square',
            json={
                "row": row,
                "column": column
            }).json()