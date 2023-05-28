import { Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar'
import Home from './pages/Home';
import About from './pages/About';
import React from 'react';

function App() {
    return (
        <>
        <Navbar />

        <Routes>
            <Route path="/" element={<Home/>}/>
            <Route path="/about" element={<About/>}/>
        </Routes>

        </>
    );
}

export default App;
