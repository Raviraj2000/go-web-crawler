import os
import json

def count_unique_urls(output_folder='./output'):
    urls = list()

    # Iterate through all files in the output directory
    for filename in os.listdir(output_folder):
        if filename.endswith('.json'):
            filepath = os.path.join(output_folder, filename)
            with open(filepath, 'r', encoding='utf-8') as file:
                for line in file:
                    try:
                        data = json.loads(line.strip())
                        url = data.get("URL", "").strip()
                        if url:
                            urls.append(url)
                        else:
                            print(f"Empty or missing URL in file {filename}: {data}")  # Debugging line
                    except json.JSONDecodeError:
                        print(f"Error decoding JSON in file {filename}")

    # Output the count of unique URLs
    print(f"Total URLs across all files: {len(urls)}")
    print(f"Total unique across all files: {len(set(urls))}")
    return len(urls)

if __name__ == "__main__":
    count_unique_urls()
