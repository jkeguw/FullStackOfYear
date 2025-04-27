#!/usr/bin/env python3
# eloshape_scraper.py - A script to scrape mouse data from eloshapes.com

import requests
from bs4 import BeautifulSoup
import pandas as pd
import json
import time
from datetime import datetime

def scrape_eloshapes():
    """
    Scrapes mouse data from eloshapes.com and saves it to CSV and JSON formats.
    """
    # URL to scrape
    url = "https://www.eloshapes.com/mouse/database"
    
    # Send HTTP request with headers to mimic a browser
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Connection': 'keep-alive',
    }
    
    try:
        print(f"Fetching data from {url}...")
        response = requests.get(url, headers=headers)
        response.raise_for_status()  # Raise an exception for bad status codes
        
        # Parse HTML content
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # Extract the table with mouse data
        table = soup.find('table', {'class': 'table'})
        
        if not table:
            print("Table not found on the page.")
            return None
        
        # Extract table headers
        headers = []
        for th in table.find('thead').find_all('th'):
            headers.append(th.text.strip())
        
        # Extract table rows
        rows = []
        for tr in table.find('tbody').find_all('tr'):
            row = []
            for td in tr.find_all('td'):
                # Extract text, handle empty cells
                cell_text = td.text.strip() if td.text.strip() else None
                row.append(cell_text)
            rows.append(row)
        
        # Create a pandas DataFrame
        df = pd.DataFrame(rows, columns=headers)
        
        # Save to CSV
        csv_filename = f"eloshapes_mice_{datetime.now().strftime('%Y%m%d_%H%M%S')}.csv"
        df.to_csv(csv_filename, index=False)
        print(f"Data saved to {csv_filename}")
        
        # Save to JSON
        json_filename = f"eloshapes_mice_{datetime.now().strftime('%Y%m%d_%H%M%S')}.json"
        df.to_json(json_filename, orient='records')
        print(f"Data saved to {json_filename}")
        
        # Print summary
        print(f"Successfully scraped data for {len(df)} mice.")
        return df
        
    except requests.exceptions.RequestException as e:
        print(f"Error fetching the URL: {e}")
    except Exception as e:
        print(f"Error processing data: {e}")
    
    return None

if __name__ == "__main__":
    print("Starting EloShapes mouse database scraper...")
    scrape_eloshapes()
    print("Scraping completed.")
    time.sleep(10)

    