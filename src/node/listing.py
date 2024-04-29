# listing.py
import os
import fnmatch

def find_source_files(directory):
    matches = []
    for root, dirnames, filenames in os.walk(directory):
        dirs_to_skip = ['.bundle', 'android', 'scripts', 'ios', 'certs', 'assets', 'libs', 'localization', 'esm', 'lib', 'node_modules', '.next', '__generated__', '.git', '.idea', '.vscode', 'build', 'dist', 'coverage', 'public', 'out', 'tmp', 'temp']
        if any(skip_dir in root.split(os.sep) for skip_dir in dirs_to_skip):
            continue
     
        for pattern in ['*.go', '*.mod', '*.ts', '*.tsx', '*.js', '*.jsx', '*.json', '*.md', '*.yml', '*.yaml', '*.env']:
            for filename in fnmatch.filter(filenames, pattern):
                if filename in ['blockchain_data.json', 'listing.py','package-lock.json'] or fnmatch.fnmatch(filename, '*.test.*'):  # Skip these files
                    continue
                matches.append(os.path.join(root, filename))
    return matches

def print_file_listings(directory):
    source_files = find_source_files(directory)
    with open('results.txt', 'w', encoding='utf-8') as results_file:
        for file_path in source_files:
            file_name = os.path.basename(file_path)
            results_file.write(f"File Path: {file_path}\n")
            with open(file_path, 'r', encoding='utf-8') as file:
                source_code = file.read()
                results_file.write(f"Source Code:\n{source_code}\n\n")

# Warning message to inform the user
print("Check result on file 'results.txt'.")

# Specify the directory of your project here
project_directory = '.'

print_file_listings(project_directory)