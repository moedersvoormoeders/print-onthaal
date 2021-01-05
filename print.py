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
            
            totaalPrijs = 0.0

            p.set(width=3, height=3)
            p.text(content['klant']['mvmNummer']+"\n")
            p.set(width=2, height=2)
            p.text(content['klant']['voornaam']+" "+content['klant']['naam']+"\n")
            p.text("\n")
            p.text("Materiaal\n")
            
            for item in content['items']:
                p.text("\n")
                p.text("----------------")
                p.text("\n")
                p.text(item['object'] + "\n")
                if 'prijs' in item and item['prijs'] is not None:
                    totaalPrijs += item['prijs']
                
                if 'ontvanger' in item and item['ontvanger'] is not None:
                    p.text(item['ontvanger']["naam"] + "\n")
                    p.text(item['ontvanger']["geslacht"] + "\n")

                if 'maat' in item and item['maat'] is not None:
                    p.text("Maat: "+ item['maat'] + "\n")
                
                p.text(item['opmerking'] + "\n")
                
                p.text("\n")
                p.text("----------------")
                p.text("\n")
            p.text("Totaal " + str(totaalPrijs) + " EUR")
            p.cut()
            return {'status': 'ok'}
        except:
            return {'status': 'error'}
        finally:
            mutex.release()
            pass

class Sinterklaas(Resource):
    def post(self):
        mutex.acquire()
        try:
            content = request.json
            if content == None:
                return {'status': 'error', 'error': 'No Content'}

            # Speelgoed

            p.set(width=3, height=3)
            p.text(content['speelgoed']['mvmNummer']+"\n")
            p.set(width=2, height=2)
            p.text(content['speelgoed']['naam']+"\n")
            p.text("\n")
            p.text("Sinterklaas\n")
            
            for item in content['speelgoed']['paketten']:
                p.text("\n")
                p.text("----------------")
                p.text("\n")
                p.text(item['naam'] + "\n")
                p.text(item["geslacht"] + "\n")
                p.text(str(item["leeftijd"]) + " jaar\n")
                p.text(item['opmerking'] + "\n")
                p.text("\n")
                p.text("----------------")
                p.text("\n")
            p.cut()


            # Snoep
            p.set(width=3, height=3)
            p.text(content['snoep']['mvmNummer']+"\n")
            p.set(width=2, height=2)
            p.text(content['snoep']['naam']+"\n")
            p.text("\n")
            p.text("Sinterklaas Snoep\n")
            p.text("\n")
            p.text("volwassenen: " + str(content['snoep']['volwassenen'])+"\n")
            p.text("kinderen: " + str(content['snoep']['kinderen'])+"\n")
            p.text("\n")
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
api.add_resource(Sinterklaas, '/sinterklaas')
api.add_resource(Eenmaligen, '/eenmaligen')

if __name__ == '__main__':
     app.run(port='8080')
