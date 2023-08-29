// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

contract QuickSort {
    uint256 public seed;

    function random() internal returns (uint256) {
        seed = (1103515245 * uint(seed) + 12345) % (2 ** 32);
        return seed;
    }

    function randomArray(uint256 size) public returns (uint256[] memory) {
        uint256[] memory arr = new uint256[](size);
        for (uint256 i = 0; i < size; i++) {
            arr[i] = random();
        }
        return arr;
    }

    function quickSort(
        uint256[] memory arr,
        uint256 left,
        uint256 right
    ) internal pure {
        uint256 i = left;
        uint256 j = right;
        uint256 pivot = arr[(left + right) / 2];
        while (i <= j) {
            while (arr[i] < pivot) i++;
            while (pivot < arr[j]) j--;
            if (i <= j) {
                (arr[i], arr[j]) = (arr[j], arr[i]);
                i++;
                j--;
            }
        }
        if (left < j) quickSort(arr, left, j);
        if (i < right) quickSort(arr, i, right);
    }

    function benchmark() public returns (uint256) {
        uint256 checksum = 0;
        for (uint256 i = 0; i < 100; i++) {
            uint256[] memory arr = randomArray(1000);
            quickSort(arr, 0, arr.length - 1);
            checksum += arr[100];
        }
        return checksum;
    }
}
