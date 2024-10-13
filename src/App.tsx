import { Routes, Route } from 'react-router-dom'
import Navbar from './components/Navbar'
import Home from './pages/Home';
import Article from './pages/Article';
import "./index.css";

function App() {
  return (
    <>
      <Navbar />

      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/articles/:articleId" element={<Article />} />
      </Routes>

      {/* 
        <Footer />
        */}
    </>
  );
}

export default App;
