import json
import sys
import pandas
import numpy as np
from sklearn.preprocessing import PolynomialFeatures
from sklearn.linear_model import LinearRegression

def predict(dimension, linear_reg, poly, pc, nb, ht, st, year, dht, dpc, dt):
    predictedHousePrice = linear_reg.predict(poly.fit_transform([[int(pc), int(nb), int(ht), int(st), int(year)]]))

    print("{} prediction for postcode:{} #bedrooms:{} type:{} tenure:{} year:{} £{}".format(
        dimension,
        dpc[str(pc)],
        nb,
        dht[str(ht)],
        dt[str(st)],
        year,
        round(predictedHousePrice[0], 2),
    ))

def loadDecoder():
    with open('encodedHouseType.json') as f:
        lines = f.read()
        eht = json.loads(lines)
        
    with open('encodedPostcode.json') as f:
        lines = f.read()
        epc = json.loads(lines)

    with open('encodedTenure.json') as f:
        lines = f.read()
        et = json.loads(lines)

    with open('decodeHouseType.json') as f:
        lines = f.read()
        dht = json.loads(lines)
        
    with open('decodePostcode.json') as f:
        lines = f.read()
        dpc = json.loads(lines)

    with open('decodeTenure.json') as f:
        lines = f.read()
        dt = json.loads(lines)
    
    return eht, epc, et, dht, dpc, dt

def learn(degree, X, y, pc, nb, ht, st, year, dht, dpc, dt):
    poly = PolynomialFeatures(degree=degree)
    X1 = poly.fit_transform(X)
    linear_reg = LinearRegression()
    linear_reg.fit(X1, y)
    predict('{}d'.format(degree), linear_reg, poly, pc, nb, ht, st, year, dht, dpc, dt)

def main():
    args = sys.argv[1:len(sys.argv)]
    if len(args) != 5:
        print('make sure postcode, num bedroom, type, tenure and year')
        sys.exit(1)

    eht, epc, et, dht, dpc, dt = loadDecoder()

    df = pandas.read_csv("encoded-house-data.csv")

    X = df[['postcode_first_part', 'num_bedrooms', 'house_type', 'sale_tenure', 'sale_year']].values
    y = df['sale_amount'].values

    learn(1, X, y, epc[args[0]], args[1], eht[args[2]], et[args[3]], args[4], dht, dpc, dt)
    learn(2, X, y, epc[args[0]], args[1], eht[args[2]], et[args[3]], args[4], dht, dpc, dt)
    learn(3, X, y, epc[args[0]], args[1], eht[args[2]], et[args[3]], args[4], dht, dpc, dt)
    learn(4, X, y, epc[args[0]], args[1], eht[args[2]], et[args[3]], args[4], dht, dpc, dt)

if __name__ == "__main__":
    main()
