import json
import requests
import argparse

def add_movie_titles(json_file_path, api_key):
    
    with open(json_file_path, 'r') as file:
        data = json.load(file)
    
    for movie in data['ids']:
        movie_id = movie['movieID']
        url = f"http://www.omdbapi.com/?i={movie_id}&apikey={api_key}"
        
        response = requests.get(url)
        if response.status_code == 200:
            movie_data = response.json()
            movie['title'] = movie_data.get('Title', 'N/A')
        else:
            movie['title'] = 'Error fetching title'
    
    
    with open(json_file_path, 'w') as file:
        json.dump(data, file, indent=4)

    print("Movie titles have been added to the JSON file.")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Add movie titles to JSON file using OMDb API")
    parser.add_argument("json_file", help="Path to the JSON file containing movie IDs")
    parser.add_argument("api_key", help="OMDb API key")
    
    args = parser.parse_args()
    
    add_movie_titles(args.json_file, args.api_key)
