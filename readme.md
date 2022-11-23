# Right Move-Scraper
## Parse The Data
This can be used to scrape the sold houses data from Right Move
1. Go to https://www.rightmove.co.uk/house-prices/kt5/king-charles-crescent.html?soldIn=2&radius=0.5&page=1;
1. Amend the data as you wish and search for it;
1. Manually download all the result pages;
1. import them into here.

## Use Linear Regression

### Setup
```bash
cd . # cd into the root of this project

virtualenv -p python3 . # Setup virtual env

source bin/activate # Activate virtual env

pip3 install -r requirements.txt # Install dependencies

python3 main.py # Run Linear Regression
```
