import openpyxl
import pandas as pd


def workbook2df(path, first_row_header, first_column_index):
    workbook = openpyxl.load_workbook(path).active
    rows = list(workbook.iter_rows(values_only=True))
    df = pd.DataFrame(rows)
    if first_row_header:
        new_header = df.iloc[0]
        df = df[1:]
        df.columns = new_header
    if first_column_index:
        df = df.set_index("SL")
    return df
