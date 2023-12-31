// SPDX-License-Identifier: MIT
pragma solidity ^0.8.13;

// import "forge-std/console2.sol";

contract QuickSortBenchmark {
    uint256 public seed;

    function random() internal returns (uint256) {
        seed = (1103515245 * uint(seed) + 12345) % (1 << 31);
        return seed;
    }

    function randomizeArray(uint256[] memory arr) internal {
        for (uint256 i = 0; i < arr.length; i++) {
            arr[i] = random() % 1000;
        }
    }

    function quicksort(uint[] memory arr, int left, int right) internal {
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
        if (left < j) quicksort(arr, left, j);
        if (i < right) quicksort(arr, i, right);
    }

    function run(uint256 _seed, uint256 arrLen, uint256 iter) public returns (uint256) {
        seed = _seed;
        uint256 checksum = 0;
        uint256[] memory arr = new uint256[](arrLen);
        for (uint256 i = 0; i < iter; i++) {
            randomizeArray(arr);
            quicksort(arr, 0, int256(arrLen - 1));
            checksum += arr[arrLen / 2];
        }
        return checksum;
    }
}
