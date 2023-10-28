import json
import os
import csv
from tqdm import tqdm


def extract_variable_relationship(json_data, state_variables):
    function_variable_map = []

    def is_transfer_capable(node):
        if 'call' in str(node) or 'transfer' in str(node) or 'send' in str(node):
            return True
        return False

    def traverse(node):
        nonlocal function_variable_map
        if isinstance(node, dict):
            if node.get('type') == 'FunctionDefinition':
                function_visibility = node.get('visibility', '')
                if function_visibility in ['public', 'external']:
                    function_name = node.get('name')
                    function_params = node.get('parameters', {}).get('parameters', [])
                    params_str = ",".join([p.get('typeName', {}).get('name', '') for p in function_params])
                    full_function_name = f"{function_name}({params_str})"
                    function_body = node.get('body', {})
                    variable_references = extract_variable_references(function_body)
                    function_state_variables = [var for var in variable_references if var in state_variables]

                    can_transfer = is_transfer_capable(function_body)

                    function_variable_map.append({
                        'function_name': full_function_name,
                        'variables': function_state_variables,
                        'is_transfer_capable': can_transfer
                    })

            for value in node.values():
                traverse(value)
        elif isinstance(node, list):
            for item in node:
                traverse(item)

    def extract_variable_references(node):
        references = []

        if isinstance(node, dict):
            if node.get('type') == 'Identifier':
                reference_name = node.get('name')
                references.append(reference_name)

            for value in node.values():
                references.extend(extract_variable_references(value))
        elif isinstance(node, list):
            for item in node:
                references.extend(extract_variable_references(item))

        return references

    traverse(json_data)
    return function_variable_map


def extract_state_variables(json_data):
    state_vars = {}

    def traverse(node):
        nonlocal state_vars
        if isinstance(node, dict):
            if node.get('type') == 'StateVariableDeclaration':
                for var_dec in node.get('variables', []):
                    if var_dec.get('visibility') in ('public', 'internal', 'private'):
                        state_vars[var_dec.get('name')] = 1

            for value in node.values():
                traverse(value)
        elif isinstance(node, list):
            for item in node:
                traverse(item)

    traverse(json_data)
    return state_vars


if __name__ == '__main__':
    json_files = [f for f in os.listdir("AST_JSON") if f.endswith("_AST.json")]
    print(f"Found {len(json_files)} JSON files in AST_JSON.")

    os.makedirs("CSV", exist_ok=True)

    for json_file in tqdm(json_files):
        json_path = os.path.join("AST_JSON", json_file)

        try:
            with open(json_path, "r", encoding='utf-8') as f:
                json_data = json.load(f)
        except (FileNotFoundError, json.JSONDecodeError, UnicodeDecodeError) as e:
            print(f"Error while opening {json_path}: {e}")
            continue

        contract_address = json_file.split('_AST.json')[0]
        state_vars = extract_state_variables(json_data)
        function_var_map = extract_variable_relationship(json_data, state_vars)

        output_csv_path = os.path.join("CSV", f"{contract_address}.csv")

        with open(output_csv_path, 'w', newline='') as csvfile:
            csvwriter = csv.writer(csvfile)
            csvwriter.writerow(['Function Name', 'State Variables Modified', 'Can Transfer Ether'])
            for item in function_var_map:
                csvwriter.writerow(
                    [item['function_name'], ', '.join(item['variables']), item['is_transfer_capable']])
    