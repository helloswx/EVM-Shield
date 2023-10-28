import os
import json
from tqdm import tqdm
from solidity_parser import parser

def make_serializable(obj):
    if isinstance(obj, dict):
        return {key: make_serializable(value) for key, value in obj.items()}
    elif isinstance(obj, list):
        return [make_serializable(elem) for elem in obj]
    elif hasattr(obj, '__dict__'):
        return make_serializable(obj.__dict__)
    else:
        return obj

contracts_folder = 'contracts' 
ast_folder = 'AST_JSON'

# 创建AST文件夹，如果不存在的话
if not os.path.exists(ast_folder):
    os.makedirs(ast_folder)

# 获取所有.sol文件
sol_files = [f for f in os.listdir(contracts_folder) if f.endswith('.sol')]



# 使用tqdm库显示进度条
for sol_file in tqdm(sol_files):
    with open(os.path.join(contracts_folder, sol_file), 'r', encoding='utf-8') as f:
        source_code = f.read()

    # 使用solidity_parser解析源代码
    try:
        ast = parser.parse(source_code)
        ast = make_serializable(ast)  # 转换为可序列化对象
    except Exception as e:
        print(f"Error parsing {sol_file}: {e}")
        continue

    # 将AST保存为JSON文件
    ast_file_name = sol_file.replace('.sol', '_AST.json')
    with open(os.path.join(ast_folder, ast_file_name), 'w', encoding='utf-8') as f:
        json.dump(ast, f, ensure_ascii=False, indent=4)
