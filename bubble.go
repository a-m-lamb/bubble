package main

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {

	driver, err := selenium.NewChromeDriverService("./chromedriver", 4444)
	if err != nil {
		panic(err)
	}
	defer driver.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"window-size=1920x1080",
	}})

	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		panic(err)
	}
	wd.DeleteAllCookies()

	if err := wd.Get("some website"); err != nil {
		fmt.Printf("Could not log into site: %s", err)
	}

	if err = login(wd); err != nil {
		fmt.Printf("Could not log into site: %s", err)
	}

	if err = search(wd); err != nil {
		fmt.Printf("Unable to submit search: %s", err)
	}

	if err = getPosts(wd); err != nil {
		fmt.Printf("Failed to get Posts: %s", err)
	}

}

func login(wd selenium.WebDriver) error {
	wd.SetImplicitWaitTimeout(60000000000 * time.Nanosecond)

	loginBtn, err := wd.FindElement(selenium.ByXPATH, "//location of login btn")
	if err != nil {
		return err
	}

	if err = loginBtn.Click(); err != nil {
		return err
	}

	email, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"email\"]")
	if err != nil {
		return err
	}

	password, err := wd.FindElement(selenium.ByXPATH, "//*[@id=\"password\"]")
	if err != nil {
		return err
	}

	if err = email.SendKeys(""); err != nil {
		return err
	}

	if err = password.SendKeys(""); err != nil {
		return err
	}

	loginBtn, err = wd.FindElement(selenium.ByXPATH, "//second login btn")

	if err = loginBtn.Click(); err != nil {
		return err
	}

	fmt.Println("We're logged in!!!")
	return nil
}

func search(wd selenium.WebDriver) error {
	wd.SetImplicitWaitTimeout(60000000000)

	searchBar, err := wd.FindElement(selenium.ByXPATH, "//searchBar")
	if err != nil {
		return err
	}

	searchBar.SendKeys("whatever you want to search\n")
	return err
}

func getPosts(wd selenium.WebDriver) error {
	var wg sync.WaitGroup
	div := 1
	var err error

	for i := 1; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			xpath := fmt.Sprintf("//post", div)
			ad := fmt.Sprintf("//ad", div)
			_, err := wd.FindElement(selenium.ByXPATH, xpath)

			for err != nil {
				div++
				xpath := fmt.Sprintf("//post", div)
				err = scrollForElement(wd, xpath, "//page")
				if err != nil {
					return
				}
				err = savePost(wd, err)
				if err != nil {
					return
				}
				fmt.Printf("\nWe have successfully saved a post %d", i)
				div++
			}

			if err == nil {
				_, err := wd.FindElement(selenium.ByXPATH, ad)
				if err != nil {
					fmt.Println("\nWe don't have an ad")
				} else {
					div++
				}

				post, err := wd.FindElement(selenium.ByXPATH, fmt.Sprintf("//post", div))

				post.Click()
				err = savePost(wd, err)
				if err != nil {
					return
				} else {
					fmt.Printf("\nWe have successfully saved a post %d", i)
					div++
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("\nI think this works")
	return err
}

func scrollForElement(wd selenium.WebDriver, parentXpath string, xpath string) error {
	var err error
	script := fmt.Sprintf(`
        var elem = document.evaluate("%s", parent, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;
		var parent = document.evaluate("%s", document, null, XPathResult.FIRST_ORDERED_NODE_TYPE, null).singleNodeValue;

        if (elem) {
            elem.click();
            return true;
        } else {
            parent.scrollTop = parent.scrollHeight;
            return false;
        }

    `, parentXpath, xpath)

	found := false
	for !found {
		if _, err := wd.ExecuteScript(script, nil); err != nil {
			return fmt.Errorf("failed to execute script: %w", err)
		}

		elem, err := wd.FindElement(selenium.ByXPATH, xpath)
		if elem != nil {
			found = true
		}
		if err != nil {
			return errors.New("Failed to find element")
		}
	}

	return err
}

func savePost(wd selenium.WebDriver, err error) error {
	saveButton := "//savebtn"

	err = scrollForElement(wd, saveButton, "//page")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	wd.Back()
	fmt.Println("Saved the element going to previous page")
	return nil // no error
}
