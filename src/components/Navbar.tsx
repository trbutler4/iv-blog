import React from 'react';
import { Link } from 'react-router-dom'

const Navbar = () => {
    return (
        <div>
            <h1>IV</h1>
            <li>
                <Link to="/">Home</Link>
            </li>
            <li>
                <Link to="/about">About</Link>
            </li>
        </div>
    )
}

export default Navbar;
