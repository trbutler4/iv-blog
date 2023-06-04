import * as React from 'react';
import {
    Button,
} from '@chakra-ui/react';
import detectEthereumProvider from '@metamask/detect-provider';

declare global {
    interface Window {
        ethereum?: any;
    }
}

function ConnectWalletButton() {
    const [buttonText, setButtonText] = React.useState("Connect Wallet");

    async function connectWallet() {
        const provider = await detectEthereumProvider();
        if (provider) {

            // get accounts
            const accounts = await window.ethereum.request({ 
                method: 'eth_requestAccounts' 
            }).catch((err: any) => {
                if (err.code === 4001) {
                    console.log("user rejected request")
                } else {
                    console.log(err)
                }
            })

            // connect account 
            const account = accounts[0]
            console.log(`account: ${account}`)
            setButtonText(`${account.slice(0, 6)}...`)
            
        } else {
            console.log("metamask not detected")
        }
    }

    return (
        <Button onClick={connectWallet}>{`${buttonText}`}</Button>
    )
}

export default ConnectWalletButton;
