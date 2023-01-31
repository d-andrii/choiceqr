package main

import "log"

func GetOnlineSum(c *Client, date string) float64 {
	checks, err := c.GetChecksForDay(date)
	if err != nil {
		log.Fatal(err)
	}

	online := 0
	for _, x := range checks {
		if (x.Status == "closed" || x.Status == "pre_closed") && x.PayBy == "online" {
			online += x.Total
		}
	}

	return float64(online) / 100
}
