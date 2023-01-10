from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By


def main():
    driver = webdriver.Firefox()
    driver.get("https://league.poolplayers.com")
    login = driver.find_element(By.ID, "__next")
    login.clear()
    login.send_keys("berryscottr@gmail.com")
    login.send_keys(Keys.TAB)
    login.clear()
    login.send_keys("")
    login.send_keys(Keys.ENTER)
    enter = driver.find_element(By.CLASS_NAME, "btn-primary")
    enter.click()
    driver.close()


if __name__ == '__main__':
    main()
