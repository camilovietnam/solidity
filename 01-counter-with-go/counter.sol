// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;

contract Counter {
        uint8 count;

        constructor () {
            reset();
        }

        function increase() public {
            count++;
        }

        function getCount() public view returns (uint8) {
            return count;
        }

        function reset() public {
            count = 0;
        }
}