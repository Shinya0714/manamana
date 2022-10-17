# -*- coding: utf-8 -*-

import time
import subprocess
from selenium import webdriver
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from dotenv import load_dotenv
import os

# .envファイルの内容を読み込見込む
load_dotenv()

path = "/Users/Owner/Desktop/python_project/driver/chromedriver"

options = Options()
options.headless = True
service = Service(path)

driver = webdriver.Chrome(options=options, service=service)

# 対象サイトに移動
driver.get('https://trade.smbcnikko.co.jp/Login/0/login/ipan_web/hyoji/')

time.sleep(3)

driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/div[2]/input').send_keys(os.environ['SMBC_SHITEN_CODE'])

driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/div[3]/input').send_keys(os.environ['SMBC_KOZA_NUMBER'])

driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/div[4]/div/input').send_keys(os.environ['SMBC_PASS_WORD'])

driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/p[3]/input').click()

print(driver.find_element(by=By.XPATH, value='/html/body/div[3]/div[1]/table/tbody/tr[1]/td/div/div[2]/table/tbody/tr/td[1]/div[1]/table/tbody/tr/td/div[2]/table/tbody/tr[1]/td/table/tbody/tr/td[2]/div/div[1]').text)
