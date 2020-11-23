from selenium import webdriver
from selenium.common.exceptions import TimeoutException
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from time import sleep
import pickle
import json

driver = webdriver.Firefox()
wait = WebDriverWait(driver, 10)

driver.get('https://www.weibo.com/')
input_username = wait.until(EC.presence_of_element_located((By.ID, 'loginname')))
input_password = wait.until(EC.presence_of_element_located((By.NAME, 'password')))
login = wait.until(EC.element_to_be_clickable((By.XPATH, '//*[@id="pl_login_form"]/div/div[3]/div[6]/a')))
input_username.clear()
input_password.clear()
input_username.send_keys("your username")
sleep(1)
input_password.send_keys("your password")
sleep(1)
login.click()
sleep(15)

sleep(80)
with open("utils/cookies.txt","w") as f:
    for cookie in driver.get_cookies():
        f.write(str(cookie).replace("\'","\"").replace("False","false").replace("True","true")+"\n")    

driver.quit()
