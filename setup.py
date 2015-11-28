
from distutils.core import setup, Extension

module1 = Extension('fastopen', sources = ['fastopenmodule.c'])

setup(
        name = 'TcpFastOpen',
        version = '0.1',
        ext_modules = [module1]
        )
