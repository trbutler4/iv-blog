import * as React from 'react';
import {
    Button,
} from '@chakra-ui/react';
import { getAccount } from '../interfaces/metamask_interface';

declare global {
    interface Window {
        ethereum?: any;
    }
}

function ConnectWalletButton() {
    const [buttonText, setButtonText] = React.useState("Connect Wallet");

    function connectWallet() {
        getAccount().then((account) => {
            setButtonText(`${account.slice(0,6)}...`)
        })
    }

    return (
        <Button onClick={connectWallet}>{`${buttonText}`}</Button>
    )
}

export default ConnectWalletButton;
