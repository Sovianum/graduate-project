
import pandas as pd


def read_single_compressor_data(file_path):
    return pd.read_csv(file_path, names=['pi', 'G', 'N_e', 'eta'])


def read_double_compressor_data(file_path):
    return pd.read_csv(file_path, names=['pi', 'pi_factor', 'G', 'N_e', 'eta'])


def get_max_eta_df(double_compressor_df):
    pi_factor = double_compressor_df[double_compressor_df.eta == double_compressor_df.eta.max()].pi_factor.values[0]
    return double_compressor_df.groupby(['pi_factor']).get_group(pi_factor)


def get_max_power_df(double_compressor_df):
    pi_factor = double_compressor_df[double_compressor_df.N_e == double_compressor_df.N_e.max()].pi_factor.values[0]
    return double_compressor_df.groupby(['pi_factor']).get_group(pi_factor)