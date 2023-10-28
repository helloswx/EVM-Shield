from src.state_extraction.state_extractor import extract_contract_state, extract_regular_variables
from src.key_approx_analysis.key_approx_analyzer import get_slot_details
from configparser import ConfigParser
import json

def read_source_code(contract_name, input_dir):
    input_path = input_dir + contract_name + ".sol"
    f = open(input_path)
    source_code = f.read()
    return source_code


if __name__ == "__main__":
    config = ConfigParser()
    config.read("config.ini")
    input_dir = config.get('directories', 'contract_directory')
    with open(input_dir+"/contracts.json") as f:
        contracts = json.load(f)
    for index in range(len(contracts)):
        print(index+1, " ", contracts[index]['Address'], " ", contracts[index]['Contract Name'])
    num = input("\nSelect contract no from above to run SmartMuv -> ")
    contract_name = contracts[int(num)-1]['Contract Name']
    cont_addr = contracts[int(num)-1]['Address']
    compiler_version = contracts[int(num)-1]['Compiler Version']
    source_code = read_source_code(contract_name, input_dir)
    print("\n-----------------------------------------------------------\n")
    print("Select the SmartMuv feature you want to use (1-3):")
    options = """
    1 - Get Variable Slot Layout Details.
    
    Select Option - """
    option = input(options)
    print("\n-----------------------------------------------------------\n")

    if option == "1":
        slots_details = get_slot_details(contract_name, source_code, compiler_version)
        print("\nContract Slot Layout is:\n")
        for slot in slots_details:
            print(slot)
