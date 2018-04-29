#!/usr/bin/python
# -*- coding:utf-8 -*-
import numpy as np
import matplotlib.pyplot as plt
import statsmodels.api as sm
import sys
import os

filename = "result/round5/cache/benchmark_read_info_172.16.1.92.txt"

data = []

def getData(path):
    data = []
    for f in os.listdir(path):
        filename = os.path.join(path,f)
        if os.path.isfile(filename) and f.endswith(".log"):
            with open(filename, 'r') as fp:
                lines = fp.read().split('\n')
                for line in lines:
                    try:
                        if line.startswith("Completed"):
                            data.append(float(line.split(' ')[-1][:-4]))
                    except ValueError,e:
                        print("error append",e, line)
    return data

def main():
    for path in sys.argv[1:]:
        data = getData(path)
        dataLen = len(data)
        print dataLen
        data = data[:dataLen/500 * 500]
        newData = data[::dataLen/500]
        newDataLen = len(newData)
        print newDataLen

        x = [x+1 for x in range(newDataLen)]
        #y = ecdf(x)
        plt.plot(x,newData,label=path.split('/')[-1])

    plt.legend(loc="best")
    #plt.show()
    plt.savefig('bw.eps', format='eps', dpi=1000)


if __name__ == '__main__':
    main()
