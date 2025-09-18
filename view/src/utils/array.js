export function isArrEqual(arr1, arr2){
    return arr1.length === arr2.length && arr1.every((ele) => arr2.includes(ele));
};