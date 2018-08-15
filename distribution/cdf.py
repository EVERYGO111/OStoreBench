#!/usr/bin/python
# -*- coding:utf-8 -*-
import numpy as np
import matplotlib.pyplot as plt
import statsmodels.api as sm
import sys
import os


data = []

def getData(filename):
    data = []
    if os.path.isfile(filename) :
        with open(filename, 'r') as fp:
            while 1:
                line = fp.readline()
                if not line:
                    break
                try:
                    #data.append(float(line)/1000000)
                    data.append(float(line)+1)
                except ValueError,e:
                    print("error append",e, line)
    return data


def main():
    for path in sys.argv[1:]:
        data = getData(path)
        ecdf = sm.distributions.ECDF(data)
        x = np.linspace(min(data), 15, 100)
        y = ecdf(x)
        plt.plot(x,y, label=path.split('/')[-1])

    plt.legend(bbox_to_anchor=(0.65,0.3), loc=2,borderaxespad=0.)
    plt.show()


if __name__ == '__main__':
    main()
