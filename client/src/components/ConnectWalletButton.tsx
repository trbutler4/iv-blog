import React, {useEffect, useState} from 'react';
import {
    Button,
} from '@chakra-ui/react';
import { getAccount } from '../interfaces/metamask_interface';

function ConnectWalletButton() {
    const [account, setAccount] = useState("");

    const fetchAccount = async () => {
        const account = await getAccount();
        if (account) {
            setAccount(account)
        }
    }

    const connectAccount = async () => {
        const account = await getAccount();
        if (account) {
            setAccount(account)
        }
    }

    useEffect(() => {
        fetchAccount()
    }, []);


    return (
        <>
        {account ? (
            <Button> {`${account.slice(0,6)}...`} </Button>
        ) : (
            <Button onClick={connectAccount}>Connect Account</Button>
        )}
        </>
    )
}

export default ConnectWalletButton;
