import requests
import bs4
import argparse
import json

parser = argparse.ArgumentParser()

parser.add_argument('-u', '--url', required=True)
parser.add_argument('-o', '--output', default='output.json')

args = parser.parse_args()

r = requests.get(args.url)
soup = bs4.BeautifulSoup(r.content, features='html.parser')

text = soup.get_text().split()
words = set()

for word in text:
    valid = True
    for char in word:
        if ord(char) < 97 or ord(char) > 122:
            valid = False
            break
    if valid:
        words.add(word)
        
with open(args.output, 'w+') as f:
    json.dump(list(words), f)