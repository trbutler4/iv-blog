import { Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar'
//import Home from './pages/Home';
import About from './pages/About';
import Blog from './pages/Blog';
import Footer from './components/Footer';
import Post from './pages/Post';
import UploadPost from './pages/UploadPost';
import { ChakraProvider } from '@chakra-ui/react';

function App() {
    return (
        <ChakraProvider>
            <Navbar />

            <Routes>
                <Route path="/" element={<Blog/>}/>
                <Route path="/about" element={<About/>}/>
                <Route path="/blog" element={<Blog/>}/>
                <Route path="/post" element={<Post/>}/>
                <Route path="/upload" element={<UploadPost/>}/>
            </Routes>

            <Footer />
        </ChakraProvider>
    );
}

export default App;
