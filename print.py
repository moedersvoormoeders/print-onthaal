from escpos.printer import Usb, Dummy

from flask import Flask, request
from flask_restful import Resource, Api
from json import dumps
from flask_jsonpify import jsonify
from flask_cors import CORS
import json

from threading import Lock

mutex = Lock()

app = Flask(__name__)
CORS(app)
api = Api(app)

p = Usb(0x04b8, 0x0202, 0)

class Print(Resource):
    def post(self):
        mutex.acquire()
        try:
            content = request.json
            if content == None:
                return {'status': 'error', 'error': 'No Content'}

            p.set(width=3, height=3)
            p.text(content['doelgroepnummer']+"\n")
            p.set(width=2, height=2)
            p.text("kleding\n")
            p.text(content['naam'] + " " + content['voornaam']+"\n")
            p.text(content['kind'] + "\n")
            # p.text(str(content['leeftijd']) + " jaar\n")
            p.cut()
            return {'status': 'ok'}
        except:
            return {'status': 'error'}
        finally:
            mutex.release()
            pass

class Eenmaligen(Resource):
    def post(self):
        mutex.acquire()
        try:
            content = request.json
            if content == None:
                return {'status': 'error', 'error': 'No Content'}

            p.set(width=3, height=3)
            p.text(content['eenmaligenNummer']+"\n")
            p.set(width=2, height=2)
            if content['naam'] != None:
                p.text(content['naam'] +"\n")
            if content['bericht'] != None:
                p.text(content['bericht'] +"\n")
            p.cut()
            return {'status': 'ok'}
        except:
            return {'status': 'error'}
        finally:
            mutex.release()
            pass



api.add_resource(Print, '/print')
api.add_resource(Eenmaligen, '/eenmaligen')

if __name__ == '__main__':
     app.run(port='8080')
