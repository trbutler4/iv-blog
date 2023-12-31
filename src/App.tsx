import { Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar'
import Home from './pages/Home';
import About from './pages/About';
import Footer from './components/Footer';
import Article from './pages/Article';
import React from 'react';
import "./index.css";

function App() {
    return (
        <>
        <Navbar />

        <Routes>
            <Route path="/" element={<Home/>}/>
            <Route path="/about" element={<About/>}/>
            <Route path="/articles/:articleId" element={<Article/>}/>
        </Routes>
        
        {/* 
        <Footer />
        */}
        </>
    );
}

export default App;
