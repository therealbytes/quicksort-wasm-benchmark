// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

import "forge-std/console2.sol";

contract QuickSort {
    uint256 public seed;

    function random() internal returns (uint256) {
        seed = (1103515245 * uint(seed) + 12345) % (1 << 31);
        return seed;
    }

    function randomizeArray(uint256[] memory arr) internal {
        for (uint256 i = 0; i < arr.length; i++) {
            arr[i] = random();
        }
    }

    function quickSort(uint[] memory arr, int left, int right) internal {
        int i = left;
        int j = right;
        if (i == j) return;
        uint pivot = arr[uint(left + (right - left) / 2)];
        while (i <= j) {
            while (arr[uint(i)] < pivot) i++;
            while (pivot < arr[uint(j)]) j--;
            if (i <= j) {
                (arr[uint(i)], arr[uint(j)]) = (arr[uint(j)], arr[uint(i)]);
                i++;
                j--;
            }
        }
        if (left < j) quickSort(arr, left, j);
        if (i < right) quickSort(arr, i, right);
    }

    function benchmark() public returns (uint256) {
        seed = 7;
        uint256 checksum = 0;
        uint256[] memory arr = new uint256[](1000);
        for (uint256 i = 0; i < 100; i++) {
            randomizeArray(arr);
            quickSort(arr, 0, 999);
            checksum += arr[100];
        }
        return checksum;
    }
}
