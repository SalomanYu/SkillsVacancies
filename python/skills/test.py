import re
from collections import Counter

with open("bot_log.txt", "r", encoding="utf-8") as f:
    f = f.readlines()

regexp = "INFO \[(.*?)\]"
names = []
for i in f:
    name = re.findall(regexp, i)
    if name: names.append(name[0])


count = Counter()
for item in names:
    count[item] += 1
from pprint import pprint

pprint(sorted(count, key=count.get, reverse=True))