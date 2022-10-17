# -*- coding: utf-8 -*-

import time
import subprocess
from selenium import webdriver
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.common.by import By
from dotenv import load_dotenv
from selenium.common.exceptions import NoSuchElementException
from dotenv import load_dotenv
import os

# .envファイルの内容を読み込見込む
load_dotenv()

path = "/Users/Owner/Desktop/python_project/driver/chromedriver"

options = Options()
# options.headless = True
service = Service(path)

driver = webdriver.Chrome(options=options, service=service)

try:

    # 対象サイトに移動
    driver.get('https://trade.smbcnikko.co.jp/Login/0/login/ipan_web/hyoji/')

    time.sleep(3)

    driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/div[2]/input').send_keys(os.environ['SMBC_SHITEN_CODE'])

    driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/div[3]/input').send_keys(os.environ['SMBC_KOZA_NUMBER'])

    driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/div[4]/div/input').send_keys(os.environ['SMBC_PASS_WORD'])

    driver.find_element(by=By.XPATH, value='/html/body/div[1]/div[6]/div/div[1]/div[2]/div[1]/form/div/div[1]/p[3]/input').click()

    driver.get('https://trade.smbcnikko.co.jp/StockOrderConfirmation/1CC1H0450321/ez_ipo/meigara/ichiran')

    dict = {}

    for i in range(1, 50):

        targetImgXpath = ''
        cd = ''

        try:
            targetImgXpath = '/html/body/table[1]/tbody/tr/td[2]/div[1]/div[2]/table[2]/tbody/tr[1]/td/div[2]/table/tbody/tr/td/div[' + str(i) + ']/table[2]/tbody/tr[2]/td[2]/span/a/img'
            cd = driver.find_element(by=By.XPATH, value='/html/body/table[1]/tbody/tr/td[2]/div[1]/div[2]/table[2]/tbody/tr[1]/td/div[2]/table/tbody/tr/td/div[' + str(i) + ']/table[1]/tbody/tr[2]/td/div/table/tbody/tr/td[3]/table/tbody/tr/td/span[2]').text[-4:]
        except NoSuchElementException:
            break

        dict.setdefault(cd, targetImgXpath)

    target = '5132'

    driver.find_element(by=By.XPATH, value=dict[target]).click()

    driver.find_element(by=By.XPATH, value='//*[@id="mcChk"]').click()

    driver.find_element(by=By.XPATH, value='//*[@id="printzone"]/form/div/table/tbody/tr/td/div[4]/table/tbody/tr[1]/td/input[3]').click()

    driver.find_element(by=By.XPATH, value='//*[@id="printzone"]/div[2]/form/div[1]/table/tbody/tr[4]/td/div[1]/div[3]/table/tbody/tr[3]/td[2]/table/tbody/tr/td[1]/input').send_keys("100")

    driver.find_element(by=By.XPATH, value='//*[@id="printzone"]/div[2]/form/div[1]/table/tbody/tr[4]/td/div[1]/div[3]/table/tbody/tr[5]/td[2]/span/select/option[2]').click()

    driver.find_element(by=By.XPATH, value='//*[@id="printzone"]/div[2]/form/div[1]/table/tbody/tr[4]/td/div[2]/table/tbody/tr/td/table/tbody/tr[1]/td/input[1]').click()

    driver.find_element(by=By.XPATH, value='//*[@id="printzone"]/div[2]/table/tbody/tr/td/div[3]/table/tbody/tr[4]/td/div[2]/table/tbody/tr/td/table/tbody/tr[1]/td/form/input[1]').click()

    driver.find_element(by=By.XPATH, value='//*[@id="printzone"]/table/tbody/tr/td/div/table/tbody/tr/td[1]/a/div').click()

    print('成功')
except:

    print('失敗')
    
    driver.close()