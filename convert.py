from pathlib import Path
import json
import sys
import argparse

def find_extension_jsons(root_dir):
    """
    Recursively find all 'extension.json' files in the given directory.
    
    Args:
        root_dir (str): Root directory to start the search from
        
    Returns:
        list: List of Path objects pointing to all found extension.json files
    """
    root_path = Path(root_dir)
    return list(root_path.rglob("extension.json"))

def validate_directory(dir_path):
    """
    Validate that the path exists and is a directory.
    
    Args:
        dir_path (str): Path to validate
        
    Returns:
        Path: Validated Path object
        
    Raises:
        ArgumentTypeError: If path doesn't exist or isn't a directory
    """
    path = Path(dir_path)
    if not path.exists():
        raise argparse.ArgumentTypeError(f"Path does not exist: {dir_path}")
    if not path.is_dir():
        raise argparse.ArgumentTypeError(f"Path is not a directory: {dir_path}")
    return path

def main():
    parser = argparse.ArgumentParser(description='Convert old extension.json format to new format')
    parser.add_argument('search_dir', type=validate_directory, 
                       help='Directory to recursively search for extension.json files')
    parser.add_argument('out_dir', type=validate_directory,
                       help='Directory where the converted json will be written')
    
    args = parser.parse_args()
    
    extension_files = find_extension_jsons(args.search_dir)
    
    # Print found files for testing
    for file_path in extension_files:
        print(f"Found: {file_path}")
    
    # TODO: Add JSON conversion logic here
    """
    combined_data = {}
    for file_path in extension_files:
        with file_path.open() as f:
            data = json.load(f)
            # Process and merge data as needed
            
    # Write the result
    output_path = args.out_dir / "converted_extensions.json"
    with output_path.open('w') as f:
        json.dump(combined_data, f, indent=2)
    """

if __name__ == "__main__":
    main()