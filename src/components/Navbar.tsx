import ConnectWalletButton from './ConnectWalletButton';

const Navbar = () => {
    return (

        <div className="flex flex-col justify-center  items-center">
            <div className="text-left text-xl font-bold font-mono">
                IV
            </div>

            <div>
                <button className="ml-2">Home</button>
                <button className="ml-2">About</button>
            </div>
        </div>

    )
}

export default Navbar;
