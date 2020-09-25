pragma solidity 0.4.24;

interface IMeterNative {
    function native_mtr_totalSupply() external view returns(uint256); 
    function native_mtr_totalBurned() external view  returns(uint256); 
    function native_mtr_get(address addr) external view returns(uint256); 
    function native_mtr_add(address addr, uint256 amount) external;
    function native_mtr_sub(address addr, uint256 amount) external returns(bool); 

    //@@@@@
    function native_mtrg_totalSupply() external view returns(uint256);
    function native_mtrg_totalBurned() external view returns(uint256); 
    function native_mtrg_get(address addr) external view returns(uint256); 
    function native_mtrg_add(address addr, uint256 amount) external; 
    function native_mtrg_sub(address addr, uint256 amount) external returns(bool);

    //@@@
    function native_master(address addr) external view returns(address); 
}
