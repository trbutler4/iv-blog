import detectEthereumProvider from '@metamask/detect-provider';

declare global {
    interface Window {
        ethereum?: any;
    }
}

export async function getAccount() {
    const provider = await detectEthereumProvider();
    if (provider) {

        // get accounts
        const accounts = await window.ethereum.request({ 
            method: 'eth_requestAccounts' 
        }).catch((err: any) => {
            if (err.code === 4001) {
                console.log("user rejected request")
                return err
            } else {
                console.log(err)
                return err
            }
        })

        // connect account 
        const account = accounts[0]
        console.log(`connected to: ${account}`)
        return account;
        
    } else {
        console.log("metamask not detected")
    }
}

export async function isAdmin() {
    const account = await getAccount();
    if (account === "0xf0017157524793829f29fb67a4df21eaeee87578") {
        console.log(`${account} is admin`)
        return true;
    }
    
    console.log(`${account} is not admin`)
    return false;
}
