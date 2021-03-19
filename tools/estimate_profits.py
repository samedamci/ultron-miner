#!/usr/bin/python3
"""
Simple script to calculate estimated miner profits.

Usage:
    python3 estimate_profits.py <username1> <username2> <...>
"""

import requests
import json
import time
import sys

ADDR = "http://51.15.127.80/"
DELAY = 30  # seconds


def main():
    del sys.argv[0]
    usernames = sys.argv
    earnings = 0
    last_balances = {}
    for r in range(2):
        balances = json.loads(requests.get(f"{ADDR}/balances.json").text)
        for user in usernames:
            current_balance = float(balances[user].split("")[0])
            if r == 1:
                earning = current_balance - last_balances[user]
                earnings += earning
            last_balances[user] = current_balance
        if r == 0:
            time.sleep(30)
    per_minute = round(earnings * 2, 4)
    per_hour = round(per_minute * 120, 4)
    per_day = round(per_hour * 24, 4)
    print(
        f"""Earnings per minute: {per_minute} DUCO
Earnings per hour: {per_hour} DUCO
Earnings per day: {per_day} DUCO"""
    )


if __name__ == "__main__":
    main()
