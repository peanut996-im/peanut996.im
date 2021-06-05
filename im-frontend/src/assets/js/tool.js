import {sha1} from './sha1.js';
let appKey = '624D8D49C3ECC20E9467A0395BCF0D6A';
let makeSign = (obj) =>{
    let plain = '';
    Object.keys(obj).filter(key => key !== 'sign'&& key !== 'EIO' && key !== 'transport' && key !== 't')
        .sort().forEach(key => {
        plain+=`${key}${obj[key]}`;
    });
    plain += appKey;
    console.log('befor sha1: ',plain);
    console.log(plain);
    return sha1(plain).toUpperCase();
}

export {
    makeSign,
    sha1
}
