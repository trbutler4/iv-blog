import ConnectWalletButton from './ConnectWalletButton';
import { Link } from 'react-router-dom';

const Navbar = () => {
    return (

        <div className="flex flex-col justify-center  items-center">
            <div className="text-left text-xl font-bold font-mono mt-8 mb-4">
                IV
            </div>

            <div>
                <Link to="/" className="ml-2">Home</Link>
                <Link to="/about" className="ml-2">About</Link>
            </div>
        </div>

    )
}

export default Navbar;
